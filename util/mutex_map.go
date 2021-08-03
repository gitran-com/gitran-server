//Refer to https://stackoverflow.com/a/62562831

package util

import (
	"fmt"
	"sync"
)

//MutexMap map multiple mutex using key
type MutexMap struct {
	mtx sync.Mutex
	mp  map[interface{}]*payload
}

type payload struct {
	mm  *MutexMap
	mtx sync.Mutex
	cnt uint
	key interface{}
}

// Unlocker provides an Unlock method to release the lock.
type Unlocker interface {
	Unlock()
}

//NewMutexMap create a new MutexMap
func NewMutexMap() *MutexMap {
	return &MutexMap{mp: make(map[interface{}]*payload)}
}

//Lock try to lock the (key, value)
func (mm *MutexMap) Lock(key interface{}) Unlocker {
	mm.mtx.Lock()
	// defer fmt.Printf("I'm locked!\n")
	e, ok := mm.mp[key]
	if !ok {
		e = &payload{mm: mm, key: key}
		mm.mp[key] = e
	}
	e.cnt++
	mm.mtx.Unlock() //顺序不能换
	e.mtx.Lock()
	return e
}

//Unlock try to unlock the (key, value). Use this after get the lock!
func (p *payload) Unlock() {
	mm := p.mm
	mm.mtx.Lock()
	// defer fmt.Printf("I'm unlocked!\n")
	e, ok := mm.mp[p.key]
	if !ok {
		panic(fmt.Errorf("Unlock requested for key=%v but no entry found", p.key))
	}
	e.cnt--
	if e.cnt < 1 {
		delete(mm.mp, p.key)
	}
	mm.mtx.Unlock()
	e.mtx.Unlock()

}
