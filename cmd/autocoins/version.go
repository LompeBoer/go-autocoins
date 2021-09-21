package main

import (
	"fmt"
	"log"

	"github.com/LompeBoer/go-autocoins/internal/github"
	"golang.org/x/mod/semver"
)

const (
	Owner = "LompeBoer"    // Owner of the GitHub repo
	Repo  = "go-autocoins" // Repo on GitHub
)

func checkLatestVersion() {
	githubAPI := github.NewAPI(Owner, Repo)
	release, err := githubAPI.LatestRelease()
	if err != nil {
		log.Printf("Unable to get version information: %s\n", err.Error())
	}

	if semver.Compare("v"+VersionNumber, "v"+release.TagName) < 0 {
		fmt.Printf("!\n! Update available (version: %s)\n! Download at: %s\n!\n", release.TagName, release.HTMLURL)
	}
}
