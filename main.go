package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	shellquote "github.com/kballard/go-shellquote"
)

// configs ...
type configs struct {
	Packages string `env:"packages,required"`
	Options  string `env:"options"`
	Upgrade  bool   `env:"upgrade,opt[yes,no]"`

	VerboseLog bool `env:"verbose_log,opt[yes,no]"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func cmdArgs(options, packages string, upgrade, verboseLog bool) (args []string) {
	if upgrade {
		args = append(args, "reinstall")
	} else {
		args = append(args, "install")
	}
	if verboseLog && !strings.Contains(options, "-v") && !strings.Contains(options, "---verbose") {
		args = append(args, "-v")
	}

	if options != "" {
		o, err := shellquote.Split(options)
		if err != nil {
			fail("Can't split options: %s", err)
		}
		args = append(args, o...)
	}
	p := strings.Split(packages, " ")
	args = append(args, p...)
	return
}

func main() {
	var cfg configs
	if err := stepconf.Parse(&cfg); err != nil {
		fail("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()
	log.SetEnableDebugLog(cfg.VerboseLog)

	log.Infof("$ brew %s", command.PrintableCommandArgs(false, []string{"update"}))
	if err := command.RunCommand("brew", "update"); err != nil {
		log.Errorf("Can't update brew: %s", err)
		os.Exit(1)
	}

	fmt.Println()

	args := cmdArgs(cfg.Options, cfg.Packages, cfg.Upgrade, cfg.VerboseLog)
	log.Infof("$ brew %s", command.PrintableCommandArgs(false, args))

	if err := command.RunCommand("brew", args...); err != nil {
		log.Errorf("Can't install formulas:  %s", err)
		os.Exit(1)
	}
}
