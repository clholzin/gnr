package gnr

// TuneInOnce allows the channel returned from Subscribe to listen to the event but not block allowing for infinit loop checking
// ** event.TuneIn()
func (event Channel) TuneIn() bool {
	select {
	case ev := <-event:
		ev.Reply <- PubSub{}
		return true
	default:
	}
	return false
}

// TuneInBlock allows the channel returned from Subscribe to listen to the event and block current thread till triggered
// ** event.TuneInBlock()
func (event Channel) TuneInBlock() bool {
	select {
	case ev := <-event:
		ev.Reply <- PubSub{}
		return true
	}
}

// TuneInOnce allows the channel returned from Subscribe to listen to the event but not block allowing for infinit loop checking, once triggered event unsubscribes
// ** event.TuneInOnce()
func (event Channel) TuneInOnce() bool {
	select {
	case ev := <-event:
		ev.Reply <- SubStop{}
		return true
	default:
	}
	return false
}

// TuneInBlockOnce allows the channel returned from Subscribe to listen to the event and block once on event trigger till triggered, once triggered event unsubscribes
// ** event.TuneInBlockOnce()
func (event Channel) TuneInBlockOnce() bool {
	select {
	case ev := <-event:
		ev.Reply <- SubStop{}
		return true
	}
}

// Drop allows the channel to unsubscribe to the event key name,
// ** event.Drop()
func (event Channel) Drop(key string) bool {
	var carrier ChanList
	carrier = append(carrier, event)
	action := deleteChans
	finalize := syncRadios{action, key, carrier, make(chan DataReply, 1)}
	radioFlow <- finalize
	// used for dropping On events
	select {
	case event <- PubSub{key, kill, make(chan Event, 1)}:
	case <-event:
		event <- PubSub{key, kill, make(chan Event, 1)}
	default:
	}

	return true
}

// On allows the channel to independently instantiate a function callback upon being triggered
func (event Channel) On(fn func(PubSub)) {
	go func() {
		for {
			select {
			case p, ok := <-event:
				p.Reply <- PubSub{}
				if !ok { // catch closed channel
					return
				}
				if p.Action == kill {
					return
				}
				fn(p)
			}
		}
	}()
}
