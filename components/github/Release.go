package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	} else if release.Name == "" {
		return Release{}, errors.New("specified repository doesn't exist or has no releases")
	}

	return release, nil
}

func (r Release) FindAsset(name string) (Asset, error) {
	for _, v := range r.Assets {
		if v.Name == name {
			return v, nil
		}
	}
	return Asset{}, errors.New("asset not found")
}

func (a Asset) Download(path ...string) (string, error) {
	dest := filepath.Dir(os.Args[0])

	for _, v := range path {
		dest = fmt.Sprint(dest, "\\", v)
		os.Mkdir(dest, 0755)
	}

	dest = fmt.Sprint(dest, "\\", a.Name)

	bytes, err := Fetch(a.Url)
	if err != nil {
		return "", err
	}

	os.WriteFile(dest, bytes, 0664)

	return dest, nil
}
