package globalvar

import "sync"

type RUNARRAY struct{
	sync.RWMutex
	m map[string]int64
}

type SESSION struct {
	m map[string]string
}

var COUNT int
var USERARRAY = &RUNARRAY{}
var CACHESESSION = &SESSION{}

func InitGlov(){
	//清空加锁
	USERARRAY = &RUNARRAY{m: make(map[string]int64)}
	COUNT = 0
	CACHESESSION = &SESSION{m: make(map[string]string)}
}

func GETCACHESESSION() *SESSION{
	return CACHESESSION
}

func (b *SESSION) RemoveSession(){
	b.m = make(map[string]string)
}

func (b *SESSION) GetSession(key string) string{
	if val, ok := b.m[key]; ok {
		return val
	}
	return "None"
}

func (b *SESSION) SetSession(key string, value string){
	b.m[key] = value
}

func (b *SESSION)  GetWeblock() string{
	if val, ok := b.m["lock"]; ok {
		return val
	}
	return "unpass"
}

func (b *SESSION) SetWeblock(lock string){
	b.m["lock"] = lock
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

