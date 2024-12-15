# A Powerline style prompt for your shell

## DISCLAIMER

This is a fork of [powerline-go](https://github.com/gentoomaniac/powerline-go) by [@justjanne](https://github.com/justjanne).

## powerline-go

A [Powerline](https://github.com/Lokaltog/vim-powerline) like prompt for Bash,
ZSH and Fish. Based on [Powerline-Shell](https://github.com/banga/powerline-shell) by @banga.
Ported to golang by @justjanne.

![Solarized+Powerline](https://raw.github.com/gentoomaniac/powerline-go/main/preview.png)

- Shows some important details about the git/hg branch (see below)
- Changes color if the last command exited with a failure code
- If you're too deep into a directory tree, shortens the displayed path with an ellipsis
- Shows the current Python [virtualenv](http://www.virtualenv.org/) environment
- Shows the current Ruby version using [rbenv](https://github.com/rbenv/rbenv) or [rvm](https://rvm.io/)
- Shows if you are in a [nix](https://nixos.org/) shell
- It's easy to customize and extend. See below for details.

## Table of Contents

- [A Powerline style prompt for your shell](#a-powerline-style-prompt-for-your-shell)
  - [DISCLAIMER](#disclaimer)
  - [powerline-go](#powerline-go)
  - [Table of Contents](#table-of-contents)
  - [Version Control](#version-control)
  - [Installation](#installation)
    - [Precompiled Binaries](#precompiled-binaries)
    - [Other Platforms](#other-platforms)
    - [Bash](#bash)
    - [ZSH](#zsh)
    - [Fish](#fish)
    - [Nix](#nix)
    - [Powershell](#powershell)
  - [Customization](#customization)
    - [Eval](#eval)
        - [Bash](#bash-1)
        - [ZSH](#zsh-1)
        - [Fish](#fish-1)
    - [Path Aliases](#path-aliases)
    - [Duration](#duration)
      - [Bash](#bash-2)
      - [Zsh](#zsh-2)
      - [Fish](#fish-2)
  - [License](#license)

## Version Control

All of the version control systems supported by powerline shell give you a
quick look into the state of your repo:

- The current branch is displayed and changes background color when the
  branch is dirty.
- When the local branch differs from the remote, the difference in number
  of commits is shown along with `⇡` or `⇣` indicating whether a git push
  or pull is pending

In addition, git has a few extra symbols:

- `✎` -- a file has been modified, but not staged for commit
- `✔` -- a file is staged for commit
- `✼` -- a file has conflicts
- `+` -- untracked files are present
- `⚑` -- stash is present

Each of these will have a number next to it if more than one file matches.

## Installation

Requires Go 1.15+

`powerline-go` uses ANSI color codes, these should nowadays work everywhere,
but you may have to set your $TERM to `xterm-256color` for it to work.

If you want to use the "patched" mode (which is the default, and provides
improved UI), you'll need to install a powerline font, either as fallback,
or by patching the font you use for your terminal: see
[powerline-fonts](https://github.com/Lokaltog/powerline-fonts).
Alternatively you can use "compatible" or "flat" mode.

### Precompiled Binaries

I provide precompiled binaries for x64 Linux and macOS in the
[releases tab](https://github.com/gentoomaniac/powerline-go/releases)

### Other Platforms

- Install (and update) the package with

```bash
go install github.com/gentoomaniac/powerline-go@latest
```

- By default it will be in `$GOPATH/bin`, if you want to change that, you can set
  your `$GOPATH` and/or `$GOBIN`, but will need to change the path in the
  following scripts, too.

### Bash

Add the following to your `.bashrc`:

```bash
function _update_ps1() {
    PS1="$($GOPATH/bin/powerline-go -error $? -jobs $(jobs -p | wc -l))"

    # Uncomment the following line to automatically clear errors after showing
    # them once. This not only clears the error for powerline-go, but also for
    # everything else you run in that shell. Don't enable this if you're not
    # sure this is what you want.

    #set "?"
}

if [ "$TERM" != "linux" ] && [ -f "$GOPATH/bin/powerline-go" ]; then
    PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"
fi
```

Currently, right prompt support is not available when using bash.

### ZSH

Add the following to your `.zshrc`:

```bash
function powerline_precmd() {
    PS1="$($GOPATH/bin/powerline-go -error $? -jobs ${${(%):%j}:-0})"

    # Uncomment the following line to automatically clear errors after showing
    # them once. This not only clears the error for powerline-go, but also for
    # everything else you run in that shell. Don't enable this if you're not
    # sure this is what you want.

    #set "?"
}

function install_powerline_precmd() {
  for s in "${precmd_functions[@]}"; do
    if [ "$s" = "powerline_precmd" ]; then
      return
    fi
  done
  precmd_functions+=(powerline_precmd)
}

if [ "$TERM" != "linux" ] && [ -f "$GOPATH/bin/powerline-go" ]; then
    install_powerline_precmd
fi
```

### Fish

Redefine `fish_prompt` in `~/.config/fish/config.fish`:

```bash
function fish_prompt
    eval $GOPATH/bin/powerline-go -error $status -jobs (count (jobs -p))
end
```
### Nix

When using `nix-shell --pure`, `powerline-go` will not be accessible, and
your prompt will disappear.

To work around this you can add this snippet to your `.bashrc`,
which should re-enable the prompt in most cases:

```bash
# Workaround for nix-shell --pure
if [ "$IN_NIX_SHELL" == "pure" ]; then
    if [ -x "$HOME/.nix-profile/bin/powerline-go" ]; then
        alias powerline-go="$HOME/.nix-profile/bin/powerline-go"
    elif [ -x "/run/current-system/sw/bin/powerline-go" ]; then
        alias powerline-go="/run/current-system/sw/bin/powerline-go"
    fi
fi
```

### Powershell

Redefine `prompt` function on your profile:

```powershell
# Load powerline-go prompt
function global:prompt {
    $pwd = $ExecutionContext.SessionState.Path.CurrentLocation
    $startInfo = New-Object System.Diagnostics.ProcessStartInfo
    $startInfo.FileName = "powerline-go"
    $startInfo.Arguments = "-shell bare"
    $startInfo.Environment["TERM"] = "xterm-256color"
    $startInfo.CreateNoWindow = $true
    $startInfo.StandardOutputEncoding = [System.Text.Encoding]::UTF8
    $startInfo.RedirectStandardOutput = $true
    $startInfo.UseShellExecute = $false
    $startInfo.WorkingDirectory = $pwd
    $process = New-Object System.Diagnostics.Process
    $process.StartInfo = $startInfo
    $process.Start() | Out-Null
    $standardOut = $process.StandardOutput.ReadToEnd()
    $process.WaitForExit()
    $standardOut
}
```

Use `ProcessStartInfo` is needed to allow fill the enviromnet variables required by powerline-go.

## Customization

There are a few optional arguments which can be seen by running
`powerline-go -help`. These can be used by changing the command you have set
in your shell’s init file.

```
Usage: powerline-go [flags]

Flags:
  -h, --help                                                              Show context-sensitive help.
  -v, --verbosity                                                         Increase verbosity.
  -q, --[no-]quiet                                                        Disable all logging
      --[no-]json                                                         Log as json
      --[no-]debug                                                        shortcut for -vvvv
      --save-config                                                       Save default config and exit
      --cwd-mode="fancy"                                                  How to display the current directory
      --cwd-max-depth=5                                                   Maximum number of directories to show in path
      --cwd-max-dir-size=-1                                               Maximum number of letters displayed for each directory in the path
      --colorize-hostname                                                 Colorize the hostname based on a hash of itself, or use the PLGO_HOSTNAMEFG and PLGO_HOSTNAMEBG env vars (both need to be set).
      --hostname-only-if-ssh                                              Show hostname only for SSH connections
      --ssh-alternate-icon                                                Show the older, original icon for SSH connections
      --east-asian-width                                                  Use East Asian Ambiguous Widths
      --newline                                                           Show the prompt on a new line
      --static-prompt-indicator                                           Always show the prompt indicator with the default color, never with the error color
      --venv-name-size-limit=0                                            Show indicator instead of virtualenv name if name is longer than this limit (defaults to 0, which is unlimited)
      --jobs=0                                                            Number of jobs currently running
      --git-assume-unchanged-size=2048                                    Disable checking for changed/edited files in git repositories where the index is larger than this size (in KB), improves performance
      --git-disable-stats=GIT-DISABLE-STATS,...                           Comma-separated list to disable individual git statuses, (valid choices: ahead, behind, staged, notStaged, untracked, conflicted, stashed)
      --git-mode="fancy"                                                  How to display git status, (valid choices: fancy, compact, simple)
      --mode="patched"                                                    The characters used to make separators between segments, (valid choices: patched, compatible, flat)
      --theme="default"                                                   Set this to the theme you want to use, (valid choices: default, low-contrast, gruvbox, solarized-dark16, solarized-light16)
      --shell="autodetect"                                                overwrite the shell (valid choices: autodetect, bare, bash, zsh)
      --modules=venv,user,host,ssh,cwd,perms,git,hg,jobs,exit,root,...    The list of modules to load, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable plugins and should output
                                                                          a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs.
      --modules-right=MODULES-RIGHT,...                                   The list of modules to load anchored to the right, separated by ','. Unrecognized modules will be invoked as 'powerline-go-MODULE' executable
                                                                          plugins and should output a (possibly empty) list of JSON objects that unmarshal to powerline-go's Segment structs.
      --priority=root,cwd,user,host,ssh,perms,git-branch,git-status,hg,jobs,exit,cwd-path,...
                                                                          Segments sorted by priority, if not enough space exists, the least priorized segments are removed first. Separate with ','
      --max-width-percentage=0                                            Maximum width of the shell that the prompt may use, in percent. Setting this to 0 disables the shrinking subsystem.
      --truncate-segment-width=16                                         Maximum width of a segment, segments longer than this will be shortened if space is limited. Setting this to 0 disables it.
      --prev-error=0                                                      Exit code of previously executed command
      --numeric-exit-codes                                                Shows numeric exit codes for errors.
      --ignore-repos=IGNORE-REPOS,...                                     A list of git repos to ignore. Separate with ','. Repos are identified by their root directory.
      --shorten-gke-names                                                 Shortens names for GKE Kube clusters.
      --shorten-eks-names                                                 Shortens names for EKS Kube clusters.
      --shorten-openshift-names                                           Shortens names for Openshift Kube clusters.
      --shell-var=STRING                                                  A shell variable to add to the segments.
      --shell-var-no-warn-empty                                           Disables warning for empty shell variable.
      --trim-ad-domain
      --duration=STRING                                                   The elapsed clock-time of the previous command
      --duration-min="0"                                                  The minimal time a command has to take before the duration segment is shown
      --duration-low-precision                                            Use low precision timing for duration with milliseconds as maximum resolution
      --eval                                                              Output prompt in 'eval' format.
      --condensed                                                         Remove spacing between segments
      --time-format="15:04:05"                                            The layout string how a reference time should be represented. The reference time is predefined and not user choosen. Consult the golang
                                                                          documentation for details: https://pkg.go.dev/time#example-Time.Format
      --vi-mode=STRING                                                    The current vi-mode (eg. KEYMAP for zsh) for vi-module module
      --path-aliases=KEY=VALUE,...                                        One or more aliases from a path to a short name. Separate with ','. An alias maps a path like foo/bar/baz to a short name like FBB. Specify these
                                                                          as key/value pairs like foo/bar/baz=FBB. Use '~' for your home dir. You may need to escape this character to avoid shell substitution.
```

### Eval

If using `eval` and `--modules-right` is desired, the shell setup must be modified slightly, as shown below:

##### Bash

Add the following to your `.bashrc`:

```bash
function _update_ps1() {
    eval "$($GOPATH/bin/powerline-go -error $? -shell bash -eval -modules-right git)"
}

if [ "$TERM" != "linux" ] && [ -f "$GOPATH/bin/powerline-go" ]; then
    PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"
fi
```

##### ZSH

Add the following to your `.zshrc`:

```bash
function powerline_precmd() {
    eval "$($GOPATH/bin/powerline-go -error $? -shell zsh -eval -modules-right git)"
}

function install_powerline_precmd() {
  for s in "${precmd_functions[@]}"; do
    if [ "$s" = "powerline_precmd" ]; then
      return
    fi
  done
  precmd_functions+=(powerline_precmd)
}

if [ "$TERM" != "linux" ]; then
    install_powerline_precmd
fi
```

##### Fish

Eval mode (and `modules-right` support) for Fish is not currently available.

### Path Aliases

The point of the path aliases feature is to allow you to replace long paths
with a shorter string that you can understand more quickly. This is useful if
you're often in deep path hierarchies that end up consuming most of your
terminal width, even when some portions are replaced by an ellipsis.

For example, you might want to replace the string `$GOPATH/src/github.com` with
`@GOPATH-GH`. When you're in a directory like
`$GOPATH/src/github.com/gentoomaniac/powerline-go`, you'll instead see `@GOPATH-GH >
justjanne > powerline-go` in the shell prompt.

Aliases are defined as comma-separated key value pairs, like this:

```bash
powerline-go ... -path-aliases \$GOPATH/src/github.com=@GOPATH-GH,\~/work/projects/foo=@FOO,\~/work/projects/bar=@BAR
```

Note that you should use `~` instead of `/home/username` when specifying the
path. Also make sure to escape the `~` character. Otherwise your shell will
perform interpolation on it before `powerline-go` can see it!

### Duration

The duration segment requires some assistance from the shell.  The shell must have a hook that gets executed immediately before the command.

#### Bash

Bash 4.4 includes an easy way to get a start-time, using `$PS0`.  However, not all operating systems come with a sufficiently recent version of Bash installed.  This example only has seconds precision.  Add or modify your `.bashrc` file to include the following:

```bash
INTERACTIVE_BASHPID_TIMER="/tmp/${USER}.START.$$"

PS0='$(echo $SECONDS > "$INTERACTIVE_BASHPID_TIMER")'

function _update_ps1() {
  local __ERRCODE=$?

  local __DURATION=0
  if [ -e $INTERACTIVE_BASHPID_TIMER ]; then
    local __END=$SECONDS
    local __START=$(cat "$INTERACTIVE_BASHPID_TIMER")
    __DURATION="$(($__END - ${__START:-__END}))"
    rm -f "$INTERACTIVE_BASHPID_TIMER"
  fi

  PS1="$($GOPATH/bin/powerline-go -modules duration -duration $__DURATION -error $__ERRCODE -shell bash)"
}

if [ "$TERM" != "linux" ] && [ -f "$GOPATH/bin/powerline-go" ]; then
  PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"
fi
```

#### Zsh

Using `$EPOCHREALTIME` requires loading the 'datetime' module in your `.zshrc` file, for example:

```bash
zmodload zsh/datetime

function preexec() {
  __TIMER=$EPOCHREALTIME
}

function powerline_precmd() {
  local __ERRCODE=$?
  local __DURATION=0

  if [ -n $__TIMER ]; then
    local __ERT=$EPOCHREALTIME
    __DURATION="$(($__ERT - ${__TIMER:-__ERT}))"
  fi

  PS1="$(powerline-go -modules duration -duration $__DURATION -error $__ERRCODE -shell zsh)"
  unset __TIMER
}
```

If the 'datetime' module is unavailable or unwanted, you may replace `$EPOCHREALTIME` with `$SECONDS`, at the loss of precision.

#### Fish

The fish prompt, in `~/.config/fish/config.fish`, will require a minimum of changes, as Fish automatically provides `$CMD_DURATION`, although with only milliseconds accuracy.

```bash
function fish_prompt
    set duration (math -s6 "$CMD_DURATION / 1000")
    $GOPATH/bin/powerline-go -modules duration -duration $duration -error $status -shell bare
end
```

## License

> This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
>
> This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see <http://www.gnu.org/licenses/>.
