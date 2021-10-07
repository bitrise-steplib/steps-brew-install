package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	shellquote "github.com/kballard/go-shellquote"
)

// configs ...
type configs struct {
	Packages string `env:"packages"`
	Options  string `env:"options"`
	Upgrade  bool   `env:"upgrade,opt[yes,no]"`

	UseBrewfile  bool   `env:"use_brewfile,opt[yes,no]"`
	BrewfilePath string `env:"brewfile_path"`
	CacheEnabled bool   `env:"cache_enabled,opt[yes,no]"`

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
	if packages == "" {
		fail("No packages provided, and not using a Brewfile")
	}

	p := strings.Split(packages, " ")
	args = append(args, p...)
	return
}

func brewFileArgs(options string, verboseLog bool, path string) (args []string) {
	args = append(args, "bundle")
	if verboseLog && !strings.Contains(options, "-v") && !strings.Contains(options, "---verbose") {
		args = append(args, "-v")
	}

	if path != "" {
		if strings.HasSuffix(path, "Brewfile") {
			args = append(args, "--file", path)
		} else {
			fail("Brewfile path must include the filename")
		}
	}

	if options != "" {
		o, err := shellquote.Split(options)
		if err != nil {
			fail("Can't split options: %s", err)
		}
		args = append(args, o...)
	}
	return
}

func collectCache() error {
	cmd := command.New("brew", "--cache")
	log.Debugf("$ %s", cmd.PrintableCommandArgs())

	brewCachePth, err := cmd.RunAndReturnTrimmedOutput()
	if err != nil {
		return fmt.Errorf("failed to find homebrew chache directory, error: %s", err)
	}

	brewCache := cache.New()
	brewCache.IncludePath(brewCachePth)
	if err := brewCache.Commit(); err != nil {
		return fmt.Errorf("failed to commit cache paths, error: %s", err)

	}
	return nil
}

func cleanCache() error {
	cmd := command.New("brew", "cleanup").SetStdout(os.Stdout).SetStderr(os.Stderr)

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clean homebrew chache directory, error: %s", err)
	}
	return nil
}

func main() {
	var cfg configs
	if err := stepconf.Parse(&cfg); err != nil {
		fail("Issue with input: %s", err)
	}

	stepconf.Print(cfg)
	fmt.Println()
	log.SetEnableDebugLog(cfg.VerboseLog)

	log.Infof("Run brew command")
	if cfg.UseBrewfile {
		args := brewFileArgs(cfg.Options, cfg.VerboseLog, cfg.BrewfilePath)
		cmd := command.New("brew", args...).SetStdout(os.Stdout).SetStderr(os.Stderr)

		log.Donef("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			fail("Can't install with Brewfile:  %s", err)
		}
	} else {
		args := cmdArgs(cfg.Options, cfg.Packages, cfg.Upgrade, cfg.VerboseLog)
		cmd := command.New("brew", args...).SetStdout(os.Stdout).SetStderr(os.Stderr)

		log.Donef("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			fail("Can't install formulas:  %s", err)
		}
	}

	// Collecting caches
	if cfg.CacheEnabled {
		fmt.Println()
		log.Infof("Collecting homebrew cache")

		if err := collectCache(); err != nil {
			log.Warnf("Cache collection skipped: %s", err)
		} else {
			log.Donef("Cache path added to $BITRISE_CACHE_INCLUDE_PATHS")
			log.Printf("Add '%s' step to upload the collected cache for the next build.", colorstring.Yellow("Bitrise.io Cache:Push"))

			fmt.Println()
			log.Infof("Cleanup homebrew cache")
			if err := cleanCache(); err != nil {
				log.Warnf("Cache cleanup skipped: %s", err)
			}
		}
	}
}
