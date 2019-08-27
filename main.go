package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/kballard/go-shellquote"
)

// configs ...
type configs struct {
	Packages string `env:"packages,required"`
	Options  string `env:"options"`
	Upgrade  string `env:"upgrade,opt[yes,no]"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	var cfg configs
	if err := stepconf.Parse(&cfg); err != nil {
		fail("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()
	// log.SetEnableDebugLog(cfg.VerboseLog)

	log.Infof("$ brew %s", command.PrintableCommandArgs(false, []string{"update"}))
	if err := command.RunCommand("brew", "update"); err != nil {
		log.Errorf("Can't update brew: %s", err)
		os.Exit(1)
	}

	cmdArgs := []string{}
	if cfg.Upgrade == "yes" {
		cmdArgs = append(cmdArgs, "reinstall")
	} else {
		cmdArgs = append(cmdArgs, "install")
	}
	if cfg.Options != "" {
		args, err := shellquote.Split(cfg.Options)
		if err != nil {
			log.Errorf("Can't split options: %s", err)
			os.Exit(1)
		}
		cmdArgs = append(cmdArgs, args...)
	}
	packages := strings.Split(cfg.Packages, " ")
	cmdArgs = append(cmdArgs, packages...)

	fmt.Println()
	log.Infof("$ brew %s", command.PrintableCommandArgs(false, cmdArgs))
	if err := command.RunCommand("brew", cmdArgs...); err != nil {
		log.Errorf("Can't install formulas:  %s", err)
		os.Exit(1)
	}
}
