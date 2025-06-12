package lark_bot

import (
	"flag"
	"github.com/tengskyline/lark-bot/conf"
	"github.com/tengskyline/lark-bot/lark"
	"log"
)

func main() {
	configFile := flag.String("c", "conf/server.yaml", "default conf/server.yaml")
	flag.Parse()
	err := conf.ConfigInit(*configFile)
	if err != nil {
		log.Fatalf("config init failed, err %v", err)
	}

	eventHandler := lark.NewLarkHandler()
	app := lark.NewLark(eventHandler, conf.GlobalConfig)
	eventHandler.Bot = app
	err = app.Start()
	log.Println(err)
}
