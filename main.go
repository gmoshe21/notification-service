package main

import (
	"log"
	_ "notification/conn"
	"notification/routing"
	"time"
)

func main() {
	log.Println(time.Now())
	go routing.StartHttp()
	routing.CheckNotification()
}