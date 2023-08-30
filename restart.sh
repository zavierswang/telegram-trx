#!/usr/bin/env bash

/bin/ps -ef | grep "telegram-trx" | grep -v "grep"
if [ "$?" -eq 1 ]; then
  echo "正在启动服务..."
  nohup ./telegram-trx >/dev/null 2>&1 &
  echo "服务启动成功"
else
  # shellcheck disable=SC2046
  # shellcheck disable=SC2006
  /bin/kill -9 `/bin/ps -ef | grep "telegram-trx" | grep -v "grep" | awk '{print $1}'`
  echo "服务已经关闭"
fi