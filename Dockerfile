#
# docker buildx build --platform linux/amd64,linux/arm64 -t registry.cn-hangzhou.aliyuncs.com/zhangyicloud/go-profiling-demo:0.4.0 --push .

FROM golang:1.23-bullseye as builder

COPY . /build/

RUN go env -w GOPROXY="https://goproxy.cn,direct" && cd /build && go build



FROM debian:bullseye
LABEL authors="guance.com" email="zhangyi905@guance.com"

ENV DD_SERVICE go-profiling-demo
ENV DD_VERSION v0.1.0
ENV DD_ENV testing
ENV DD_AGENT_HOST 127.0.0.1
ENV DD_TRACE_AGENT_PORT 9529
ENV DD_TRACE_ENABLED true
ENV DD_PROFILING_ENABLED true

WORKDIR /usr/local/go-profiling-demo

COPY --from=builder /build/go-profiling-demo ./
COPY movies5000.json.gz run.sh run.datadog.sh ./

CMD ./go-profiling-demo
