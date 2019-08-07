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
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6061", nil))
	}()

	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Printf("Config file not found.\nPlease press enter.")
		return
	}
	var config bot.Config
	json.Unmarshal(file, &config)

	// setting up bot
	client, err := bot.NewBot(&config)
	if err != nil {
		panic(err)
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
