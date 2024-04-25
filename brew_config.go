package main

import (
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-io/go-utils/v2/log/colorstring"
)

type brewConfigPrinter struct {
	cmdFactory command.Factory
	envRepo    env.Repository
	logger     log.Logger
}

func (p brewConfigPrinter) printBrewConfig() {
	p.logger.Infof("Homebrew configuration:")

	for _, env := range []string{
		"HOMEBREW_NO_INSTALLED_DEPENDENTS_CHECK",
		"HOMEBREW_NO_INSTALL_FROM_API",
		"HOMEBREW_NO_INSTALL_CLEANUP",
		"HOMEBREW_NO_AUTO_UPDATE",
		"HOMEBREW_CORE_GIT_REMOTE",
	} {
		p.printEnv(env)
	}

	p.logger.Printf("%s: Default values are documented at https://docs.brew.sh/Manpage#environment", colorstring.Yellow("Note"))

	versionStr, err := p.cmdFactory.Create("brew", []string{"--version"}, nil).RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		p.logger.Warnf("Failed to query Homebrew version: %s", err)
	} else {
		p.logger.Printf("Version: %s", colorstring.Cyan(versionStr))
	}
}

func (p brewConfigPrinter) printEnv(env string) {
	value := p.envRepo.Get(env)
	if value == "" {
		value = "<unset>"
	} else {
		value = colorstring.Cyan(value)
	}

	p.logger.Printf("%s: %s", env, value)
}
