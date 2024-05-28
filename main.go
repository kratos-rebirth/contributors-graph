package main

import (
	"log"
	"os"
)

const (
	DefaultRepo = "Candinya/Kratos-Rebirth"
)

func main() {
	// Get target repo
	repo := os.Getenv("REPO")
	if repo == "" {
		repo = DefaultRepo
	}

	// Get API token
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("Missing token env")
	}

	// Exec
	contributors, err := ListContributors(repo, token)
	if err != nil {
		log.Fatal(err)
	}

	users := FilterUsersOnly(contributors)

	dl := DownloadInfo(users)

	g := GenerateGraph(dl)

	err = os.WriteFile("contributors.svg", []byte(g), 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Done")
}
