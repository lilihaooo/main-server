#!/usr/bin/env bash

# absolute path
# shellcheck disable=SC2046
# shellcheck disable=SC2164
CURRENT_PATH=$(cd `dirname $0`; pwd)
PROJECT_NAME="${CURRENT_PATH##*/}"
if [[ "${1}" = "help" ]]; then
  echo "如果不加参数，取默认值执行"
  echo "\$1参数 新项目名 default:当前目录名"
  echo "\$2参数 待替换项目名 default:go-template"
  exit
fi

if [[ "${1}" != "" ]]; then
  PROJECT_NAME="${1}"
fi

CURRENT_PROJECT_NAME="go-template"
if [[ "${2}" != "" ]]; then
  CURRENT_PROJECT_NAME="${2}"
fi

echo "current-name:"$CURRENT_PROJECT_NAME
echo "re-name:"$PROJECT_NAME

#将项目名进行修改
COMMAND="s/$CURRENT_PROJECT_NAME/"$PROJECT_NAME"/g"

sed -i "" $COMMAND cmd/root.go
sed -i "" $COMMAND cmd/service.go
sed -i "" $COMMAND main.go
sed -i "" $COMMAND go.mod