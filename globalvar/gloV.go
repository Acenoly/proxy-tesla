package globalvar

import "sync"

type RUNARRAY struct{
	sync.RWMutex
	m map[string]int64
}


var COUNT int
var USERARRAY = &RUNARRAY{}
var WebLock bool
var Session map[string]string

func InitGlov(){
	//清空加锁
	USERARRAY = &RUNARRAY{m: make(map[string]int64)}
	COUNT = 0
	WebLock = false
	Session = map[string]string{}
}

func GetSession(key string) string{
	if val, ok := Session[key]; ok {
		return val
	}
	return "None"
}

func SetSession(key string, value string){
	Session[key] = value
}

func GetWeblock() bool{
	return WebLock
}

func SetWeblock(weblock bool){
	WebLock = weblock
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

