package main

import (
	"os"
	"testing"
)

func TestListContributors(t *testing.T) {

	token := os.Getenv("TOKEN")

	contributorsList, err := ListContributors("Candinya/Kratos-Rebirth", token)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(contributorsList)
}
