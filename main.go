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
	Packages string
	Options  string
	Upgrade  string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		Packages: os.Getenv("packages"),
		Options:  os.Getenv("options"),
		Upgrade:  os.Getenv("upgrade"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")
	log.Printf("- Packages: %s", configs.Packages)
	log.Printf("- Options: %s", configs.Options)
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
	fmt.Println()

	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	log.Infof("$ brew %s", command.PrintableCommandArgs(false, []string{"update"}))
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
	if configs.Options != "" {
		args, err := shellquote.Split(configs.Options)
		if err != nil {
			log.Errorf("Can't split options: %s", err)
			os.Exit(1)
		}
		cmdArgs = append(cmdArgs, args...)
	}
	packages := strings.Split(configs.Packages, " ")
	cmdArgs = append(cmdArgs, packages...)

	fmt.Println()
	log.Infof("$ brew %s", command.PrintableCommandArgs(false, cmdArgs))
	if err := command.RunCommand("brew", cmdArgs...); err != nil {
		log.Errorf("Can't install formulas:  %s", err)
		os.Exit(1)
	}
}
