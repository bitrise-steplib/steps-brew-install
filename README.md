# Brew install

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-brew-install?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-brew-install/releases)

Install or upgrade dependencies with Homebrew.

<details>
<summary>Description</summary>


Install or upgrade dependencies using Homebrew, a package manager for MacOS. 

### Configuring the Step 

Homebrew defines the available packages as formulae. Our Step needs the name of the Homebrew formulae you want to use, either specified as a step input, or from a Brewfile in the project's source.

To specify formulae in the step configuration

1. In the **Formula name** input, put the name of the formula you want to download. 
1. In the **Upgrade formula?** input, set the default behavior for previously installed packages. If the input is set to `yes`, the Step will call `brew reinstall` to upgrade them to the latest version.
1. In the **Brew install/reinstall options** input, you can set additional flags for the `brew install` or `brew reinstall` commands. 
   For the possible options, see [Homebrew's documentation](https://docs.brew.sh/Manpage#install-options-formulacask).

Alternatively you can install formulae using a Brewfile

1. Add a `Brewfile` to the root of the project's source. For the format of the Brewfile, see the [Homebrew Bundle documentation](https://github.com/Homebrew/homebrew-bundle#usage)
1. Set the **Use a Brewfile to install packages?** input to "yes". 
1. (optional) Set the **Path to the Brewfile** input if it is not in the root of the project's source

### Useful links

- [Homebrew documentation](https://docs.brew.sh/Manpage)
- [Caching Homebrew installers](https://devcenter.bitrise.io/builds/caching/caching-homebrew-installers/)

### Related Steps 

- [Run yarn command](https://www.bitrise.io/integrations/steps/yarn)
- [Run npm command](https://www.bitrise.io/integrations/steps/npm)
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `packages` | Name of the formulas to install. Multiple formulas can be specified by separating them with a space, e.g. `git-lfs sqlite pipx`  This input must be specified when `use_brewfile` is `no` |  |  |
| `upgrade` | If set to `"yes"`, the step will upgrade the defined packages by calling `brew reinstall [options] [packages]` command. Otherwise the step calls `brew install [options] [packages]`.  |  | `yes` |
| `use_brewfile` | If set to `"yes"`, the step will install packages in the Brewfile by running `brew bundle`. If no Brewfile path is set, it assumes a Brewfile exists in the current directory.  |  | `no` |
| `brewfile_path` | If set, `use_brewfile` must be set to `yes`. Path must end with `Brewfile`  |  |  |
| `options` | Flags to pass to the brew install/reinstall command. `brew install/reinstall [options] [packages]`  |  |  |
| `cache_enabled` | If set to `"yes"` the contents of `~/Library/Caches/Homebrew` directory will be cached.  | required | `no` |
| `verbose_log` | Should the step print more detailed log? | required | `no` |
</details>

<details>
<summary>Outputs</summary>
There are no outputs defined in this step
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-brew-install/pulls) and [issues](https://github.com/bitrise-steplib/steps-brew-install/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
