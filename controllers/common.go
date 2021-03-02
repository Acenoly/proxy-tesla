package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"tesla/config"
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

func AuthController(c *gin.Context) {
	user := c.Query("user")
	password := c.Query("pass")
	//client_addr := c.Query("client_addr")
	//service := c.Query("service")
	//sps := c.Query("sps")
	//target := c.Query("target")
	//fmt.Println(user, password, client_addr, service, sps, target)

	infos := strings.Split(user, "-")
	user_username := infos[0]
	//user_password := password
	country := infos[1]
	level := infos[2]
	session := infos[3]
	itype := infos[4]
	rate := infos[5]

	//fmt.Println(user_username, user_password, country, level, session, itype, rate)

	key := ""
	if level == "basic" {
		key = "userBaseAuthOf" + user_username
	} else if level == "super" {
		key = "userSuperAuthOf" + user_username
	} else {
		utils.Log.WithField("level", level).Error("level is not basic or super")
		c.JSON(http.StatusCreated, "level is not basic or super")
		return
	}

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
		utils.Log.WithField("key", key).Error("redis cache value is null")
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

	//用多了
	total, _ := strconv.ParseFloat(res[1], 8)
	used, _ := strconv.ParseFloat(res[2], 8)
	if used > total {
		c.JSON(http.StatusCreated, "current traffic is oversize")
		return
	}

	//优化版本
	t := ""
	key = user_username + session
	val, err := utils.GetRedisValueByPrefix(key)
	//redis value is not found

	//redis server error
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, "redis server is not available")
		return
	}

	// value is nil
	if err == redis.Nil {
		if level == "basic" {
			key = "BasicAccountInfo" + user_username
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
				c.JSON(http.StatusCreated, "totalNumber is 0")
				return
			}
			pick := session_number % totalNumber
			accounts_info := accounts_value[pick+1]
			accounts_array := strings.Split(accounts_info, "-")
			if accounts_array[0] == "geo" {
				t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
			}
			if accounts_array[0] == "oxy" {
				rand.Seed(time.Now().UnixNano())
				number := rand.Intn(1)
				if number != 1{
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				} else if itype == "Rotate" || country == "usf" || country == "au" || country == "sg" || country == "mo" || country == "cn"  || country == "hk" || country == "cz"  {
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				} else{
					t = svc.CreateOneOxy(country, itype, session, accounts_array[1], accounts_array[2])
				}
			}
			if accounts_array[0] == "smart" {
				rand.Seed(time.Now().UnixNano())
				number := rand.Intn(5)
				if number != 1{
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				}else if itype == "Rotate" || country == "usf" || country == "au" || country == "sg" || country == "mo" || country == "cn"  || country == "hk" || country == "cz"  {
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				} else{
					t = svc.CreateOneSmart(country, itype, session, accounts_array[1], accounts_array[2])
				}
			}

		} else {
			key = "SuperAccountInfo" + user_username
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
				c.JSON(http.StatusCreated, "totalNumber is 0")
				return
			}
			pick := session_number % totalNumber
			accounts_info := accounts_value[pick+1]
			accounts_array := strings.Split(accounts_info, "-")
			if accounts_array[0] == "lumi" {
				rand.Seed(time.Now().UnixNano())
				number := rand.Intn(10)
				if number != 1{
					accounts_info := accounts_value[1]
					accounts_array := strings.Split(accounts_info, "-")
					t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
				} else{
					t = svc.CreateLumi(accounts_array[3], session, country, accounts_array[1], accounts_array[2])
				}
			}
			if accounts_array[0] == "geo" {
				t = svc.CreateOneGeo(country, itype, session, accounts_array[1], accounts_array[2])
			}

			if accounts_array[0] == "oxy" {
				t = svc.CreateOneOxy(country, itype, session, accounts_array[1], accounts_array[2])

			}

			if accounts_array[0] == "smart" {
				t = svc.CreateOneSmart(country, itype, session, accounts_array[1], accounts_array[2])
			}
		}
		redis_key := user_username + session
		err = utils.SetRedisValueByPrefix(redis_key, t, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "redis set value error key is "+redis_key+", value is "+t)
			return
		}
	} else {
		t = val
	}

	//t := ""
	//if level == "basic" {
	//	key := "geoAccountNumber"
	//	val, err := utils.GetRedisValueByPrefix(key)
	//	//redis value is not found
	//
	//	//redis server error
	//	if err != nil && err != redis.Nil {
	//		c.JSON(http.StatusInternalServerError, "redis server is not available")
	//		return
	//	}
	//
	//	// value is nil
	//	if err == redis.Nil {
	//		t = svc.CreateOneGeo(country, itype, session, "abc", "123")
	//	} else {
	//		number, _ := strconv.Atoi(val)
	//		pick := session_number % number
	//		oxy_key := "geoAccount_" + strconv.Itoa(pick)
	//		oxy_val, err := utils.GetRedisValueByPrefix(oxy_key)
	//		if err != nil && err != redis.Nil {
	//			c.JSON(http.StatusInternalServerError, "redis server is not available")
	//			return
	//		}
	//		account_info := strings.Split(oxy_val, ":")
	//		t = svc.CreateOneGeo(country, itype, session, account_info[0], account_info[1])
	//	}
	//} else {
	//	key := "lumiAccountNumber"
	//	val, err := utils.GetRedisValueByPrefix(key)
	//	//redis value is not found
	//
	//	//redis server error
	//	if err != nil && err != redis.Nil {
	//		c.JSON(http.StatusInternalServerError, "redis server is not available")
	//		return
	//	}
	//
	//	// value is nil
	//	if err == redis.Nil {
	//		t = svc.CreateLumi(country, itype, session, "abc", "123")
	//	} else {
	//		number, _ := strconv.Atoi(val)
	//		pick := session_number % number
	//		oxy_key := "lumiAccount_" + strconv.Itoa(pick)
	//		oxy_val, err := utils.GetRedisValueByPrefix(oxy_key)
	//		if err != nil && err != redis.Nil {
	//			c.JSON(http.StatusInternalServerError, "redis server is not available")
	//			return
	//		}
	//		account_info := strings.Split(oxy_val, ":")
	//		t = svc.CreateLumi(account_info[2], session, country, account_info[0], account_info[1])
	//	}
	//}
	//
	c.Header("userconns", config.AppConfig.UserConns)
	c.Header("ipconns", config.AppConfig.IPConns)
	c.Header("userrate", rate)
	c.Header("iprate", rate)
	c.Header("upstream", "http://"+t)
	c.JSON(http.StatusNoContent, "success")

}

