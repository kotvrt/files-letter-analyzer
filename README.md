# files-letter-analyser

Counts how many times each letter of English alphabet appears in a lodash/lodash GitHub repository and displays the statistics
in standard output.

## Prerequisites

Installed _Golang_ so that you can compile the sources. Check the official documentation for guides on how to install
Golang across different OS: https://go.dev/doc/install

### Quick Instructions for Mac

Assuming you have brew.

> brew install go

You can install brew by running:

> /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

## Compile

> go build -o analyser

## Execute

> ./analyser

## GitHub Token

You're going to need GitHub token exported in local environment.

On Unix systems:

> export GITHUB_TOKEN=<YOUR_GITHUB_TOKEN> 
