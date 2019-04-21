package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/intrntsrfr/meido/bot"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	jeff := zap.NewDevelopmentConfig()
	jeff.OutputPaths = []string{"./logs.txt"}
	jeff.ErrorOutputPaths = []string{"./logs.txt"}
	logger, err := jeff.Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer logger.Sync()
	logger.Info("Logger construction succeeded")

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Config file not found.\nPlease press enter.")
		return
	}
	var config bot.Config
	json.Unmarshal(file, &config)

	client, err := bot.NewBot(&config, logger.Named("discord"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	err = client.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
