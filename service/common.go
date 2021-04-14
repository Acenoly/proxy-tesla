package service

import (
	"github.com/Shopify/sarama"
	"math/rand"
	"strconv"
	"strings"
	"tesla/config"
	"tesla/utils"
	"time"
)

var (
	support_country = []string{"ca", "us", "usf", "au", "sg", "gb", "de", "it", "es", "se", "hu", "dk", "cz", "pl", "nl", "mo", "hk", "cn"}
	geo_sticky      = map[string]string{
		"usf": "us-30m.geosurf.io:8000",
		"us":  "us-30m.geosurf.io:8000",
		"au":  "au-30m.geosurf.io:8000",
		"sg":  "ctr-2-30m.geosurf.io:8000",
		"gb":  "gb-30m.geosurf.io:8000",
		"de":  "de-30m.geosurf.io:8000",
		"it":  "it-30m.geosurf.io:8000",
		"es":  "es-30m.geosurf.io:8000",
		"se":  "ctr-2-30m.geosurf.io:8000",
		"hu":  "ctr-2-30m.geosurf.io:8000",
		"dk":  "ctr-2-30m.geosurf.io:8000",
		"cz":  "cz-30m.geosurf.io:8000",
		"pl":  "pl-30m.geosurf.io:8000",
		"nl":  "ctr-2-30m.geosurf.io:8000",
		"mo":  "ctr-2-30m.geosurf.io:8000",
		"hk":  "hk-30m.geosurf.io:8000",
		"cn":  "cn-30m.geosurf.io:8000",
		"ca":  "ca-30m.geosurf.io:8000",
	}
	geo_sticky_10      = map[string]string{
		"usf": "us-10m.geosurf.io:8000",
		"us":  "us-10m.geosurf.io:8000",
		"au":  "au-10m.geosurf.io:8000",
		"sg":  "ctr-2-10m.geosurf.io:8000",
		"gb":  "gb-10m.geosurf.io:8000",
		"de":  "de-10m.geosurf.io:8000",
		"it":  "it-10m.geosurf.io:8000",
		"es":  "es-10m.geosurf.io:8000",
		"se":  "ctr-2-10m.geosurf.io:8000",
		"hu":  "ctr-2-10m.geosurf.io:8000",
		"dk":  "ctr-2-10m.geosurf.io:8000",
		"cz":  "cz-10m.geosurf.io:8000",
		"pl":  "pl-10m.geosurf.io:8000",
		"nl":  "ctr-2-10m.geosurf.io:8000",
		"mo":  "ctr-2-10m.geosurf.io:8000",
		"hk":  "hk-10m.geosurf.io:8000",
		"cn":  "cn-10m.geosurf.io:8000",
		"ca":  "ca-10m.geosurf.io:8000",
	}

	geo_rotate = map[string]string{
		"usf": "us-1m.geosurf.io:8000",
		"us":  "us-1m.geosurf.io:8000",
		"au":  "au-1m.geosurf.io:8000",
		"sg":  "ctr-2-1m.geosurf.io:8000",
		"gb":  "gb-1m.geosurf.io:8000",
		"de":  "de-1m.geosurf.io:8000",
		"it":  "it-1m.geosurf.io:8000",
		"es":  "es-1m.geosurf.io:8000",
		"se":  "ctr-2-1m.geosurf.io:8000",
		"hu":  "ctr-2-1m.geosurf.io:8000",
		"dk":  "ctr-2-1m.geosurf.io:8000",
		"cz":  "cz-1m.geosurf.io:8000",
		"pl":  "pl-1m.geosurf.io:8000",
		"nl":  "ctr-2-1m.geosurf.io:8000",
		"mo":  "ctr-2-1m.geosurf.io:8000",
		"hk":  "hk-1m.geosurf.io:8000",
		"cn":  "cn-1m.geosurf.io:8000",
	}

	geoEu = []string{"gb","de","it","cz","pl","nl","hu","se","es"}
	lumiEu = []string{"UK", "AL", "AD", "AT", "BY", "BE", "BA", "BG", "HR", "DK", "EE", "FI", "FR", "DE", "GR", "HU", "IS", "IE", "IT", "LV", "LI", "LT", "LU", "MT", "MC", "ME", "NL", "NO", "PL", "PT", "RO", "RS", "SK", "SI", "ES", "SE", "CH", "UA"}

	oxyEuRotate = []string{"gb-pr.oxylabs.io:20000", "de-pr.oxylabs.io:30000", "fr-pr.oxylabs.io:40000", "es-pr.oxylabs.io:10000", "it-pr.oxylabs.io:20000", "se-pr.oxylabs.io:30000", "gr-pr.oxylabs.io:40000", "pt-pr.oxylabs.io:10000", "nl-pr.oxylabs.io:20000", "be-pr.oxylabs.io:30000", "ua-pr.oxylabs.io:10000", "pl-pr.oxylabs.io:20000", "dk-pr.oxylabs.io:19000", "al-pr.oxylabs.io:49000", "ad-pr.oxylabs.io:10000", "at-pr.oxylabs.io:11000", "ba-pr.oxylabs.io:13000", "bg-pr.oxylabs.io:14000", "by-pr.oxylabs.io:15000", "hr-pr.oxylabs.io:16000", "dk-pr.oxylabs.io:19000", "ee-pr.oxylabs.io:20000", "fi-pr.oxylabs.io:21000", "hu-pr.oxylabs.io:23000", "is-pr.oxylabs.io:24000", "ie-pr.oxylabs.io:25000", "lv-pr.oxylabs.io:26000", "li-pr.oxylabs.io:27000", "lt-pr.oxylabs.io:28000", "lu-pr.oxylabs.io:29000", "mt-pr.oxylabs.io:30000", "mc-pr.oxylabs.io:31000", "me-pr.oxylabs.io:33000", "no-pr.oxylabs.io:34000", "ro-pr.oxylabs.io:35000", "rs-pr.oxylabs.io:36000", "sk-pr.oxylabs.io:37000", "si-pr.oxylabs.io:38000", "ch-pr.oxylabs.io:39000"}
	oxyEuSticky = []string{"GB", "DE", "FR", "ES", "IT", "SE", "GR", "PT", "NL", "BE", "UA", "PL", "DK", "AL", "AD", "AT", "BA", "BG", "BY", "HR", "DK", "EE", "FI", "HU", "IS", "IE", "LV", "LI", "LT", "LU", "MT", "MC", "ME", "NO", "RO", "RS", "SK", "SI", "CH"}

	smartEuRotate = []string{"gb.smartproxy.com:30000", "de.smartproxy.com:20000", "fr.smartproxy.com:40000", "es.smartproxy.com:10000", "it.smartproxy.com:20000", "se.smartproxy.com:20000", "gr.smartproxy.com:30000", "pt.smartproxy.com:20000", "nl.smartproxy.com:10000", "be.smartproxy.com:40000", "ua.smartproxy.com:40000", "pl.smartproxy.com:20000", "dk.smartproxy.com:27000", "al.smartproxy.com:33000", "ad.smartproxy.com:34000", "at.smartproxy.com:35000", "ba.smartproxy.com:37000", "bg.smartproxy.com:38000", "by.smartproxy.com:39000", "hr.smartproxy.com:40000", "dk.smartproxy.com:27000", "ee.smartproxy.com:28000", "fi.smartproxy.com:41000", "hu.smartproxy.com:43000", "is.smartproxy.com:23000", "ie.smartproxy.com:24000", "lv.smartproxy.com:22000", "li.smartproxy.com:23000", "lt.smartproxy.com:24000", "lu.smartproxy.com:25000", "mt.smartproxy.com:49000", "mc.smartproxy.com:10000", "me.smartproxy.com:12000", "no.smartproxy.com:13000", "ro.smartproxy.com:13000", "rs.smartproxy.com:14000", "sk.smartproxy.com:15000", "si.smartproxy.com:16000", "ch.smartproxy.com:29000"}
	smartEuSticky = []string{"gb", "de", "fr", "es", "it", "se", "gr", "pt", "nl", "be", "ua", "pl", "dk", "al", "ad", "at", "ba", "bg", "by", "hr", "dk", "ee", "fi", "hu", "is", "ie", "lv", "li", "lt", "lu", "mt", "mc", "me", "no", "ro", "rs", "sk", "si", "ch"}
)

