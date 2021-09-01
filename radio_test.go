package gnr

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestPubSubSync(t *testing.T) {
	subscribed, err := Subscribe("ListenForIt")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		select {
		case v := <-subscribed.C.Signal:
			testClose := SubStop{}
			v.Reply <- testClose
		}
	}()
	Trigger("ListenForIt", "sync")
}

func TestPubSubAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	subscribed, err := Subscribe("ListenForAsync")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		select {
		case <-subscribed.C.Signal:
			wg.Done()
		}
	}()
	Trigger("ListenForAsync", "async") // won't block
	wg.Wait()
}

func TestPubSubSyncTuneIn(t *testing.T) {
	subscribed, err := Subscribe("ListenForItTuneIn")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		for {
			if subscribed.TuneIn() {
				break
			}
		}
	}()
	Trigger("ListenForItTuneIn", "sync")
}

func TestPubSubSyncTuneInBlock(t *testing.T) {
	subscribed, err := Subscribe("ListenForItTuneInBlock")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		if subscribed.TuneInBlock() {
		}
	}()
	Trigger("ListenForItTuneInBlock", "sync")
}

func TestPubSubSyncTuneInBlockOnce(t *testing.T) {
	subscribed, err := Subscribe("ListenForItTuneInBlockOnce")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		if subscribed.TuneInBlockOnce() {
		}
	}()
	Trigger("ListenForItTuneInBlockOnce", "sync")
}

func TestPubSubAsyncTuneIn(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	subscribed, err := Subscribe("ListenForAsyncTuneIn")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		for {
			if subscribed.TuneIn() {
				wg.Done()
				break
			}
		}

	}()
	Trigger("ListenForAsyncTuneIn", "async") // won't block
	wg.Wait()
}

func TestPubSubAsyncReply(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	subscribed, err := Subscribe("ListenForPubSubAsyncReply")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		select {
		case v := <-subscribed.C.Signal:
			testClose := SubStop{}
			v.Reply <- testClose
			wg.Done()
		}
	}()
	Trigger("ListenForPubSubAsyncReply", "async") // won't block
	wg.Wait()
}

func TestPubSubAsyncReplyTunInOnce(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	eventName := "ListenForPubSubAsyncReplyTunInOnce"
	subscribed, err := Subscribe(eventName)
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		for {
			if subscribed.TuneInOnce() {
				wg.Done()
				break
			}
		}
	}()
	Trigger(eventName, "async") // won't block
	wg.Wait()
}

func TestPubSubAsyncReplyTuneInBlockOnce(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	eventName := "ListenForPubSubAsyncReplyTuneInBlockOnce"
	subscribed, err := Subscribe(eventName)
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		if subscribed.TuneInBlockOnce() {
			wg.Done()
		}
	}()
	Trigger(eventName, "async") // won't block
	wg.Wait()
}

func TestPubSubClose(t *testing.T) {
	subscribed, err := Subscribe("ListenForItClose")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		select {
		case v := <-subscribed.C.Signal:
			v.Reply <- SubStop{}
		}
	}()
	Trigger("ListenForItClose", "sync")
}

func TestPubSubListenAndRemove(t *testing.T) {
	eventName := "ListenAndRemoveFromTriggerList"
	subscribedOne, err := Subscribe(eventName)
	if err != nil {
		t.Fail()
		return
	}
	subscribedTwo, err := Subscribe(eventName)
	if err != nil {
		t.Fail()
		return
	}
	subscribedOne.Drop()
	subscribedTwo.Drop()
	channels := Access(eventName)
	if channels.Size > 0 {
		t.Fail()
	}
}

func TestPubSubCloseAsync(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	subscribed, err := Subscribe("ListenForItCloseAsync")
	if err != nil {
		t.Fail()
		return
	}
	go func() {
		select {
		case v := <-subscribed.C.Signal:
			v.Reply <- SubStop{}
			wg.Done()
		}
	}()
	Trigger("ListenForItCloseAsync", "async")
	wg.Wait()
}

func TestPubSubRestrictSync(t *testing.T) {
	subscriber, err := Subscribe("PubSubRestrictSync")
	if err != nil {
		t.Fail()
	}
	go func() {

		var pubsub PubSub
		// Listen for event 10 times
		length := 10
		for i := 1; i <= length; i++ {
			pubsub = listnerForEventName(subscriber.C.Signal)
			if i == length {
				pubsub.Reply <- SubStop{}
			} else {
				pubsub.Reply <- PubSub{} // reply back to bypass timeout
			}
		}

	}()

	for i := 1; i <= 11; i++ { // call one extra time to test availabilty
		Trigger("PubSubRestrictSync", "sync")
	}
}

func TestPubSubRestrictAsync(t *testing.T) {

	subscriber, err := Subscribe("PubSubRestrictAsync")
	if err != nil {
		t.Fail()
	}
	go func() {
		var pubsub PubSub
		// Listen for event 10 times
		for i := 1; i <= 10; i++ {
			pubsub = listnerForEventName(subscriber.C.Signal)
		}
		pubsub.Reply <- SubStop{}
	}()

	for i := 0; i < 11; i++ {
		Trigger("PubSubRestrictAsync", "async")
	}

}

