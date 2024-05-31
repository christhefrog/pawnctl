package compiler

import (
	"christhefrog/sampman/components/github"
	"christhefrog/sampman/components/sampman"
	"christhefrog/sampman/components/util"
	"fmt"
	"path/filepath"
)

func FetchLatestRelease() (github.Release, error) {
	release, err := github.FetchLatestRelease("pawn-lang", "compiler")

	if err != nil {
		return github.Release{}, err
	}

	return release, nil
}

func Download(release github.Release) error {
	name := fmt.Sprint("pawnc-", release.Name, "-windows")

	asset, err := release.FindAsset(fmt.Sprint(name, ".zip"))
	if err != nil {
		return err
	}

	path, err := asset.Download("compilers", release.Name)
	if err != nil {
		return err
	}

	util.Unzip(path, fmt.Sprint("compilers\\", release.Name))

	exec := fmt.Sprint(filepath.Dir(path), "\\", name, "\\bin\\pawncc.exe")
	sampman.AddCompiler(release.Name, exec)
	sampman.AddCompiler("latest", release.Name)

	sampman.Save()

	return nil
}
