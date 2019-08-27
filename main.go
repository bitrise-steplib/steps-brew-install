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

// configs ...
type configs struct {
	Packages string
	Options  string
	Upgrade  string
}

func createConfigsModelFromEnvs() configs {
	return configs{
		Packages: os.Getenv("packages"),
		Options:  os.Getenv("options"),
		Upgrade:  os.Getenv("upgrade"),
	}
}

func (c configs) print() {
	log.Infof("Configs:")
	log.Printf("- Packages: %s", c.Packages)
	log.Printf("- Options: %s", c.Options)
	log.Printf("- Upgrade: %s", c.Upgrade)
}

func (c configs) validate() error {
	if c.Packages == "" {
		return errors.New("no Packages parameter specified")
	}
	if c.Upgrade != "" && c.Upgrade != "yes" && c.Upgrade != "no" {
		return fmt.Errorf("invalid 'Upgrade' specified (%s), valid options: [yes no]", c.Upgrade)
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
