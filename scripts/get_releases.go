package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Release struct {
	Tag string `json:"tag_name"`
}

func main() {
	res, err := http.Get("https://api.github.com/repos/hashicorp/terraform/releases")
	if err != nil {
		log.Fatalln(err)
	}

	releases := []Release{}

	err = json.NewDecoder(res.Body).Decode(&releases)
	if err != nil {
		log.Fatalln(err)
	}

	versions := make([]string, len(releases))

	for i, release := range releases {
		versions[i] = strings.TrimLeft(release.Tag, "v")
	}


	json.NewEncoder(os.Stdout).Encode(versions)

}
