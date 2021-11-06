# bashcord

discord rich presence for bash

## Installation

`go get github.com/deletescape/bashcord@latest`

## Setup

Put the following into your `.bashrc`, this makes sure your bash history file is immediately updated after every command

```bash
shopt -s histappend
shopt -s cmdhist
export HISTCONTROL=ignoreboth
PROMPT_COMMAND='history -a;history -n'
```

## Running

Literally just run `bashcord`. If you wanna go wild with it you could probably make it a system service so it runs on boot.