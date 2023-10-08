# Shell Settings for fcqs

Add the following function for fcqs to the shell script setting file.

## Linux and bash

You have to install xclip to copy the note to clip board.

for Unix editors (Vim, Emacs, nano, gedit, etc.):

```bash
export VISUAL="vim" # your editor

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

for Visual Studio Code:

```bash
fcqs() {
  local title=$(fcqs-cli | \
    fzf --preview "fcqs-cli {}" \
        --bind "ctrl-y:execute-silent(fcqs-cli {} | xclip -selection c),ctrl-o:execute-silent(fcqs-cli -u {} | xargs xdg-open),ctrl-e:execute-silent(fcqs-cli -l {} | awk '{printf \"%s:%s\n\",\$1,\$2}' | xargs -o code -g)+abort")
  fcqs-cli "$title"
  local command=$(fcqs-cli -c "$title")
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${command}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$(( READLINE_POINT + ${#command} ))
}

# You can customize the key binding
bind -x '"\C-o":fcqs'
```
