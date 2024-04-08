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

	repoOwner := "eamonnk418"
	repoName := "spring-boot"

	repo, _, err := client.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Searching Repository: %s\n", repo.GetFullName())

	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), repoOwner, repoName, "", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Determine the language of the repository based on the files
	repoLanguage := getLanguage(directoryContent)
	if repoLanguage == "" {
		fmt.Println("Unable to determine repository language.")
		return
	}
	fmt.Printf("Detected Repository Language: %s\n", repoLanguage)

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

// Function to determine the language of the repository based on the files present
func getLanguage(contents []*github.RepositoryContent) string {
	for _, content := range contents {
		if content.GetType() == "file" {
			switch content.GetName() {
			case "pom.xml":
				return "Java"
			case "build.gradle", "build.gradle.kts":
				return "Java"
			case "build.sbt":
				return "Scala"
			case "package.json":
				return "JavaScript"
			case "go.mod":
				return "Go"
			}
		}
	}
	return "" // No specific language detected
}
