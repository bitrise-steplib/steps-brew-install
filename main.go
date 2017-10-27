package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/kballard/go-shellquote"
)

// ConfigsModel ...
type ConfigsModel struct {
	Packages  string
	Arguments string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Packages:  os.Getenv("packages"),
		Arguments: os.Getenv("arguments"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Packages: %s", configs.Packages)
	log.Printf("- Arguments: %s", configs.Arguments)
}

func (configs ConfigsModel) validate() error {
	if configs.Packages == "" {
		return errors.New("no Packages parameter specified")
	}
	return nil
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	cmdArgs := []string{"install"}
	if configs.Arguments != "" {
		args, err := shellquote.Split(configs.Arguments)
		if err != nil {
			log.Errorf("Can't split arguments: %s", err)
			os.Exit(1)
		}
		cmdArgs = append(cmdArgs, args...)
	}
	packages := strings.Split(configs.Packages, " ")
	cmdArgs = append(cmdArgs, packages...)

	if err := command.RunCommand("brew", cmdArgs...); err != nil {
		log.Errorf("Can't run command %s", err)
		os.Exit(1)
	}
}
