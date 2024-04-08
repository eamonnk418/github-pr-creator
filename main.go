package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v61/github"
)

type GitHubRepository struct {
	Name     string
	Language string
}

var languageEcosystem = map[string][]string{
	"Java":   {"maven", "gradle"},
	"Kotlin": {"gradle"},
	"Scala":  {"maven", "gradle"},
}

var packageEcosystem = map[string][]string{
	"maven":  {"pom.xml"},
	"gradle": {"build.gradle", "build.gradle.kts"},
	"npm":    {"package.json"},
	"gomod":  {"go.mod"},
}

func main() {
	client := github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))

	repo, _, err := client.Repositories.Get(context.Background(), "eamonnk418", "spring-boot")
	if err != nil {
		log.Fatal(err)
	}

	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), "eamonnk418", "jenkins", "", nil)
	if err != nil {
		log.Fatal(err)
	}

	repoLanguage := repo.GetSource().GetLanguage()
	packageManagers, ok := languageEcosystem[repoLanguage]
	if !ok {
		fmt.Printf("Language '%s' not supported.\n", repoLanguage)
		return
	}

	for _, packageManager := range packageManagers {
		files := packageEcosystem[packageManager]
		if files == nil {
			fmt.Printf("No files associated with package manager '%s'.\n", packageManager)
			continue
		}

		found := false
		for _, content := range directoryContent {
			for _, file := range files {
				if content.GetName() == file {
					fmt.Printf("Found file '%s' associated with package manager '%s'.\n", file, packageManager)
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			fmt.Printf("No '%s' file found associated with package manager '%s'.\n", files, packageManager)
		}
	}
}
