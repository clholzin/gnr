## G.N.R - Go.National.Radio

G.N.R allows golang applications to subscribe and trigger events globally. No limit
 on Subscribers which means many locations within an application can listen for the event
 and take action separately.

##### Subscribe
Pass event name as a string to receive a channel
as a subscriber to the event name. This allows for a handle on
global events taking place when that Name is triggered.

```go
channel,err := gnr.Subscribe("EventName")
```

##### On
On allows a channel to be continuously listened to for when Trigger occurs and a callback function will be called, On is none blocking to main thread
```go
channel.On(func(p PubSub){
   // code here will be instantiated each time Trigger is called on that nameed event
})
```
##### Close a ```On``` event listner
Once channel.Drop is successful, when dropping any event one last send will be received and validated for On() events to close the goroutine
```go
if channel.Drop("EventName") {
    //success
}
```

##### TuneIn
TuneIn will provide a way to listen for the event as a method on the
channel type. Function returns true for event triggered and false. Good for using within loops
for continuous checking. Function returns true for event triggered
and false.
```go
if channel.TuneIn() {
   // channel triggered
}
```
##### TuneInBlock
TuneInBlock will provide a way to listen for the event as a method on
the channel type while blocking the thread. Function returns true for event triggered
and false.
```go
channel.TuneInBlock()
```

##### TuneInOnce
TuneInOnce will provide a way to listen for the event as a method on the
channel one time before automatically notifying to be removed, similar TuneIn.
Good for using within loops for continuous checking. Function returns true for event triggered
and false.
```go
if channel.TuneInOnce(){
    // channel triggered
}
```

##### TuneInBlockOnce
TuneInBlockOnce will provide a way to listen for the event as a method on the
channel while blocking the thread one time before automatically notifying to be removed,
similar TuneInBlock. Function returns true for event triggered
and false.
```go
channel.TuneInBlockOnce()
```

##### Drop
Drop will provide a way to remove the subscriber from event list memory.
```go
channel.Drop()
```

##### Trigger
Pass Event name and if the type is "async" or "sync"
to wait.

Sync will cause Trigger to wait until all channels have been served and replied to
if they are sending a response.

```go
import "path/to/gnr"

gnr.Trigger("EventName","sync")
```

Async will not wait for confirmation and will continue down the call stack.

```go
import "path/to/gnr"

gnr.Trigger("EventName","async")
```


Basic Application Usage:

```go
import "path/to/gnr"

// Infinitly listen for event

event,err := gnr.Subscribe("EventName")
if err != nil {
  fmt.Println("EventName :"+err.Error())
  //"EventName: key name empty"
}

// Manual implementation to listen for event trigger
go func(){
    for {
         select {
         case e := <- event:
            // Optional Reply
            // This is used to help synchronize and also to remove subscribers
            e.Reply <- gnr.PubSub{}

         }
    }
}()

////////////// OR use On() on the channel as a none blocking callback function

event.On(func(p PubSub){
   // code here will be instantiated each time Trigger is called on that nameed event
})

////////////// OR use TuneIn() on the channel

go func(){
   for {
        if event.TuneIn() {
           // handle event trigger
        }
   }

}()

////////////// OR use TuneInOnce() on the channel

go func(){
   for {
        if event.TuneInOnce() {
            break
           // handle event trigger
        }
   }

}()

////////////// OR use TuneInBlock() on the channel

go func(){
   for {
        if event.TuneInBlock() {
            break
            // handle event trigger
        }
   }

}()

////////////// OR use TuneInBlockOnce() on the channel

if event.TuneInBlockOnce() {
    // handle event trigger
}


// Finally trigger the event name

gnr.Trigger("EventName")

```

```go

event,err := gnr.Subscribe("ListenForNoEvent")

//It is also possible to listen for the false value. This allows for listening for
no event trigger.

if !event.TuneIn() {
    // Having the option to catch a none triggered event allows the user to
    // decide when the program removes the event subscriber using eventName.Drop()

    event.Drop()
    
    // handle event trigger
}

// Finally trigger the event name

gnr.Trigger("ListenForNoEvent")

```

Listen for only 10 consecutive events

```go
import "path/to/gnr"

func listnerForEventName (subscriber chan PubSub) PubSub{
	pubSub :=<- subscriber
	log.Println("Received from event: "+pubSub.Key)
	return pubSub
}

func main(){

    subscriber,err := Subscribe("PubSubRestrict")
    if err != nil{
        t.Fail()
    }

    // Listen for event 10 times
    var pubsub PubSub
    go func(){
        for i := 0;i < 10;i++{
            pubsub = listnerForEventName(subscriber)
        }
        pubsub.Reply <- SubStop{}
    }()


    for i := 0;i < 10;i++{
        Trigger("PubSubRestrict","async")
    }
}



```



##### Close Subscriber
Reply back on a channel with an empty "SubStop" struct to close and remove placement in event array

```go

subscribed,err := Subscribe("ListenForItClose")
if err != nil {
    log.Println(err)
}
go func(){
    select {
    case v := <- subscribed:
        log.Println("Got it, received trigger")
        closeMe := SubStop{}
        v.Reply <- closeMe
    }
}()
log.Println("Trigger Event: ListenForItClose")
Trigger("ListenForItClose","sync")
log.Println("Trigger Event Close...")
