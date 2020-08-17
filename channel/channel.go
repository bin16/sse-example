package channel

import (
	"fmt"
	"strings"

	"github.com/bin16/sse-demo/store"
)

type Channel chan string

var channelMap = map[string]Channel{}

func Subscribe(id, sid string) Channel {
	fmt.Println(id, sid, "Subscribe")
	longID := strings.Join([]string{id, sid}, "/")
	store.Add(id, sid)
	if channelMap[longID] == nil {
		channelMap[longID] = make(Channel, 5)
	}
	return channelMap[longID]
}

func UnSubscribe(id, sid string) {
	fmt.Println(id, sid, "UnSubscribe")
	longID := strings.Join([]string{id, sid}, "/")
	if channelMap[longID] != nil {
		close(channelMap[longID])
	}
	store.Del(id, sid)
}

func Post(id, message string) {
	fmt.Println(id, "POST")
	l := store.Get(id)
	for _, d := range l {
		sid := d[1]
		longID := strings.Join([]string{id, sid}, "/")
		if channelMap[longID] == nil {
			UnSubscribe(id, sid)
		} else {
			channelMap[longID] <- message
		}
	}
}
