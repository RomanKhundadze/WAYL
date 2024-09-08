package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"local.wayl/config"
	"local.wayl/handler"
)

var conf config.Config

func init() {
	if _, err := toml.DecodeFile("config/config.toml", &conf); err != nil {
	}
}

func main() {
	http.HandleFunc("/", handler.Root)
	http.HandleFunc("/login", handler.Login(&conf))
	http.HandleFunc("/callback", handler.Callback)
	http.HandleFunc("/playback", handler.Playback)

	fmt.Println("Server started at http://localhost" + conf.Port)

	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
