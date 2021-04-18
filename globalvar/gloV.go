package globalvar

import "sync"

type RUNARRAY struct{
	sync.RWMutex
	m map[string]float64
}

var COUNT int
var USERARRAY = RUNARRAY{}

func InitGlov(){
	//清空加锁
	USERARRAY.Lock()
	USERARRAY = RUNARRAY{m: make(map[string]float64)}
	USERARRAY.Unlock()
	COUNT = 0
}

func UpdateUSERARRAYVal(key string, used float64){
	//读数据锁
	USERARRAY.RLock()
	if val, ok := USERARRAY.m[key]; ok{
		//解了
		USERARRAY.RUnlock()
		//写数据锁
		USERARRAY.Lock()
		val += used
		USERARRAY.Unlock()
	}else{
		//没有也解了
		USERARRAY.RUnlock()
		//写数据锁
		USERARRAY.Lock()
		USERARRAY.m[key] = used
		USERARRAY.Unlock()
	}
}

func CopyMap() map[string]float64{
	targetMap  := make(map[string]float64)
	USERARRAY.RLock()
	for key, value := range USERARRAY.m {
		targetMap[key] = value
	}
	USERARRAY.Lock()
	USERARRAY.m = make(map[string]float64)

	//全打开
	USERARRAY.Unlock()
	USERARRAY.RUnlock()

	return targetMap
}


func AddCOUNT() int{
	COUNT += 1
	return COUNT
}

