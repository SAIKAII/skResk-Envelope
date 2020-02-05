package main

import (
	"time"

	"github.com/tietang/go-eureka-client/eureka"
)

func main() {
	cfg := eureka.Config{
		DialTimeout: time.Second * 10,
	}
	client := eureka.NewClientByConfig([]string{
		"http://127.0.0.1:8761/eureka",
	}, cfg)
	appName := "Go-Example"
	instance := eureka.NewInstanceInfo(
		"test.com", appName, "127.0.0.2",
		8083, 30, false)
	client.RegisterInstance(appName, instance)
	client.Start()
	c := make(chan int)
	<-c
}
