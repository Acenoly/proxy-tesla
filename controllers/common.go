package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"tesla/config"
	"tesla/globalvar"
	svc "tesla/service"
	"tesla/utils"
	"time"
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
			key := "userAuthOf" + user_username
			value, err := utils.GetRedisValueByPrefix(key)
			if err == redis.Nil {
				utils.Log.WithField("key", key).Error("redis cache value is null")
				continue
			}
			//redis get value success
			res := strings.Split(value, ":")
			//用多了
			total, _ := strconv.ParseFloat(res[1], 8)
			used, _ := strconv.ParseFloat(res[2], 8)
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
	//client_addr := c.Query("client_addr")
	local_addr := c.Query("local_addr")
	//service := c.Query("service")
	//sps := c.Query("sps")
	target := c.Query("target")
	//fmt.Println(user, password, client_addr, service, sps, target)
	flag := utils.GetSneakerMap(target)
	infos := strings.Split(user, "-")
	if len(infos) < 2{
		utils.Log.WithField("account", infos).Error("account err")
		c.JSON(http.StatusCreated, "account err")
		return
	}

	local_addrs := strings.Split(local_addr, ":")
	port := local_addrs[1]
	//fmt.Println(user_username, user_password, country, level, session, itype, rate)
	//key拼接token
	user_username := ""
	//user_password := password
	country := ""
	session := ""
	itype := ""
	rate := ""

	if local_addrs[1] == "15000"{
		user_username = infos[0]
		country = infos[1]
		session = infos[3]
		itype = infos[4]
		rate = infos[5]
	}else if local_addrs[1] == "15001"{
		user_username = infos[0]
		country = infos[1]
		session = infos[3]
		itype = infos[4]
		rate = "0"
	}else{
		user_username = infos[0]
		country = infos[2]
		itype = infos[4]
		session = infos[6]
		rate = "0"
	}
	fmt.Print("here5")

	key := port + user_username
	session_number, err := strconv.Atoi(session)
	session_number = session_number
	if err != nil {
		utils.Log.WithField("session", session).Error("session parse to int err")
		c.JSON(http.StatusCreated, "session parse to int err")
		return
	}

	value, err := utils.GetRedisValueByPrefix(key)
	//redis value is not found
	if err == redis.Nil {
		utils.Log.WithField("key", key).Error("Not this provider")
		c.JSON(http.StatusCreated, "redis cache value is null, redis key is  "+key)
		return
	}
	fmt.Print("here3")

	//redis server error
	if err != nil {
		c.JSON(http.StatusInternalServerError, "redis server is not available")
		return
	}

	key = "userAuthOf" + user_username
	value, err = utils.GetRedisValueByPrefix(key)
	fmt.Print("here2")

	//redis get value success
	res := strings.Split(value, ":")
	//密码不正确
	if password != res[0] {
		utils.Log.WithField("password", res[0]).Error("password is not right")
		c.JSON(http.StatusCreated, "password is not right")
		return
	}

	//用多了
	total, _ := strconv.ParseInt(res[1],10, 64)
	used, _ := strconv.ParseInt(res[2],10, 64)	
        if used > total {
		c.JSON(http.StatusCreated, "current traffic is oversize")
		return
	}

	//优化版本
	t := ""
	key = "session:" + user_username + session
	val, err := utils.GetRedisValueByPrefix(key)
	//redis value is not found
	fmt.Print("here1")

	//redis server error
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, "redis server is not available")
		return
	}

	// value is nil
	if err == redis.Nil {
		key =  "AccountInfo" + res[3]
		fmt.Print("here10")
		val, err := utils.GetRedisValueByPrefix(key)
		if err == redis.Nil {
			c.JSON(http.StatusCreated, "redis value is nil , key is "+key)
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, "redis server is not available")
			return
		}
		accounts_value := strings.Split(val, ":")
		totalNumber, err := strconv.Atoi(accounts_value[0])
		if err != nil {
			c.JSON(http.StatusInternalServerError, "accounts_value[0] can not parse to int, accounts_value[0] is"+accounts_value[0])
			return
		}
		if totalNumber == 0 {
			c.Header("userconns", config.AppConfig.UserConns)
			c.Header("ipconns", config.AppConfig.IPConns)
			c.Header("userrate", rate)
			c.Header("iprate", rate)
			c.Header("upstream", "")
			c.JSON(http.StatusNoContent, "success")
		}

		pick := session_number % totalNumber
		accounts_info := accounts_value[pick+1]
		accounts_array := strings.Split(accounts_info, "-")
		if accounts_array[0] == "geo"{
			t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
		}
		if accounts_array[0] == "lumi" {
			if itype == "Rotate" || country == "usf" || !flag{
				accounts_info := accounts_value[1]
				accounts_array := strings.Split(accounts_info, "-")
				t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
			} else{
				t = svc.CreateLumi(accounts_array[3], session, country, accounts_array[1], accounts_array[2])
			}
		}
		if accounts_array[0] == "oxy" {
			if itype == "Rotate" || country == "usf" || country == "mo" || country == "cn"  || country == "hk" || country == "cz"  {
				accounts_info := accounts_value[1]
				accounts_array := strings.Split(accounts_info, "-")
				t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
			}else{
				rand.Seed(time.Now().UnixNano())
				number := rand.Intn(3)
				if number != 1{
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				} else{
					t = svc.CreateOneOxy(country, itype, session, accounts_array[1], accounts_array[2])
				}
			}
		}
		if accounts_array[0] == "smart" {
			if itype == "Rotate" || country == "usf" || country == "mo" || country == "cn"  || country == "hk" || country == "cz"  {
				accounts_info := accounts_value[1]
				accounts_array := strings.Split(accounts_info, "-")
				t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
			}else{
				t = svc.CreateOneSmart(country, itype, session, accounts_array[1], accounts_array[2])
			}
		}
		redis_key := "session:" + user_username + session
		fmt.Print("here")
		err = utils.SetRedisValueByPrefix(redis_key, t, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "redis set value error key is "+redis_key+", value is "+t)
			return
		}
	} else {
		t = val
	}

	c.Header("userconns", config.AppConfig.UserConns)
	c.Header("ipconns", config.AppConfig.IPConns)
	c.Header("userrate", rate)
	c.Header("iprate", rate)
	c.Header("upstream", t)
	c.JSON(http.StatusNoContent, "success")

}

func TrafficController(c *gin.Context) {
	//server_addr := c.Query("server_addr")
	//client_addr := c.Query("client_addr")
	//target_addr := c.Query("target_addr")
	username := c.Query("username")
	bytes := c.Query("bytes")

	//这里是拿Key
	infos := strings.Split(username, "-")
	user_username := infos[0]
	userkey := ""
	userkey = "userAuthOf" + user_username

	//计算
	byteUse, err := strconv.ParseInt(bytes,10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "res[2] cannot parse to float, res[2] is "+ bytes)
		return
	}
	USERVALUE := globalvar.GETRUNARRAY()
	USERVALUE.Deposit(userkey, byteUse)
	//上传
	if globalvar.AddCOUNT() > 1000 {
		UploadToKafka()
	}
	c.JSON(http.StatusNoContent, "success")
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
