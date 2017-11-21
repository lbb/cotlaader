package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	fetchURLDefault = "http://localhost:8080"
	fetchInterval   = time.Second * 20
	volumeEnv       = "COT_VOLUME"
	fetchURLEnv     = "COT_URL"
)

func generateFetchURL() string {
	e := os.Getenv(fetchURLEnv)
	if e != "" {
		return e
	}
	return fetchURLDefault
}

func generateFilePath() string {
	return filepath.Join(os.Getenv(volumeEnv), time.Now().Format(time.RFC3339))
}

// TODO: get mime type and add to file as extension
func newCatPicture() error {
	u := generateFetchURL()
	log.Printf("Will fetch from URL %s\n", u)
	res, err := http.Get(u)
	if err != nil {
		return err
	}

	filename := generateFilePath()
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0444)
	if err != nil {
		return fmt.Errorf("error opening file: %s: %s", filename, err)
	}

	defer res.Body.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return fmt.Errorf("error copying response to file: %s: %s", filename, err)
	}
	return nil
}

func main() {
	log.Printf("Will get cats from %s\n", generateFetchURL())
	log.Printf("Will save cats like %s\n", generateFilePath())
	ticker := time.Tick(fetchInterval)
	for {
		log.Println("Wait for new cat pictures")
		select {
		case <-ticker:
			log.Println("Getting new cats! Yay")
			err := newCatPicture()
			if err != nil {
				log.Printf("Couldn't write cat picture: %s\n", err)
			}
		}
	}
}
