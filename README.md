# fcs

fcs is a quick searcher for flashcards-like notes with fzf.

## Usage

``` sh
fcs
```

You can search for the title of the note with fzf.
The preview screen shows the contents of the note.
Pressing enter outputs the note to standard output.

## Installation

Install [fzf](https://github.com/junegunn/fzf) which is required to use fcs.

Build the command `fcs-cli` with Go 1.21.

``` sh
go build -o fcs-cli
```

Copy the command to a directory with a path.

``` sh
install fcs-cli ~/.local/bin/
```

Add the following function for fcs to `~/.bashrc`.

``` bash
fcs() {
  index=$(fcs-cli | fzf --preview "fcs-cli {}") &&
    fcs-cli "$index"
}
```

## Notes specification

### File

The default notes file is `~/fcnotes.md`.
The file can be changed by the environment variable `FCS_NOTES_FILE`.

### Format

The format of notes is similar to markdown.
However, all you really need are the titles of each note in the heading
and the content below it.

``` markdown
## title1

contents1

## title2

contents2
```
