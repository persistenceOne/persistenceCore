# Locally testing the workflows

Why would you want to do this? Two reasons:

- **Fast Feedback** - Rather than having to commit/push every time you want to test out the changes you are making to your `.github/workflows/` files (or for any changes to embedded GitHub actions), you can use `act` to run the actions locally. The [environment variables](https://help.github.com/en/actions/configuring-and-managing-workflows/using-environment-variables#default-environment-variables) and [filesystem](https://help.github.com/en/actions/reference/virtual-environments-for-github-hosted-runners#filesystems-on-github-hosted-runners) are all configured to match what GitHub provides.
- **Local Task Runner** - I love [make](<https://en.wikipedia.org/wiki/Make_(software)>). However, I also hate repeating myself. With `act`, you can use the GitHub Actions defined in your `.github/workflows/` locally!


# Installation

## Necessary prerequisites for running `act`

`act` depends on `docker` to run workflows. So install and run `docker` first to get act to work!

### [Homebrew](https://brew.sh/) (Linux/macOS)

```shell
brew install act
```

### [Scoop](https://scoop.sh/) (Windows)

```shell
scoop install act
```

### Bash script

Run this command in your terminal:

```shell
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

### Manual download

Download the [latest release](https://github.com/nektos/act/releases/latest) and add the path to your binary into your PATH.

# Example commands

```sh
# Command structure:
act [<event>] [options]
If no event name passed, will default to "on: push"

# List the actions for the default event:
act -l

# List the actions for a specific event:
act workflow_dispatch -l

# Run the default (`push`) event:
act

# Run a specific event:
act pull_request

# Run a specific job:
act -j test

# Run in dry-run mode:
act -n

# Enable verbose-logging (can be used with any of the above commands)
act -v
```

Additionally, act supports loading environment variables from an .env file. The default is to look in the working directory for the file but can be overridden by:

```shell
act --env-file my.env

```

# Example command to run a workflow

- If you are using Apple M1 chip and you have not specified container architecture, you might encounter issues while running act. If so, try running it with '--container-architecture linux/amd64'.

```shell
act -j docker --env-file act.env --container-architecture linux/amd64 -P ubuntu-latest=catthehacker/ubuntu:act-latest

```

This will run the job `docker` present in any/all of the workflows present, makes use of the env file we have provided and uses the `ubuntu-latest=catthehacker/ubuntu:act-latest` image for the act configuration.


You can chcekout more about `act` [here](https://github.com/nektos/act)