package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func InitConf(PATH *string, AUTHPATH *string) error {
	editPath := *PATH + "/Website"
	if _, err := os.Stat(editPath); os.IsNotExist(err) {
		if err := os.MkdirAll(editPath, os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(*AUTHPATH); os.IsNotExist(err) {
		if err := os.MkdirAll(*AUTHPATH, os.ModePerm); err != nil {
			return err
		}
	}

	editPath = *PATH + "/config.toml"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err := createToml(&editPath)
		if err != nil {
			return err
		}
	}

	editPath = *PATH + "/Website/playback.html"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err := wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/playback.html", &editPath)
		if err != nil {
			return err
		}
	}

	editPath = *PATH + "/Website/script.js"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err := wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/script.js", &editPath)
		if err != nil {
			return err
		}
	}

	editPath = *PATH + "/Website/styles.css"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		err := wget("https://raw.githubusercontent.com/RomanKhundadze/WAYL/main/import/styles.css", &editPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createToml(PATH *string) error {
	file, err := os.Create(*PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, `ClientID     = ""`)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(file, `ClientSecret = ""`)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(file, `Port         = ""`)
	if err != nil {
		return err
	}

	return nil
}

func wget(url string, PATH *string) error {
	cmd := exec.Command("wget", url, "-O", *PATH)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
