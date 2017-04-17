#compdef mkr

_mkr ()
{

  _arguments -C \
    '1: :->cmds' \
    '*: :->args' \
    && ret=0

  local -a list

  if ! command -v mkr > /dev/null 2>&1; then
    return 0
  fi

  case $state in

    (cmds)
      list=( ${(f)"$(mkr | sed '1,/COMMANDS/d;/GLOBAL OPTIONS/,$d' | sed 's/^ *//' | sed -E s/$'\t'+/:/ | sed 's/, h//')"} )
      _describe -t commands 'Commands' list && ret=0
      ;;

    (args)
      case $line[1] in
        (status|update|retire)
          list=( ${(f)"$(mkr hosts -f '{{range .}}{{.ID}}{{":"}}{{.Name}}{{"\n"}}{{end}}')"} )
          _describe -t plugins 'Hosts' list && ret=0
          ;;
        (*)
          _nothing
          ;;
      esac
      ;;

  esac
}

return 0
