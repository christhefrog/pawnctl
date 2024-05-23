package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Release struct {
	Url       string    `json:"html_url"`
	Name      string    `json:"name"`
	Published time.Time `json:"published_at"`
	Assets    []Asset   `json:"assets"`
}

type Asset struct {
	Name string `json:"name"`
	Url  string `json:"browser_download_url"`
}

func FetchLatestRelease(user string, repo string) Release {
	url := fmt.Sprint("https://api.github.com/repos/", user, "/", repo, "/releases/latest")

	body := Fetch(url)

	var release Release
	err := json.Unmarshal(body, &release)
	if err != nil {
		log.Fatal(err)
	}

	return release
}

func Fetch(url string) []byte {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "pawnman")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	} else {
		log.Fatal("API responded with []")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}
