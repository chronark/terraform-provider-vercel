package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Release struct {
	Tag string `json:"tag_name"`
}

func main() {
	res, err := http.Get("https://api.github.com/repos/hashicorp/terraform/releases")
	if err != nil {
		panic(err)
	}

	releases := []Release{}

	err = json.NewDecoder(res.Body).Decode(&releases)
	if err != nil {
		panic(err)
	}

	testVersions := make([]string, 0)

	for _, release := range releases {
		split := strings.Split(strings.Trim(release.Tag, "v"), ".")

		major, err := strconv.Atoi(split[0])
		if err != nil {
			panic(err)
		}
		minor, err := strconv.Atoi(split[1])
		if err != nil {
			panic(err)
		}
		if major >= 0 && minor >= 15 {
			testVersions = append(testVersions, strings.TrimLeft(release.Tag, "v"))
		}

	}
	// Manually add some more versions
	testVersions = append(testVersions, "0.14.11")
	testVersions = append(testVersions, "0.13.7")

	err = json.NewEncoder(os.Stdout).Encode(testVersions)
	if err != nil {
		panic(err)
	}
}
