package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

type NowPlayingData struct {
	Title        string
	Album        string
	Artist       string
	PlaybackRate int8
	Duration     float64
	ElapsedTime  float64
}

func parseNowPlayingOutput(output string) (NowPlayingData, error) {
	values := strings.Split(strings.TrimSpace(output), "\n")
	if len(values) != 6 {
		return NowPlayingData{}, fmt.Errorf("unexpected output format")
	}

	title, album, artist, playbackRateStr, durationStr, elapsedTimeStr := values[0], values[1], values[2], values[3], values[4], values[5]

	playbackRate, err := strconv.ParseFloat(playbackRateStr, 64)
	if err != nil {
		return NowPlayingData{}, err
	}
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return NowPlayingData{}, err
	}
	elapsedTime, err := strconv.ParseFloat(elapsedTimeStr, 64)
	if err != nil {
		return NowPlayingData{}, err
	}

	return NowPlayingData{
		Title:        title,
		Album:        album,
		Artist:       artist,
		PlaybackRate: int8(playbackRate),
		Duration:     duration,
		ElapsedTime:  elapsedTime,
	}, nil
}

func main() {
	// Create a new ticker with a 1-second interval
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Create a channel to receive interrupt signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Open the DuckDB database
	db, err := sql.Open("duckdb", "my_music_data.db")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Create a table if it doesn't exist
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS now_playing (recorded_at TIMESTAMP, title TEXT, album TEXT, artist TEXT, playback_rate REAL, duration REAL, elapsed_time REAL)")
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	var previousRecord NowPlayingData

	// Loop until an interrupt signal is received
	for {
		select {
		case <-ticker.C:
			// Run the nowplaying-cli command
			cmd := exec.Command("nowplaying-cli", "get", "title", "album", "artist", "playbackRate", "duration", "elapsedTime")
			output, err := cmd.Output()
			if err != nil {
				fmt.Println("Error executing command:", err)
			} else {
				data, err := parseNowPlayingOutput(string(output))
				if err != nil {
					fmt.Println("Error parsing output:", err)
					continue
				}

				if data == previousRecord {
					continue
				}

				previousRecord = data

				// Get the current epoch time
				epochTime := time.Now()
				fmt.Println(epochTime.Unix(), data.Title, data.Album, data.Artist, data.PlaybackRate, data.Duration, data.ElapsedTime)

				// Insert the data into the table
				_, err = db.Exec("INSERT INTO now_playing VALUES (?, ?, ?, ?, ?, ?, ?)", epochTime, data.Title, data.Album, data.Artist, data.PlaybackRate, data.Duration, data.ElapsedTime)
				if err != nil {
					fmt.Println("Error inserting data:", err)
					continue
				}
			}
		case <-interrupt:
			fmt.Println("Interrupt signal received. Exiting...")
			return
		}
	}
}