func TestOnSubscriber(t *testing.T) {

	event, _ := Subscribe("name")
	tester := make(chan int, 1)
	event.On(func(p PubSub) {
		// callback
		tester <- 1
	})

	Trigger("name", "async")

	ticker := time.NewTicker(time.Second)
	select {
	case <-tester:
	case <-ticker.C:
		t.Fail()
	}
}

func TestOnSubscriberCloser(t *testing.T) {

	event, _ := Subscribe("name")
	event.On(func(p PubSub) {})

	Trigger("name", "async")
	event.Drop()
	close(event.C.Signal)   // close will initiate termination of the event breaking event.On loop
	Trigger("name", "sync") //won't trigger event as it has been removed

}

func TestSize(t *testing.T) {

	_, _ = Subscribe("evName")

	size := Size("evName")
	if size == 0 {
		t.Fail()
	}

}

func listnerForEventName(subscriber chan PubSub) PubSub {
	pubSub := <-subscriber
	return pubSub
}

func BenchmarkPubSubSyncTuneInBlock(b *testing.B) {
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubSyncTuneInBlock" + strconv.Itoa(time.Now().Nanosecond())
	for i := 0; i < b.N; i++ {
		go func() {
			subscribed, err := Subscribe(eventName)
			if err != nil {
				b.Fail()
				return
			}

			subscribed.TuneInBlock()
		}()
	}
	Trigger(eventName, "sync")
}

func BenchmarkPubSubAsyncTuneInBlock(b *testing.B) {
	wg := sync.WaitGroup{}
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubAsyncTuneInBlock" + strconv.Itoa(time.Now().Nanosecond())
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		subscribed, err := Subscribe(eventName)
		if err != nil {
			b.Fail()
			return
		}
		go func() {
			subscribed.TuneInBlock()
			wg.Done()
		}()
	}
	Trigger(eventName, "async")
	wg.Wait()
}

func BenchmarkPubSubSync(b *testing.B) {
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubSync_" + strconv.Itoa(time.Now().Nanosecond())
	for i := 0; i < b.N; i++ {
		subscribed, err := Subscribe(eventName)
		if err != nil {
			b.Fail()
			return
		}
		go func(event *Node) {
			select {
			case v := <-event.C.Signal:
				v.Reply <- PubSub{}
			}
		}(subscribed)
	}

	Trigger(eventName, "sync")
}

func BenchmarkPubSubAsync(b *testing.B) {
	wg := sync.WaitGroup{}
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubAsync" + strconv.Itoa(time.Now().Nanosecond())
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		subscribed, err := Subscribe(eventName)
		if err != nil {
			b.Fail()
			return
		}
		go func(event *Node) {
			select {
			case v := <-event.C.Signal:
				v.Reply <- PubSub{}
				wg.Done()
			}
		}(subscribed)
	}
	Trigger(eventName, "async")
	wg.Wait()
}

func BenchmarkPubSubSyncTrigger(b *testing.B) {
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubSyncTrigger" + strconv.Itoa(time.Now().Nanosecond())
	subscribed, err := Subscribe(eventName)
	if err != nil {
		b.Fail()
		return
	}
	var count int
	subscribed.On(func(sub PubSub) {
		count++
	})

	timer := time.Now()
	for i := 0; i < b.N; i++ {
		Trigger(eventName, "sync")
	}
	b.Logf("Triggered %d in %v seconds received %d", b.N, time.Now().Sub(timer).Seconds(), count)
}

func BenchmarkPubSubAsyncTrigger(b *testing.B) {
	radio = &radios{Subs: make(map[string]*NodeList)}
	eventName := "BenchmarkPubSubAsyncTrigger" + strconv.Itoa(time.Now().Nanosecond())
	subscribed, err := Subscribe(eventName)
	if err != nil {
		b.Fail()
		return
	}
	var count int
	subscribed.On(func(sub PubSub) {
		count++
	})

	timer := time.Now()
	for i := 0; i < b.N; i++ {
		Trigger(eventName, "async")
	}
	b.Logf("Triggered %d in %v seconds received %d", b.N, time.Now().Sub(timer).Seconds(), count)
}

func BenchmarkPubSubSyncTuneInBlockOnce(b *testing.B) {
	eventName := "BenchmarkPubSubSyncTuneInBlockOnce" + strconv.Itoa(time.Now().Nanosecond())
	radio = &radios{Subs: make(map[string]*NodeList)}
	for i := 0; i < b.N; i++ {
		subscribed, err := Subscribe(eventName)
		if err != nil {
			b.Fail()
			return
		}
		go func() {
			subscribed.TuneInBlockOnce()
		}()
	}
	Trigger(eventName, "sync")
}

func BenchmarkPubSubAsyncTuneInBlockOnce(b *testing.B) {
	wg := sync.WaitGroup{}
	eventName := "BenchmarkPubSubAsyncTuneInBlockOnce" + strconv.Itoa(time.Now().Nanosecond())
	radio = &radios{Subs: make(map[string]*NodeList)}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		subscribed, err := Subscribe(eventName)
		if err != nil {
			b.Fail()
			return
		}
		go func() {
			subscribed.TuneInBlockOnce()
			wg.Done()
		}()
	}
	Trigger(eventName, "async")
	wg.Wait()
}
