package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"local.wayl/config"
	"local.wayl/handler"
)

var (
	conf config.Config
	path string
)

func init() {
	var err error
	path, err = os.UserHomeDir()
	path += "/.config/WAYL"
	if err != nil {
		log.Fatal(err)
	}

	config.InitConf(&path)

	if _, err := toml.DecodeFile(path+"/config.toml", &conf); err != nil {
		log.Fatal(err)
	}
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/Website"))))

	http.HandleFunc("/", handler.Root(&path))
	http.HandleFunc("/login", handler.Login(&conf))
	http.HandleFunc("/callback", handler.Callback)
	http.HandleFunc("/playback", handler.Playback)
	http.HandleFunc("/get-playback-data", handler.HandleGetPlaybackData)

	fmt.Println("Server started at http://localhost" + conf.Port)

	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
