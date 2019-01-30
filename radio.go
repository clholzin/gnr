package gnr

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var GoroutineCount uint64

func init() {
	go running()
}

func running() {
	for channels := range radioFlow {
		switch channels.Action {
		case add:
			ch := channels.Data[0]
			list := radio[channels.Key]
			radio[channels.Key] = append(list, ch)
			channels.Reply <- PubSub{}
		case remove:
			if len(channels.Data) == 0 {
				delete(radio, channels.Key)
			}
			channels.Reply <- PubSub{}
		case deleteChans:
			if len(channels.Data) > 0 {
				if list, ok := radio[channels.Key]; ok {
					var clone ChanList
					if len(channels.Data) == 1 {
						length := len(list)
						removeChan := channels.Data[0]
						for i := 0; i < length; i++ {
							event := list[i]
							if event == removeChan {
								clone = append(clone, list[i+1:]...) //complete
								break
							} else {
								clone = append(clone, event)
							}
						}
					} else if len(channels.Data) > 1 {
						for i := 0; i < len(list); i++ {
							var found bool
							event := list[i]
							if len(channels.Data) == 0 {
								clone = append(clone, list[i:]...) //complete
								break
							}
							for k, removeChan := range channels.Data {
								if event == removeChan {
									found = true
									channels.Data = append(channels.Data[:k], channels.Data[k+1:]...)
									break
								}
							}
							if !found {
								clone = append(clone, event)
							}
						}
					}
					radio[channels.Key] = clone
					if len(clone) == 0 {
						delete(radio, channels.Key)
					}
				}
			}
			channels.Reply <- PubSub{}
		case size:
			v, ok := radio[channels.Key]
			var length int
			if ok {
				length = len(v)
				channels.Reply <- Found{length, make(ChanList, 0)}
				continue
			}
			channels.Reply <- Found{0, make(ChanList, 0)}
		case find:
			if v, ok := radio[channels.Key]; ok {
				channels.Reply <- Found{len(v), v}
			} else {
				channels.Reply <- Found{0, make(ChanList, 0)}
			}
		}

	}
}

func GoCounter() uint64 {
	return atomic.LoadUint64(&GoroutineCount)
}

// Access is a function that accepts a key string name for the event and returns the list of channels listening to that key
func Access(key string) Found {
	find := syncRadios{find, key, make(ChanList, 0), make(chan DataReply, 1)}
	radioFlow <- find
	data := <-find.Reply
	return data.(Found)
}

// Size provides the length of the map containing all the channels under that key
func Size(key string) Found {
	find := syncRadios{size, key, make(ChanList, 0), make(chan DataReply, 1)}
	radioFlow <- find
	data := <-find.Reply
	return data.(Found)
}

// Subscribe is a function that accepts a key string name for the event to subscribe to and returns a channel to listen on
func Subscribe(key string) (Channel, error) {
	if key != "" {
		ch := make(Channel, 1)
		list := ChanList{ch}
		command := syncRadios{add, key, list, make(chan DataReply, 1)}
		radioFlow <- command
		<-command.Reply
		return ch, nil
	}
	return nil, errors.New("key name empty")
}

// Trigger is a function that accepts a key and action for adding data to channels listening for the action
func Trigger(key, action string) {
	var waitGroup sync.WaitGroup
	var rmWait sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		atomic.AddUint64(&GoroutineCount, 1)
		defer atomic.AddUint64(&GoroutineCount, ^uint64(0))
		if action == "async" {
			waitGroup.Done()
		}
		channels := Access(key)
		if channels.Length > 0 {
			rmBuffer := make(chan Channel, channels.Length)
			for i := 0; i < channels.Length; i++ {
				value := PubSub{key, send, make(chan Event, 1)}
				event := channels.List[i]
				rmWait.Add(1)
				select {
				case event <- value: // makes it safe for repeat trigger calls
					go func(ev Channel) {
						atomic.AddUint64(&GoroutineCount, 1)
						defer atomic.AddUint64(&GoroutineCount, ^uint64(0))
						defer rmWait.Done()
						timeout := time.NewTimer(15 * time.Microsecond)
						select {
						case reply := <-value.Reply:
							switch reply.(type) {
							case SubStop:
								rmBuffer <- ev
							default:
								//Todo: allow custom types and actions or passing functions
							}
						case <-timeout.C:
						}

					}(event)
				default:
					rmWait.Done()
				}
			}
			rmWait.Wait()
			if len(rmBuffer) > 0 {
				go func() {
					atomic.AddUint64(&GoroutineCount, 1)
					defer atomic.AddUint64(&GoroutineCount, ^uint64(0))
					var clone ChanList
					length := len(rmBuffer)
					for t := 0; t < length; t++ {
						clone = append(clone, <-rmBuffer)
					}
					radioFlow <- syncRadios{deleteChans, key, clone, make(chan DataReply, 1)}
				}()
			}

		}
		if action == "sync" {
			waitGroup.Done()
		}
	}()

	waitGroup.Wait()
}
