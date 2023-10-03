# Shell Settings for fcs

Add the following function for fcs to the shell script setting file.

## Linux and bash

You have to install xclip to copy the note to clip board.

for Vim, Emacs, gedit, etc.

```bash
export VISUAL="vim" # your editor

fcs() {
  local title=$(fcs-cli | \
    fzf --preview "fcs-cli {}" \
        --bind "ctrl-y:execute-silent(fcs-cli {} | xclip -selection c),ctrl-o:execute-silent(fcs-cli -u {} | xargs xdg-open),ctrl-e:execute-silent(fcs-cli -l {} | awk '{printf \"+%s %s\n\",\$2,\$1}' | xargs -o $VISUAL > /dev/tty)+abort")
  fcs-cli "$title"
  local command=$(fcs-cli -c "$title")
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${command}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$(( READLINE_POINT + ${#command} ))
}

# You can customize the key binding
bind -x '"\C-o":fcs'
```

for Visual Studio Code:

```bash
fcs() {
  local title=$(fcs-cli | \
    fzf --preview "fcs-cli {}" \
        --bind "ctrl-y:execute-silent(fcs-cli {} | xclip -selection c),ctrl-o:execute-silent(fcs-cli -u {} | xargs xdg-open),ctrl-e:execute-silent(fcs-cli -l {} | awk '{printf \"+%s %s\n\",\$2,\$1}' | xargs -o $VISUAL > /dev/tty)+abort")
  fcs-cli "$title"
  local command=$(fcs-cli -c "$title")
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${command}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$(( READLINE_POINT + ${#command} ))
}

# You can customize the key binding
bind -x '"\C-o":fcs'
```
