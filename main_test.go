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

func Test_brewFileArgs(t *testing.T) {
	t.Log("Create bundle args with no options or verbose option")
	{
		options := ""
		verboseLog := false

		args := brewFileArgs(options, verboseLog, "")
		require.Equal(t, []string{"bundle"}, args)
	}

	t.Log("Create bundle args with verbose log")
	{
		options := ""
		verboseLog := true

		args := brewFileArgs(options, verboseLog, "")
		require.Equal(t, []string{"bundle", "-v"}, args)
	}

	t.Log("Create bundle args with verbose option")
	{
		options := "--verbose"
		verboseLog := true

		args := brewFileArgs(options, verboseLog, "")
		require.Equal(t, []string{"bundle", "--verbose"}, args)
	}

	t.Log("Create bundle args with options")
	{
		options := "option_1 option_2"
		verboseLog := false

		args := brewFileArgs(options, verboseLog, "")
		require.Equal(t, []string{"bundle", "option_1", "option_2"}, args)
	}

	t.Log("Create bundle args with brewfile location")
	{
		options := ""
		verboseLog := false
		path := "path/to/Brewfile"

		args := brewFileArgs(options, verboseLog, path)
		require.Equal(t, []string{"bundle", "--file", "path/to/Brewfile"}, args)
	}
}
