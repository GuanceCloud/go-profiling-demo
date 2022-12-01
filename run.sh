#!/usr/bin/env bash

nohup /usr/local/datakit/datakit >/dev/null 2>&1 &
if [ $? -ne 0 ]; then
  echo "start datakit failed :("
  exit 1
fi


DD_AGENT_HOST=127.0.0.1 DD_TRACE_AGENT_PORT=9529 DD_SERVICE=go-profiling-demo DD_ENV=demo DD_VERSION=v0.0.1 ./go-profiling-demo