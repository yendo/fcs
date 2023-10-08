# fcqs

fcqs is a quick searcher for flashcards-like notes with fzf.

## Usage

Press `Ctrl+o` (customizable) to launch fcqs on command-line.

You can search for the title of the note with fzf.
The preview screen shows the contents of the note.
The following key bindings are available.

- Enter key: Output the note to standard output.
  If the notes has shell fenced code blocks, the first block is pasted to the command-line.
- Ctrl+y: Copy the note to clip board.
- Ctrl+o: Open the first URL in the note with a browser.
- Ctrl+e: Edit the note

## Installation

Install [fzf](https://github.com/junegunn/fzf) which is required to use fcqs.

Download the fcqs archive from [GitHub Releases](https://github.com/yendo/fcqs/releases) and extract it.

Copy the command `fcqs-cli` to a directory with a path.

``` sh
install fcqs-cli ~/.local/bin/
```

Add the following function for fcqs to `~/.bashrc` for Linux, Bash and Vim users.
You have to install xclip to copy the note to clip board.
See [shell settings document](docs/shell_settings.md) for other cases.

``` bash
export VISUAL="vim"

fcqs() {
  local title=$(fcqs-cli | \
    fzf --preview "fcqs-cli {}" \
        --bind "ctrl-y:execute-silent(fcqs-cli {} | xclip -selection c),ctrl-o:execute-silent(fcqs-cli -u {} | xargs xdg-open),ctrl-e:execute-silent(fcqs-cli -l {} | awk '{printf \"+%s %s\n\",\$2,\$1}' | xargs -o $VISUAL > /dev/tty)+abort")
  fcqs-cli "$title"
  local command=$(fcqs-cli -c "$title")
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${command}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$(( READLINE_POINT + ${#command} ))
}

# You can customize the key binding
bind -x '"\C-o":fcqs'
```

## Notes specification

### File

The default notes file is `~/fcnotes.md`.
The file can be changed by the environment variable `FCQS_NOTES_FILE`.

### Format

The format of notes is similar to markdown.
However, all you really need are the titles of each note in the heading
and the content below it.

``` markdown
# title1

contents1

# title2

contents2
```

## Develop

Build the command `fcqs-cli` with Go 1.21:

``` sh
make
```

Unit test:

``` sh
make unit-test
```

Integration test:

``` sh
make integration-test
```
