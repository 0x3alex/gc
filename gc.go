package main

import (
	"time"
)

const (
	set = iota
	get
	del
)

type (
	gCObject struct {
		value        any
		neverExpires bool
		ttl          int64
	}

	GC struct {
		data       map[string]*gCObject
		DefaultTTL time.Duration
		buffer     int
		dataChan   chan gcRequestObject
		terminated chan bool
	}

	gcRequestObject struct {
		action       int
		key          string
		data         gCObject
		responseChan chan *gCObject
	}
)

func New(defaultTTL time.Duration, buffer int) *GC {
	return &GC{
		data:       make(map[string]*gCObject),
		DefaultTTL: defaultTTL,
		buffer:     buffer,
	}
}

func (gc *GC) SetNoTTL(key string, value any) {
	gc.dataChan <- gcRequestObject{
		action: set,
		key:    key,
		data: gCObject{
			value:        value,
			neverExpires: true,
		},
	}
}

func (gc *GC) Set(key string, value any, ttl time.Duration) {
	gc.dataChan <- gcRequestObject{
		action: set,
		key:    key,
		data: gCObject{
			value: value,
			ttl:   time.Now().Add(ttl).Unix(),
		},
	}
}

func (gc *GC) set(gcRequestObject gcRequestObject) {
	gc.data[gcRequestObject.key] = &gcRequestObject.data
}

func (gc *GC) Get(key string) (bool, any) {
	responseChan := make(chan *gCObject)
	gc.dataChan <- gcRequestObject{
		action:       get,
		key:          key,
		responseChan: responseChan,
	}
	m := <-responseChan
	close(responseChan)
	if m == nil {
		return false, nil
	}
	return true, m.value
}

func (gc *GC) get(gcRequestObject gcRequestObject) {
	gcRequestObject.responseChan <- gc.data[gcRequestObject.key]
}

func (gc *GC) Delete(key string) {
	gc.dataChan <- gcRequestObject{
		action: del,
		key:    key,
	}
}

func (gc *GC) delete(gcRequestObject gcRequestObject) {
	delete(gc.data, gcRequestObject.key)
}

func (gc *GC) Run() {
	gc.dataChan = make(chan gcRequestObject, gc.buffer)
	go func() {
		cleaner := time.NewTicker(gc.DefaultTTL)
		for {
			select {
			case <-cleaner.C:
				gc.clean()
			case <-gc.terminated:
				return
			case m := <-gc.dataChan:
				switch m.action {
				case set:
					gc.set(m)
				case get:
					gc.get(m)
				case del:
					gc.delete(m)
				}
			}
		}
	}()
}

func (gc *GC) StopCache() {
	gc.terminated <- true
	close(gc.dataChan)
}

func (gc *GC) clean() {
	for k, v := range gc.data {
		if !v.neverExpires && time.Now().Unix()-v.ttl >= 0 {
			delete(gc.data, k)
		}
	}
}
