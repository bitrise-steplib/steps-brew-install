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
	Upgrade   string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Packages:  os.Getenv("packages"),
		Arguments: os.Getenv("arguments"),
		Upgrade:   os.Getenv("upgrade"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Packages: %s", configs.Packages)
	log.Printf("- Arguments: %s", configs.Arguments)
	log.Printf("- Upgrade: %s", configs.Upgrade)
}

func (configs ConfigsModel) validate() error {
	if configs.Packages == "" {
		return errors.New("no Packages parameter specified")
	}
	if configs.Upgrade != "" && configs.Upgrade != "yes" && configs.Upgrade != "no" {
		return fmt.Errorf("invalid 'Upgrade' specified (%s), valid options: [yes no]", configs.Upgrade)
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

	if err := command.RunCommand("brew", "update"); err != nil {
		log.Errorf("Can't update brew: %s", err)
		os.Exit(1)
	}

	cmdArgs := []string{}
	if configs.Upgrade == "yes" {
		cmdArgs = append(cmdArgs, "reinstall")
	} else {
		cmdArgs = append(cmdArgs, "install")
	}
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
		log.Errorf("Can't install formulas:  %s", err)
		os.Exit(1)
	}
}
