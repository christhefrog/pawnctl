package github

import (
	"encoding/json"
	"fmt"
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

func FetchLatestRelease(user string, repo string) (Release, error) {
	url := fmt.Sprint("https://api.github.com/repos/", user, "/", repo, "/releases/latest")

	body, err := Fetch(url)

	if err != nil {
		return Release{}, err
	}

	var release Release
	err = json.Unmarshal(body, &release)
	if err != nil {
		return Release{}, err
	}

	return release, nil
}
