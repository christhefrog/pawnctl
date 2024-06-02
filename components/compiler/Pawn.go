package compiler

import (
	"christhefrog/pawnctl/components/github"
	"christhefrog/pawnctl/components/pawnctl"
	"christhefrog/pawnctl/components/util"
	"fmt"
	"path/filepath"

	"golang.org/x/mod/semver"
)

func FetchLatestRelease() (github.Release, error) {
	release, err := github.FetchLatestRelease("pawn-lang", "compiler")

	if err != nil {
		return github.Release{}, err
	}

	return release, nil
}

func FetchRelease(name string) (github.Release, error) {
	release, err := github.FetchRelease("pawn-lang", "compiler", name)

	if err != nil {
		return github.Release{}, err
	}

	return release, nil
}

func Download(release github.Release) error {
	config, err := pawnctl.LoadConfig()
	if err != nil {
		util.Fatalf("Couldn't load global config (%s)", err)
	}

	name := fmt.Sprint("pawnc-", release.Name, "-windows")

	asset, err := release.FindAsset(fmt.Sprint(name, ".zip"))
	if err != nil {
		return err
	}

	path, err := asset.Download("compilers", release.Name)
	if err != nil {
		return err
	}

	util.Unzip(path, filepath.Dir(path))

	exec := fmt.Sprint(filepath.Dir(path), "\\", name, "\\bin\\pawncc.exe")
	config.AddCompiler(release.Name, exec)

	latest, _ := config.GetLatestCompiler()

	if semver.Compare(fmt.Sprint("v", release.Name), fmt.Sprint("v", latest)) > 0 {
		config.AddCompiler("latest", release.Name)
	}

	config.Save()

	return nil
}
