package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"local.wayl/config"
	"local.wayl/handler"
	"local.wayl/manageRuntime"
	"log"
	"net/http"
	"os"
)

var (
	conf     config.Config
	path     string
	authPath string
)

func init() {
	newPath, err := os.UserHomeDir()
	newPath += "/.config/WAYL"
	if err != nil {
		log.Fatal("Error Creating var path:", err)
	}
	path = newPath
	authPath, err = os.UserHomeDir()
	authPath += "/.WAYL"
	if err != nil {
		log.Println("Error Creating var authPath:", err)
		return
	}

	err = config.InitConf(&path, &authPath)
	if err != nil {
		log.Fatal("Error Initialising Configs:", err)
		return
	}

	if _, err := toml.DecodeFile(path+"/config.toml", &conf); err != nil {
		log.Fatal("Error Decoding config.toml:", err)
	}
}

func processArgs() {
	if len(os.Args) < 2 {
		return
	}

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch arg {
		case "-k":
			if len(os.Args) != 2 {
				log.Fatal("Error: can only pass 1 argument when using -k")
			}
			manageRuntime.KillRunningInstances()
			os.Exit(0)

		default:
			log.Fatalf("Error: unknown argument '%s'", arg)
		}
	}
}

func main() {
	processArgs()

	fileServer := http.FileServer(http.Dir(path + "/Website"))
	http.Handle("/website/", http.StripPrefix("/website/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fileServer.ServeHTTP(w, r)
	})))

	http.HandleFunc("/", handler.Root(&path, &authPath))
	http.HandleFunc("/login", handler.Login(&conf))
	http.HandleFunc("/callback", handler.Callback)
	http.HandleFunc("/playback", handler.Playback)
	http.HandleFunc("/getPlaybackData", handler.HandleGetPlaybackData)

	fmt.Println("Server started at http://localhost" + conf.Port)

	log.Fatal(http.ListenAndServe(conf.Port, nil))
}
