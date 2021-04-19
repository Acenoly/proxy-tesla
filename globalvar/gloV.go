package globalvar

import "sync"

type RUNARRAY struct{
	sync.RWMutex
	m map[string]int64
}

var COUNT int
var USERARRAY = &RUNARRAY{}

func InitGlov(){
	//清空加锁
	USERARRAY = &RUNARRAY{m: make(map[string]int64)}
	COUNT = 0
}

func GETRUNARRAY() *RUNARRAY{
	return USERARRAY
}

func (b *RUNARRAY) GETRUNARRAYVALUE(key string)  int64{
	b.RLock()
	temp := b.m[key]
	b.RUnlock()
	return temp
}

func (b *RUNARRAY) Deposit(key string, used int64) {
	b.Lock()
	b.m[key] += used
	b.Unlock()
}

func (b *RUNARRAY) Content() map[string]int64 {
	b.Lock()
	temp := b.m
	b.m = make(map[string]int64)
	b.Unlock()
	return temp
}

func ClearCount(){
	COUNT = 0
}

func AddCOUNT() int{
	COUNT += 1
	return COUNT
}

