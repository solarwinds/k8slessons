// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

const (
	CurrentVersion = "2.0.0"
)

// RunCompose runs the Docker compose file
func RunCompose() error {
	return sh.RunV("docker-compose", "up", "--build")
}

// DockerBuild builds the container image from the Dockerfile
func DockerBuild() error {
	nameAndTag := fmt.Sprintf("%s:%s", appName(), appTag())
	return sh.RunV("docker", "build", "-t", nameAndTag, ".")
}

// Uses Helm to install Redis
func HelmInstallRedis() error {
	return sh.RunV("helm", "install", "--name", "dumbstore-redis", "-f", "_deployments/redis/values.yml", "stable/redis")
}

func webAppName() string {
	return fmt.Sprintf("%s-web", appName())
}

func appTag() string {
	at := os.Getenv("APP_TAG")
	if at != "" {
		return at
	}
	return CurrentVersion
}

func appName() string {
	return "dumbstore"
}
