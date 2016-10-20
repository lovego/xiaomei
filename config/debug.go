package config

import (
	"os"
	"regexp"
)

func Debug(name string) bool {
	matched, _ := regexp.MatchString(`\b`+name+`\b`, os.Getenv(`debug`))
	return matched
}
