package main

import (
	"flag"
	"fmt"
	"github.com/tengskyline/lark-bot/conf"
	"github.com/tengskyline/lark-bot/lark"
	"log"
)

func main() {
	configFile := flag.String("c", "conf/config.yaml", "default conf/config.yaml")
	flag.Parse()
	fmt.Printf("configFile %+v\n", *configFile)
	err := conf.ConfigInit(*configFile)
	if err != nil {
		fmt.Println("config init failed, err %v", err)
	}
	fmt.Printf("config init %+v", conf.GlobalConfig)
	eventHandler := lark.NewLarkHandler()
	app := lark.NewLark(eventHandler, conf.GlobalConfig)
	eventHandler.Bot = app
	err = app.Start()
	log.Println(err)
}
