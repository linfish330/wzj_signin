package main

import (
	"time"
	"wzj_signin/config"
	"wzj_signin/db"
	"wzj_signin/server"
	"wzj_signin/service"

	"github.com/spf13/viper"
)

func main() {
	if err := config.Load(); err != nil {
		panic(err)
	}
	db.InitRedis()
	go startTimer()
	server.Start()
}

func startTimer() {
	for {
		interval := viper.GetInt("app.interval")
		if interval < 1 {
			interval = 8
		}
		time.Sleep(time.Duration(interval) * time.Second)

		for _, openId := range db.RedisGetAllMatchedKeys("wzj:user:*") {
			openId := openId[9:]
			signList, _ := service.GetAllSigns(openId)
			for _, sign := range signList {
				go service.Signin(sign, openId)
			}
		}
	}
}
