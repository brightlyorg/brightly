package main

import (
	"github.com/brightlyorg/brightly/internal/brightly"
	"os"
)

const (
	brightlyYamlEnvVar           = "BRIGHTLY_YAML"
	defaultBrightlyYamlInputPath = "project"
)

var logger = brightly.GetLogger().Named("Validator")

func main() {
	brightlyYamlInputPath := os.Getenv(brightlyYamlEnvVar)
	if brightlyYamlInputPath == "" {
		logger.Debugf("Env var [%s] not set. Using default: %s", brightlyYamlEnvVar, defaultBrightlyYamlInputPath)
		brightlyYamlInputPath = defaultBrightlyYamlInputPath
	}
	l := logger.With("brightlyYamlInputPath", brightlyYamlInputPath)

	p, err := brightly.ValidateYamlProject(brightlyYamlInputPath)
	if err != nil {
		l.Fatalf("Failed to validate project yaml files: %v", err)
	}
	l.With("project", p).Info("Project yaml files validated successfully")
}
