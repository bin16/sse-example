package store

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestSimple(t *testing.T) {
	id := "ID_A"
	sids := []string{"SID_0", "SID_1", "SID_2"}
	Add(id, sids[0])
	Add(id, sids[1])

	d := Get(id)
	for i, sid := range sids[:2] {
		if d[i][1] != sid {
			t.Errorf("Failed to run Get(%s); got %v", id, d)
		}
	}
	if !Has(id, sids[1]) {
		t.Errorf("Failed to run Has(%s, %s); got false", id, sids[1])
	}
	if Has(id, sids[2]) {
		t.Errorf("Failed to run Has(%s, %s); got true", id, sids[2])
	}

	Del(id, sids[1])
	if Has(id, sids[1]) {
		t.Errorf("Failed to run Has(%s, %s); got true", id, sids[1])
	}
}

func TestRW(t *testing.T) {
	id := "ID_B-RW"
	sids := []string{}
	for i := 0; i < 1000; i++ {
		sids = append(sids, fmt.Sprintf("%s:SID_%d", id, i))
	}
	for _, sid := range sids {
		go Add(id, sid)
	}
	// time.Sleep(time.Second * 2)
	for _, sid := range sids {
		if !Has(id, sid) {
			t.Errorf("Failed to run Has(%s, %s); got %v; want %v", id, sid, false, true)
		}
	}

	xUsed := map[int]int{}
	xSids := []string{}
	for i := 0; i < 300; i++ {
		rand.Seed(int64(42 * i))
		ix := rand.Intn(len(sids))
		if xUsed[ix] == 0 {
			xSid := sids[ix]
			xSids = append(xSids, xSid)
			go Del(id, xSid)
			xUsed[ix] = 1
		}
	}
	// time.Sleep(time.Second * 2)
	for _, xSid := range xSids {
		if Has(id, xSid) {
			t.Errorf("Failed to run Has(%s, %s); got %v; want %v", id, xSid, true, false)
		}
	}
}
