package gnr

// TuneInOnce allows the channel returned from Subscribe to listen to the event but not block allowing for infinit loop checking
// ** event.TuneIn()
func (event *Node) TuneIn() bool {
	select {
	case ev := <-event.C.Signal:
		ev.Reply <- PubSub{}
		return true
	default:
	}
	return false
}

// TuneInBlock allows the channel returned from Subscribe to listen to the event and block current thread till triggered
// ** event.TuneInBlock()
func (event *Node) TuneInBlock() bool {
	select {
	case ev := <-event.C.Signal:
		ev.Reply <- PubSub{}
		return true
	}
}

// TuneInOnce allows the channel returned from Subscribe to listen to the event but not block allowing for infinit loop checking, once triggered event unsubscribes
// ** event.TuneInOnce()
func (event *Node) TuneInOnce() bool {
	select {
	case ev := <-event.C.Signal:
		ev.Reply <- SubStop{}
		return true
	default:
	}
	return false
}

// TuneInBlockOnce allows the channel returned from Subscribe to listen to the event and block once on event trigger till triggered, once triggered event unsubscribes
// ** event.TuneInBlockOnce()
func (event *Node) TuneInBlockOnce() bool {
	select {
	case ev := <-event.C.Signal:
		ev.Reply <- SubStop{}
		return true
	}
}

// Drop allows the channel to unsubscribe to the event key name,
// ** event.Drop()
func (event *Node) Drop() {
	event.Parent.Size--
	if event.Prev == nil {
		event.Parent.Head = event.Next
		if event.Parent.Tail == event {
			event.Parent.Tail = event.Next
		}
	} else {
		event.Prev.Next = event.Next
	}
}

// On allows the channel to independently instantiate a function callback upon being triggered
func (event *Node) On(fn func(PubSub)) {
	go func() {
		for {
			select {
			case p, ok := <-event.C.Signal:
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