func TrafficController(c *gin.Context) {
	//user := c.Query("username")
	server_addr := c.Query("server_addr")
	client_addr := c.Query("client_addr")
	target_addr := c.Query("target_addr")
	username := c.Query("username")
	bytes := c.Query("bytes")

	//infos := strings.Split(user, "-")
	//user_username := infos[0]
	//country := infos[1]
	//level := infos[2]
	//session := infos[3]
	//itype := infos[4]

	//key := ""
	//if level == "basic" {
	//	key = "userBaseAuthOf" + user_username
	//} else if level == "super" {
	//	key = "userSuperAuthOf" + user_username
	//} else {
	//	c.JSON(http.StatusCreated, "level is not basic or super")
	//	return
	//}
	//value, err := utils.GetRedisValueByPrefix(key)
	//redis value is not found
	//if err == redis.Nil {
	//	c.JSON(http.StatusCreated, "redis cache value is null, redis key is  "+key)
	//	return
	//}

	//redis server error
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, "redis server is not available")
	//	return
	//}

	//res := strings.Split(value, ":")
	//byteUse, err := strconv.ParseFloat(bytes, 8)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, "byteUse cannot parse to float, byteUse is "+bytes)
	//	return
	//}
	//if byteUse < 1000 {
	//	byteUse = 1000
	//}
	//byteUse = math.Ceil(byteUse/1000) * 1000
	//usage := float64(4) * byteUse / 1000 / 1000 / 10

	//res2Float, err := strconv.ParseFloat(res[2], 8)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, "res[2] cannot parse to float, res[2] is "+res[2])
	//	return
	//}
	//total := res2Float + usage
	//rspon, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", total), 64)

	//保存
	//rsponStr := strconv.FormatFloat(rspon, 'E', -1, 64) //float64
	//final := res[0] + ":" + res[1] + ":" + rsponStr
	//err = utils.SetRedisValueByPrefix(key, final, 0)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, "redis set kv error key is "+key+" , value is "+final)
	//	return
	//}

	params := TrafficParam{
		Username:   username,
		ServerAddr: server_addr,
		ClientAddr: client_addr,
		TargetAddr: target_addr,
		Bytes:      bytes,
	}
	paramByte, err := json.Marshal(&params)
	if err != nil {
		utils.Log.WithField("param", params).Error("TrafficParam struct marshal to json failed")
		c.JSON(http.StatusInternalServerError, "TrafficParam struct marshal to json failed")
		return
	}

	go func() {
		//push to kafka
		err = svc.PushTrafficParamToKafka(string(paramByte))
		if err != nil {
			utils.Log.WithField("err", err).Error("push to kafka err")
			return
		}
	}()

	c.JSON(http.StatusNoContent, "success")

}
