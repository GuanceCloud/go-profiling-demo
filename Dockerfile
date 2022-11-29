FROM golang:1.18

MAINTAINER zhangyi905@guance.com

COPY . /usr/local/go-profiling-demo
WORKDIR /usr/local/go-profiling-demo

RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN go build

ARG DK_DATAWAY=https://openway.guance.com?token=tkn_f5b2989ba6ab44bc988cf7e2aa4a6de3
RUN DK_DATAWAY=${DK_DATAWAY} bash -c "$(curl -L https://static.guance.com/datakit/install.sh)"
RUN cp /usr/local/datakit/conf.d/profile/profile.conf.sample /usr/local/datakit/conf.d/profile/profile.conf

CMD /bin/sh run.sh