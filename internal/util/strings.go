package util

import (
	"os"
	"regexp"
)

func ReplaceEnvVars(value string) string {
	pattern := `\$\{([^:]+):([^}]+)\}` // Regex to match ${VAR:default}
	r := regexp.MustCompile(pattern)

	return r.ReplaceAllStringFunc(value, func(match string) string {
		groups := r.FindStringSubmatch(match)
		envVar := groups[1]         // Environment variable name
		defaultVal := groups[2]     // Default value
		envVal := os.Getenv(envVar) // Return env variable or default
		if envVal == "" {
			return defaultVal
		} else {
			return envVal
		}
	})
}
