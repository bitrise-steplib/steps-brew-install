package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	v2command "github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	shellquote "github.com/kballard/go-shellquote"
)

var logger = log.NewLogger()

type Inputs struct {
	Packages          string `env:"packages"`
	Options           string `env:"options"`
	Upgrade           bool   `env:"upgrade,opt[yes,no]"`
	UpgradeDependents bool   `env:"upgrade_dependents,opt[yes,no]"`

	UseBrewfile  bool   `env:"use_brewfile,opt[yes,no]"`
	BrewfilePath string `env:"brewfile_path"`
	CacheEnabled bool   `env:"cache_enabled,opt[yes,no]"`

	VerboseLog bool `env:"verbose_log,opt[yes,no]"`
}

func fail(format string, v ...interface{}) {
	logger.Errorf(format, v...)
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

	if !strings.Contains(options, "--display-times") {
		args = append(args, "--display-times")
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
	cmd := brewCommand([]string{"--cache"}, nil, false)
	logger.Debugf("$ %s", cmd.PrintableCommandArgs())

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
	cmd := brewCommand([]string{"cleanup"}, nil, true)

	logger.Donef("$ %s", cmd.PrintableCommandArgs())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clean homebrew cache directory, error: %s", err)
	}
	return nil
}

func main() {
	var cfg Inputs
	if err := stepconf.Parse(&cfg); err != nil {
		fail("Issue with input: %s", err)
	}

	logger.EnableDebugLog(cfg.VerboseLog)
	envRepo := env.NewRepository()
	cmdFactory := v2command.NewFactory(envRepo)

	stepconf.Print(cfg)
	logger.Println()

	extraEnvs := make(map[string]string)
	var noDependentsCheck string
	if cfg.UpgradeDependents {
		// Need to use the empty string, anything else is parsed as `true`, including "0" and "false"
		noDependentsCheck = ""
	} else {
		noDependentsCheck = "1"
	}
	extraEnvs["HOMEBREW_NO_INSTALLED_DEPENDENTS_CHECK"] = noDependentsCheck
	extraEnvs["HOMEBREW_COLOR"] = "true" // Gets disabled on non-TTY outputs, but we can handle it

	configPrinter := brewConfigPrinter{cmdFactory, envRepo, logger}
	configPrinter.printBrewConfig(extraEnvs)
	logger.Println()

	logger.Infof("Run brew command")

	if cfg.UseBrewfile {
		args := brewFileArgs(cfg.Options, cfg.VerboseLog, cfg.BrewfilePath)
		cmd := brewCommand(args, extraEnvs, true)

		logger.Donef("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			fail("Can't install with Brewfile:  %s", err)
		}
	} else {
		args := cmdArgs(cfg.Options, cfg.Packages, cfg.Upgrade, cfg.VerboseLog)
		cmd := brewCommand(args, extraEnvs, true)

		logger.Donef("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			fail("Can't install formulas:  %s", err)
		}
	}

	// Collecting caches
	if cfg.CacheEnabled {
		fmt.Println()
		logger.Infof("Collecting homebrew cache")

		if err := collectCache(); err != nil {
			logger.Warnf("Cache collection skipped: %s", err)
		} else {
			logger.Donef("Cache path added to $BITRISE_CACHE_INCLUDE_PATHS")
			logger.Printf("Add '%s' step to upload the collected cache for the next build.", colorstring.Yellow("Bitrise.io Cache:Push"))

			fmt.Println()
			logger.Infof("Cleanup homebrew cache")
			if err := cleanCache(); err != nil {
				logger.Warnf("Cache cleanup skipped: %s", err)
			}
		}
	}
}

func brewCommand(args []string, extraEnvs map[string]string, setDefaultOutput bool) *command.Model {
	brewPrefix, err := command.New("brew", "--prefix").RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		logger.Warnf("Failed to get brew prefix: %s\n%s", err, brewPrefix)
	}
	logger.Debugf("Brew prefix: %s", brewPrefix)

	activeArch, err := command.New("arch").RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		logger.Warnf("Failed to get active arch: %s\n%s", err, activeArch)
	}
	logger.Debugf("Active arch: %s", activeArch)

	var effectiveCmd string
	var effectiveArgs []string
	if (activeArch == "i386" || activeArch == "x86_64") && brewPrefix == "/opt/homebrew" {
		// We are running on an Apple Silicon system, but emulated under Rosetta
		// Fix this inconsistency by running brew natively
		effectiveCmd = "arch"
		effectiveArgs = []string{"-arm64", "brew"}
		effectiveArgs = append(effectiveArgs, args...)
	} else {
		effectiveCmd = "brew"
		effectiveArgs = args
	}

	var envStrings []string
	for k, v := range extraEnvs {
		envStrings = append(envStrings, fmt.Sprintf("%s=%s", k, v))
	}
	// Go `cmd.Env` implementation detail: duplicate env vars are handled by applying the last one,
	// so this is fine.
	finalEnvs := append(os.Environ(), envStrings...)

	if setDefaultOutput {
		return command.New(effectiveCmd, effectiveArgs...).SetStdout(os.Stdout).SetStderr(os.Stderr).SetEnvs(finalEnvs...)
	}

	return command.New(effectiveCmd, effectiveArgs...).SetEnvs(finalEnvs...)
}
