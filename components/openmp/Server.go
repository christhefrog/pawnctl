package openmp

import (
	"christhefrog/sampman/components/github"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"
	"fmt"
	"path/filepath"
)

func FetchLatestServer() (github.Release, error) {
	release, err := github.FetchLatestRelease("openmultiplayer", "open.mp")

	if err != nil {
		return github.Release{}, err
	}

	return release, nil
}

func Download(release github.Release) error {
	name := "open.mp-win-x86"

	asset, err := release.FindAsset(fmt.Sprint(name, ".zip"))
	if err != nil {
		return err
	}

	path, err := asset.Download("servers", release.Name)
	if err != nil {
		return err
	}

	util.Unzip(path, fmt.Sprint("servers\\", release.Name))

	server := sampman.Server{
		Type:     "omp",
		Path:     fmt.Sprint(filepath.Dir(path), "\\Server"),
		Exec:     fmt.Sprint(filepath.Dir(path), "\\Server\\omp-server.exe"),
		Includes: fmt.Sprint(filepath.Dir(path), "\\Server\\qawno\\include"),
	}

	sampman.AddServer(release.Name, server)
	sampman.Save()

	return nil
}