func isCountrySupport(country string) (flag bool) {
	for _, v := range support_country {
		if country == v {
			flag = true
			return
		}
	}
	return
}

func CreateOneGeo(country, types, session, username, password string) string {
	session_int, _ := strconv.Atoi(session)
    if country == "usf" {
        country = "us"
    }
    if country == "eu" {
        country = geoEu[session_int%len(geoEu)]
    }
	s := username + "+" + strings.ToUpper(country) + "+" + username + "-" + session + ":" + password + "@"
	if strings.ToLower(types) == "sticky" {
		if isCountrySupport(country) {
			rand.Seed(time.Now().UnixNano())
			number := rand.Intn(3)
			if number == 1{
				s += geo_sticky[country]
			}else{
				s += geo_sticky_10[country]
			}
		}
	} else {
		if isCountrySupport(country) {
			s += geo_rotate[country]
		}
	}
	return s
}

func CreateLumi(zone, session, country, username, password string) string {
	ip_country := ""
	if country == "eu" {
		session_int, _ := strconv.Atoi(session)
		ip_country = "lum-customer-" + username + "-zone-" + zone + "-country-" + strings.ToLower(lumiEu[session_int%len(lumiEu)]) + "-session-" + session + ":" + password + "@zproxy.lum-superproxy.io:22225"
	} else {
		ip_country = "lum-customer-" + username + "-zone-" + zone + "-country-" + country + "-session-" + session + ":" + password + "@zproxy.lum-superproxy.io:22225"
	}
	return ip_country
}

