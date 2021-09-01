package gnr

import (
	"sync"

	"github.com/google/uuid"
)

type (
	PubSub struct {
		Key    string
		Action int
		Reply  chan Event
	}

	Channel struct {
		Signal chan PubSub
		Id     uuid.UUID
	}

	Node struct {
		C      *Channel
		Prev   *Node
		Next   *Node
		Parent *NodeList
	}

	NodeList struct {
		Head *Node
		Tail *Node
		Size int
	}

	Event     interface{}
	DataReply interface{}

	SubStop struct{}

	radios struct {
		sync.Mutex
		Subs map[string]*NodeList
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
	radio = &radios{Subs: make(map[string]*NodeList)}
)

func NewNode() *Node {
	return &Node{C: &Channel{Signal: make(chan PubSub, 1)}}
}
