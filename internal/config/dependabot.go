package config

type Schedule struct {
	Interval string
}

type Update struct {
	PackageEcosystem string
	Directory        string
	Schedule         *Schedule
}

type Dependabot struct {
	Version int
	Updates []*Update
}


func generateMap() {
	result := make(map[string][]string, 0)

	result["gradle"]= []string{"Java", "Kotlin", "Groovy", "Scala", "Kotlin"}
	result["gomod"] = []string{"Go"}
	result["npm"] = []string{"javascript", "nodejs", "npm", "typescript"}
} 




