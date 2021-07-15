# LSGO
Unix/Linux coreutil <b>ls</b> command implementation in golang

## Usage

### Basic usage
```
Usage:
  lsgo [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  list        Shows the files and folders in the current directory
  version     Print out the version information

Flags:
  -h, --help   help for lsgo

Use "lsgo [command] --help" for more information about a command.
```

### list usage
```
list command lists down all the files and folders present in the current directory. 
By default, it hides the files which starts with .
Use appropriate flags in order to see such kind of hidden files

Usage:
  lsgo list [flags]

Flags:
  -a, --all           List all files including hidden files
  -h, --help          help for list
  -l, --long          List all files including hidden files
  -p, --path string   See contents of a path
  -r, --readable      Prints output in human readable format. Works with long lists
  -s, --sort string   Sorts rows by column names. Works with long lists. Sorts alphabetically only (default "NAME")
```

## Installation
Download the package specific to your OS from releases page and copy it to /usr/local/bin for linux, or keep it at a particular location and set its path. 