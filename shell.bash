# fcqs

# You can customize settings:
#
# FCQS_EDITOR="default" or "vscode"
# FCQS_COPY_KEY="ctrl-y"
# FCQS_OPEN_KEY="ctrl-o"
# FCQS_EDIT_KEY="ctrl-e"
# FCQS_BASH_BIND_KEY="\C-o"
# FCQS_COPY_COMMAND="xclip -selection c"
# FCQS_COPY_WITH_TITLE=true
# FCQS_OPEN_COMMAND="open"

FCQS_EDITOR=${FCQS_EDITOR:-default}
FCQS_COPY_KEY=${FCQS_COPY_KEY:-ctrl-y}
FCQS_OPEN_KEY=${FCQS_OPEN_KEY:-ctrl-o}
FCQS_EDIT_KEY=${FCQS_EDIT_KEY:-ctrl-e}
FCQS_BASH_BIND_KEY=${FCQS_BASH_BIND_KEY:-"\C-o"}
FCQS_COPY_COMMAND=${FCQS_COPY_COMMAND:-"xclip -selection c"}
FCQS_COPY_WITH_TITLE=${FCQS_COPY_WITH_TITLE:-true}
FCQS_OPEN_COMMAND=${FCQS_BROWSE_COMMAND:-"open"}

FCQS_EDIT_COMMAND_DEFAULT="awk '{printf \"+%s %s\n\",\$2,\$1}' | xargs -o ${VISUAL} > /dev/tty"
FCQS_EDIT_COMMAND_VSCODE="awk '{printf \"%s:%s\n\",\$1,\$2}' | xargs -o code -g"
[ "${FCQS_EDITOR}" = "vscode" ] && FCQS_EDIT_COMMAND=${FCQS_EDIT_COMMAND_VSCODE} || FCQS_EDIT_COMMAND=${FCQS_EDIT_COMMAND_DEFAULT}

[ "${FCQS_COPY_WITH_TITLE}" = true ] && FCQS_COPY_COMMAND_FLAG="" || FCQS_COPY_COMMAND_FLAG="-t"

fcqs() {
  local title
  title=$(fcqs-cli |
    fzf --preview "fcqs-cli {}" \
      --bind "${FCQS_COPY_KEY}:execute-silent(fcqs-cli ${FCQS_COPY_COMMAND_FLAG} {} | ${FCQS_COPY_COMMAND}),${FCQS_OPEN_KEY}:execute-silent(fcqs-cli -u {} | xargs ${FCQS_OPEN_COMMAND}),${FCQS_EDIT_KEY}:execute-silent(fcqs-cli -l {} | ${FCQS_EDIT_COMMAND})+abort")
  fcqs-cli "$title"
  local command
  command=$(fcqs-cli -c "$title")
  READLINE_LINE="${READLINE_LINE:0:$READLINE_POINT}${command}${READLINE_LINE:$READLINE_POINT}"
  READLINE_POINT=$((READLINE_POINT + ${#command}))
}

eval "bind -x '\"${FCQS_BASH_BIND_KEY}\":fcqs'"
