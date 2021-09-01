package gnr

import (
	"errors"
	"sync"
)

// Access is a function that accepts a key string name for the event and returns the list of channels listening to that key
func Access(key string) *NodeList {
	radio.Lock()
	data := radio.Subs[key]
	radio.Unlock()
	return data
}

// Size provides the length of the map containing all the channels under that key
func Size(key string) int {
	radio.Lock()
	var length int
	if node, ok := radio.Subs[key]; ok {
		length = node.Size
	}
	radio.Unlock()
	return length
}

// Subscribe is a function that accepts a key string name for the event to subscribe to and returns a channel to listen on
func Subscribe(key string) (*Node, error) {
	if key != "" {
		node := NewNode()
		radio.Lock()
		if nodelist, ok := radio.Subs[key]; ok {
			nodelist.Tail.Next = node
			nodelist.Tail = nodelist.Tail.Next
			node.Parent = nodelist
			nodelist.Size++
		} else {
			nodelist = &NodeList{}
			nodelist.Head = node
			nodelist.Tail = node
			radio.Subs[key] = nodelist
			node.Parent = nodelist
			nodelist.Size++
		}
		radio.Unlock()
		return node, nil
	}
	return nil, errors.New("key name empty")
}

// Trigger is a function that accepts a key and action for adding data to channels listening for the action
func Trigger(key, action string) {
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		if nodelist, ok := radio.Subs[key]; ok {
			if nodelist.Size == 0 {
				return
			}
			node := nodelist.Head
			for node != nil {
				var remove bool
				data := PubSub{key, send, make(chan Event, 1)}
				select {
				case node.C.Signal <- data:
					select {
					case reply := <-data.Reply:
						switch reply.(type) {
						case SubStop:
							remove = true
						default:
						}
					default:
					}
				default:
				}
				if remove {
					node.Drop()
				} else {
					node = node.Next
				}
			}
		}
		wait.Done()
	}()

	if action == "sync" {
		wait.Wait()
	}

}
