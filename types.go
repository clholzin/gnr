package gnr

type (
	PubSub struct {
		Key    string
		Action int
		Reply  chan Event
	}
	Channel  chan PubSub
	ChanList []Channel

	Event     interface{}
	DataReply interface{}

	SubStop   struct{}
	radios    map[string]ChanList
	removeMap map[Channel]int

	Found struct {
		Length int
		List   ChanList
	}

	syncRadios struct {
		Action int
		Key    string
		Data   ChanList
		Reply  chan DataReply
	}
)

const (
	add = iota + 1
	remove
	deleteChans
	size
	find
	kill
	send
)

var (
	radio     = make(radios)
	radioFlow = make(chan syncRadios, 1000)
)
