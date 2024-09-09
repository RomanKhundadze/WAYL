package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func InitConf(PATH *string) {
	editPath := *PATH + "/Website"
	if _, err := os.Stat(editPath); os.IsNotExist(err) {
		if err := os.MkdirAll(editPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	editPath = *PATH + "/config.toml"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		createToml(&editPath)
	}

	editPath = *PATH + "/Website/playback.html"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err = wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/playback.html", &editPath)
	}

	editPath = *PATH + "/Website/script.js"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err = wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/script.js", &editPath)
	}

	editPath = *PATH + "/Website/styles.css"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err = wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/styles.css", &editPath)
	}
}

func createToml(PATH *string) {
	file, err := os.Create(*PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, `ClientID     = ""`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, `ClientSecret = ""`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, `Port         = ""`)
	if err != nil {
		log.Fatal(err)
	}
}

func wget(url string, PATH *string) error {
	cmd := exec.Command("wget", url, "-O", *PATH)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
