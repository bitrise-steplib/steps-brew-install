package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_cmdArgs(t *testing.T) {
	t.Log("Create args for: install go")
	{
		upgrade := false
		options := ""
		packages := "go"
		verboseLog := false

		args := cmdArgs(options, packages, upgrade, verboseLog)
		require.Equal(t, []string{"install", "go"}, args)
	}

	t.Log("Create args for: upgrade go")
	{
		upgrade := true
		options := ""
		packages := "go"
		verboseLog := false

		args := cmdArgs(options, packages, upgrade, verboseLog)
		require.Equal(t, []string{"reinstall", "go"}, args)
	}

	t.Log("Create args with verbose log")
	{
		upgrade := false
		options := ""
		packages := "go"
		verboseLog := true

		args := cmdArgs(options, packages, upgrade, verboseLog)
		require.Equal(t, []string{"install", "-v", "go"}, args)
	}

	t.Log("Create args with verbose option")
	{
		upgrade := false
		options := "--verbose"
		packages := "go"
		verboseLog := true

		args := cmdArgs(options, packages, upgrade, verboseLog)
		require.Equal(t, []string{"install", "--verbose", "go"}, args)
	}

	t.Log("Create args with options")
	{
		upgrade := false
		options := "option_1 option_2"
		packages := "go"
		verboseLog := false

		args := cmdArgs(options, packages, upgrade, verboseLog)
		require.Equal(t, []string{"install", "option_1", "option_2", "go"}, args)
	}
}
