package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func IniConf(path *string) {
	if _, err := os.Stat(*path); os.IsNotExist(err) {
		err := os.MkdirAll(*path, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	editPath := *path + "/playback.html"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {

	}
	editPath = *path + "/styles.css"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {

	}
	editPath = *path + "/config.toml"
	if _, err := os.Stat(editPath); errors.Is(err, os.ErrNotExist) {
		createToml(&editPath)
	}
}

func createToml(editPath *string) {
	file, err := os.Create(*editPath)
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

	fmt.Println("config.toml created successfully.")
}

func createCSS(editPath *string) {
	file, err := os.Create(*editPath)
	if err != nil {
		log.Fatal(err)
	}
	cssContent := `
		@import url('https://fonts.googleapis.com/css2?family=Space+Mono:ital,wght@0,400;0,700;1,400;1,700&display=swap');

		body {
			font-family: 'Space Mono', monospace;
			margin: 0;
			overflow: hidden;
			background: black;
			color: white;
		}

		.container {
			position: relative;
			width: 100vw;
			height: 100px;
			overflow: hidden;
		}

		.scroll-container {
			display: inline-block;
			white-space: nowrap;
			position: absolute;
			left: 0;
			opacity: 0;
			font-size: 50px;
			animation: fadeInScroll 20s linear infinite;
		}

		@keyframes fadeInScroll {
			0% {
				opacity: 0;
				transform: translateX(0);
			}

			10% {
				opacity: 1;
				transform: translateX(0);
			}

			25% {
				opacity: 1;
				transform: translateX(0);
			}

			100% {
				opacity: 1;
				transform: translateX(-100%);
			}
		}
	`
	file.WriteString(cssContent)
}

func createHTML(editPath *string) {
	file, err := os.Create(*editPath)
	if err != nil {
		log.Fatal(err)
	}
	htmlContent := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Playback State</title>
			<link rel="stylesheet" href="/static/styles.css">
			<script>
				function updateText() {
					fetch('/get-playback-data')
						.then(response => response.json())
						.then(data => {
							const trackElement = document.getElementById('track');
							const artistElement = document.getElementById('artist');
							if (data.is_playing) {
								trackElement.textContent = data.item.name;
								artistElement.textContent = data.item.artists.map(artist => artist.name).join(", ");
							} else {
								trackElement.textContent = "No track currently playing.";
								artistElement.textContent = "";
							}
							adjustScrollDuration()
						})
						.catch(error => console.error('Error fetching playback state:', error));
				}

				function adjustScrollDuration() {
					const scrollContainer = document.querySelector('.scroll-container');
					if (!scrollContainer) return;
					const containerWidth = scrollContainer.scrollWidth;
					const baseDuration = 20;
					const duration = baseDuration * (window.innerWidth / containerWidth);
					const minDuration = 10;
					const finalDuration = Math.max(duration, minDuration);
					document.documentElement.style.setProperty('--scroll-duration', duration + 's');
				}

				window.addEventListener('load', adjustScrollDuration);
				setInterval(updateText, 1000);
				window.onload = updateText;
			</script>
		</head>
		<body>
			<div class="container">
				<div class="scroll-container">
					<span id="track">Loading</span>
					<span id="connect">By:</span>
					<span id="artist">Loading</span>
				</div>
			</div>
		</body>
		</html>
	`
	file.WriteString(htmlContent)
}
