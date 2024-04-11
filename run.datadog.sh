#!/usr/bin/env bash

DD_TRACE_ENABLED=true DD_PROFILING_ENABLED=true DD_AGENT_HOST=127.0.0.1 DD_TRACE_AGENT_PORT=8126 DD_SERVICE=go-profiling-demo DD_ENV=demo DD_VERSION=v0.0.1 ./go-profiling-demo
