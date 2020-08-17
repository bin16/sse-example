package store

import "sync"

type DataItem [3]string

var data []DataItem = []DataItem{}
var rwl sync.RWMutex

func Add(id, sid string) {
	rwl.Lock()
	defer rwl.Unlock()
	if contains(id, sid) {
		return
	}
	data = append(data, DataItem{id, sid, ""})
}

func Del(id, sid string) {
	rwl.Lock()
	defer rwl.Unlock()
	for i, d := range data {
		if d[0] == id && d[1] == sid {
			data = append(data[:i], data[i+1:]...)
		}
	}
}

func Get(id string) []DataItem {
	r := []DataItem{}
	rwl.RLock()
	defer rwl.RUnlock()
	for _, d := range data {
		if d[0] == id {
			r = append(r, d)
		}
	}

	return r
}

func GetItem(id, sid string) DataItem {
	rwl.RLock()
	defer rwl.RUnlock()
	for _, d := range data {
		if d[0] == id && d[1] == sid {
			return d
		}
	}

	return DataItem{}
}

func Has(id, sid string) bool {
	rwl.RLock()
	defer rwl.RUnlock()
	return contains(id, sid)
}

func contains(id, sid string) bool {
	for _, d := range data {
		if d[0] == id && d[1] == sid {
			return true
		}
	}
	return false
}
