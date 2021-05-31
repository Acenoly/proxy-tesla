package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"strings"
	"tesla/config"
	"tesla/globalvar"
	svc "tesla/service"
	"tesla/utils"
)

type TrafficParam struct {
	Username   string `json:"username"`
	ServerAddr string `json:"server_addr"`
	ClientAddr string `json:"client_addr"`
	TargetAddr string `json:"target_addr"`
	Bytes      string `json:"bytes"`
}

type KickParam struct {
	User string `form:"user" json:"user"`
	Ip string `form:"ip" json:"ip"`
}

func KickController(c *gin.Context){
	var info KickParam
	err := c.Bind(&info)
	if err != nil{
		fmt.Println(err.Error())
	}
	if info.User != ""{
		users := strings.Split(info.User, ",")
		return_users_str := ""
		for i := 0; i < len(users); i++ {
			infos := strings.Split(users[i], "-")
			user_username := infos[0]
			key := user_username
			value, err := utils.GetRedisValueByPrefix(key)
			if err == redis.Nil {
				utils.Log.WithField("key", key).Error("redis cache value is null")
				continue
			}
			//redis get value success
			res := strings.Split(value, ":")
			//用多了
			total, _ := strconv.ParseFloat(res[2], 8)
			used, _ := strconv.ParseFloat(res[3], 8)
			if used > total {
				return_users_str += users[i] + ","
			}
		}
		sz := len(return_users_str)
		if sz > 0{
			c.JSON(http.StatusOK, gin.H{
				"user": return_users_str[:sz-1],
				"ip":"",
			})
			return
		}else{
			c.JSON(http.StatusOK, gin.H{
				"user": "",
				"ip":"",
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"user": "",
		"ip":"",
	})
}

func AuthController(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("pass")
	client_addr := c.Query("client_addr")
	local_addr := c.Query("local_addr")
	//service := c.Query("service")
	//sps := c.Query("sps")
	target := c.Query("target")
	//fmt.Println(user, password, client_addr, service, sps, target)

	//fmt.Println(user_username, user_password, country, level, session, itype, rate)
	//key拼接token

	key := user
	value, err := utils.GetRedisValueByPrefix(key)
	//redis value is not found
	if err == redis.Nil {
		utils.Log.WithField("key", key).Error("Not this provider")
		c.JSON(http.StatusCreated, "redis cache value is null, redis key is  "+key)
		return
	}
	//redis server error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "redis server is not available")
		return
	}

	//redis get value success
	res := strings.Split(value, ":")
	//密码不正确
	if password != res[0] {
		utils.Log.WithField("password", res[0]).Error("password is not right")
		c.JSON(http.StatusCreated, "password is not right")
		return
	}

	client_ip := strings.Split(client_addr, ":")[0]
	//ip wrong
	if client_ip != res[1]{
		utils.Log.WithField("ip", res[0]).Error("ip is not right")
		c.JSON(http.StatusCreated, "password is not right")
		return
	}

	//用多了
	total, _ := strconv.ParseInt(res[2],10, 64)
	used, _ := strconv.ParseInt(res[3],10, 64)
        if used > total {
		c.JSON(http.StatusCreated, "current traffic is oversize")
		return
	}

	//优化版本
	key = user+local_addr
	v := globalvar.GetSession(key)
	if v != "None"{
		value = v
	}else{
		value, err = utils.GetRedisSession(key)
		globalvar.SetSession(key, value)
	}

	if err == redis.Nil {
		utils.Log.WithField("local", key).Error("Not this provider")
		c.JSON(http.StatusCreated, "redis cache value is null, redis key is  "+key)
		return
	}
	//redis server error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "redis server is not available")
		return
	}

	weblock := globalvar.GetWeblock()
	if !weblock{
		if  strings.Contains(strings.ToLower(target), "footlocker") ||
			strings.Contains(strings.ToLower(target), "champssports") ||
			strings.Contains(strings.ToLower(target), "footaction") ||
			strings.Contains(strings.ToLower(target), "eastbay"){
			value = ""
		}
	}

	c.Header("userconns", config.AppConfig.UserConns)
	c.Header("ipconns", config.AppConfig.IPConns)
	c.Header("userrate", "1000000")
	c.Header("iprate", "1000000")
	c.Header("upstream", value)
	c.JSON(http.StatusNoContent, "")
}

func TrafficController(c *gin.Context) {
	//server_addr := c.Query("server_addr")
	//client_addr := c.Query("client_addr")
	target_addr := c.Query("target_addr")
	username := c.Query("username")
	bytes := c.Query("bytes")

	//计算
	byteUse, err := strconv.ParseInt(bytes,10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "res[3] cannot parse to float, res[3] is "+ bytes)
		return
	}
	USERVALUE := globalvar.GETRUNARRAY()
	USERVALUE.Deposit(username, byteUse)

	//记录请求网络
	go func() {
		flag := utils.GetSneakerMap(target_addr)
		if !flag{
			message := username + "@"+ target_addr
			//push to kafka
			err := svc.PushWebLogParamToKafka(message)
			if err != nil {
				utils.Log.WithField("err", err).Error("push to kafka err")
				return
			}
		}
	}()

	//上传
	if globalvar.AddCOUNT() > 1000 {
		UploadToKafka()
	}
	c.JSON(http.StatusNoContent, "success")
}

func UploadWebLock(){
	value, err := utils.GetRedisValueByPrefix("lock")
	//redis value is not found
	if err == redis.Nil {
		globalvar.SetWeblock(true)
		return
	}
	//redis server error
	if err != nil {
		globalvar.SetWeblock(true)
		return
	}

	if value == "pass"{
		globalvar.SetWeblock(false)
	}else {
		globalvar.SetWeblock(true)
	}
}

func UploadToKafka(){
	USERVALUE := globalvar.GETRUNARRAY()
	userArray := USERVALUE.Content()
	globalvar.ClearCount()
	go func() {
		message := ""
		for key, value := range userArray {
			message += key + ":"+ fmt.Sprintf("%d", value) + ","
		}
		//push to kafka
		err := svc.PushTrafficParamToKafka(message)
		if err != nil {
			utils.Log.WithField("err", err).Error("push to kafka err")
			return
		}
	}()
}
