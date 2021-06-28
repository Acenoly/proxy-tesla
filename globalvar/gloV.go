package globalvar

import "sync"

type RUNARRAY struct{
	sync.RWMutex
	m map[string]int64
}

type SESSION struct {
	sync.RWMutex
	m map[string]string
}
type SESSIONTEMP struct {
	sync.RWMutex
	m map[string]string
}

var COUNT int
var USERARRAY = &RUNARRAY{}
var CACHESESSION = &SESSION{}
var CACHETEMPSESSION = &SESSIONTEMP{}

func InitGlov(){
	//清空加锁
	USERARRAY = &RUNARRAY{m: make(map[string]int64)}
	COUNT = 0
	CACHESESSION = &SESSION{m: make(map[string]string)}
	CACHETEMPSESSION = &SESSIONTEMP{m: make(map[string]string)}
}

func GETCACHESESSION() *SESSION{
	return CACHESESSION
}

func GETCACHETEMPSESSION() *SESSIONTEMP{
	return CACHETEMPSESSION
}

func (b *SESSION) RemoveSession(){
	b.Lock()
	b.m = make(map[string]string)
	b.Unlock()
}

func (b *SESSIONTEMP) RemoveTempSession(){
	b.Lock()
	b.m = make(map[string]string)
	b.Unlock()
}

func (b *SESSION) GetSession(key string) string{
	temp := "None"
	b.RLock()
	if val, ok := b.m[key]; ok {
		temp = val
	}
	b.RUnlock()
	return temp
}

func (b *SESSIONTEMP) GetTempSession(key string) string{
	temp := "None"
	b.RLock()
	if val, ok := b.m[key]; ok {
		temp = val
	}
	b.RUnlock()
	return temp
}

func (b *SESSION) SetSession(key string, value string){
	b.Lock()
	b.m[key] = value
	b.Unlock()
}

func (b *SESSIONTEMP) SetTempSession(key string, value string){
	b.Lock()
	b.m[key] = value
	b.Unlock()
}

func (b *SESSION)  GetWeblock() string{
	temp := "unpass"
	b.RLock()
	if val, ok := b.m["lock"]; ok {
		temp = val
	}
	b.RUnlock()
	return temp
}

func (b *SESSION) SetWeblock(lock string){
	b.Lock()
	b.m["lock"] = lock
	b.Unlock()
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

