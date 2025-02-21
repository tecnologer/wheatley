package envvar

import "os"

func ValueWithDefault(envVarName string, def string) string {
	value := os.Getenv(envVarName)
	if value == "" {
		return def
	}

	return value
}