func CreateOneOxy(country, types, session, username, password string) (s string) {
	session_int, _ := strconv.Atoi(session)
	if strings.ToLower(types) == "sticky" {
		if country != "us" {
			s = "user-" + username + "-country-" + strings.ToUpper(country) + "-session-" + session + ":" + password + "@" + "pr.oxylabs.io:7777"
		}else{
			s = "user-" + username + "-country-US-session-" + session + ":" + password + "@" + "pr.oxylabs.io:7777"
		}
	} else {
		if country != "us" {
			s = username + ":" + password + "@" + oxyEuRotate[session_int%len(oxyEuRotate)]
		} else {
			s = username + ":" + password + "@us-pr.oxylabs.io:10000"
		}
	}
	return
}

func CreateOneSmart(country, types, session, username, password string) (s string) {
	session_int, _ := strconv.Atoi(session)
	if strings.ToLower(types) == "sticky" {
		if country != "us" {
			s = "user-" + username + "-country-" + country + "-session-" + session + ":" + password + "@" + "gate.smartproxy.com:7000"
		} else {
			s = "user-" + username + "-country-" + "us" + "-session-" + session + ":" + password + "@" + "gate.smartproxy.com:7000"
		}
	} else {
		if country != "us" {
			s = username + ":" + password + "@" + smartEuRotate[session_int%len(smartEuRotate)]
		} else {
			s = username + ":" + password + "@us.smartproxy.com:10000"
		}
	}
	return
}

func PushTrafficParamToKafka(s string) error {
	//utils.Log.Infof("kafka value is %s", s)
	msg := &sarama.ProducerMessage{
		Topic: config.AppConfig.Topic,
		Value: sarama.StringEncoder(s),
	}
	//partition, offset, err := producer.SendMessage(msg)
	_, _, err := utils.Producer.SendMessage(msg)
	return err
}
