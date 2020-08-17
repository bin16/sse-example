package channel

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMultiple(t *testing.T) {
	w := sync.WaitGroup{}
	id := "Foo"
	sids := []string{"a", "b", "c", "d", "e"}
	msgs := map[string]int{}
	count := 5
	for i := 0; i < count; i++ {
		msg := fmt.Sprintf("M%d", i+1)
		for r := 0; r < count; r++ {
			msgKey := fmt.Sprintf("%s#%s", msg, sids[r])
			msgs[msgKey] = 1
			w.Add(1)
		}
	}
	for _, sid := range sids {
		ch := Subscribe(id, sid)
		go func(sid string) {
			time.Sleep(time.Second * 5)
			for i := 0; i < count; i++ {
				select {
				case m := <-ch:
					w.Done()
					msgKey := fmt.Sprintf("%s#%s", m, sid)
					if msgs[msgKey] != 1 {
						t.Errorf("Message: <%s> \n(id=<%s>, sid=<%s>);\n want <%s> in <%v>", m, id, sid, msgKey, msgs)
					}
				}
			}
			UnSubscribe(id, sid)
		}(sid)
	}

	for i := 0; i < count; i++ {
		msg := fmt.Sprintf("M%d", i+1)
		Post(id, msg)
	}

	w.Wait()
}

func TestSingle(t *testing.T) {
	id := "Foo"
	sid := "sid-a"
	ch := Subscribe(id, sid)
	count := 5
	msgs := map[string]int{}
	go func() {
		time.Sleep(time.Second * 2)
		for i := 0; i < count; i++ {
			msg := fmt.Sprintf("Hello %d", i+1)
			msgs[msg] = 1
			Post(id, msg)
		}
	}()
	for i := 0; i < count; i++ {
		select {
		case m := <-ch:
			if msgs[m] != 1 {
				t.Errorf("Wrong message: %s (%d)", m, i)
			}
		}
	}
	UnSubscribe(id, sid)
}
