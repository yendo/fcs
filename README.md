# fcs

fcs is a quick searcher for memos like flashcards.

## Usage

``` sh
fcs
```

You can search for the title of the memo in fzf.
The preview screen shows the contents of the memo.
Pressing enter outputs the memo to standard output.

## Installation

Install [fzf](https://github.com/junegunn/fzf) which is required to use fcs.

Copy the script `fcs-cli` to a directory with a path.

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

## Memo format

The format of the memo is similar to markdown.
However, all you really need to do is write the title of each memo in the heading
and the content below it.

``` markdown
## title1

contents1

## title2

contents2
```
