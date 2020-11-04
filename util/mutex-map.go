package util

import "sync"

//MutexMap map multiple mutex using key
type MutexMap struct {
	mtx sync.Mutex
	mp  map[interface{}]*payload
}

type payload struct {
	mtx sync.Mutex
	cnt uint
}

//NewMutexMap create a new MutexMap
func NewMutexMap() *MutexMap {
	return &MutexMap{mp: make(map[interface{}]*payload)}
}

//Lock try to lock the (key, value)
func (mm *MutexMap) Lock(key interface{}) {
	mm.mtx.Lock()
	e, ok := mm.mp[key]
	if !ok {
		e = &payload{}
		mm.mp[key] = e
	}
	e.cnt++
	e.mtx.Lock()
	mm.mtx.Unlock()
}

//Unlock try to unlock the (key, value). Use this after get the lock!
func (mm *MutexMap) Unlock(key interface{}) {
	mm.mtx.Lock()
	e, _ := mm.mp[key]
	e.cnt--
	if e.cnt < 1 {
		delete(mm.mp, key)
	}
	mm.mtx.Unlock()
	e.mtx.Unlock()
}
