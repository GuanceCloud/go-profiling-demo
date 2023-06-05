# go-profiling-demo
One simple go app for [DataKit](https://www.guance.com) continuous profiling demonstrating

## 宿主机运行

### 构建

```shell
$ cd go-profiling-demo
$ go mod tidy
$ go build
```

### 运行

```shell
$ DD_TRACE_ENABLED=true \
DD_PROFILING_ENABLED=true \
DD_AGENT_HOST=127.0.0.1 \
DD_TRACE_AGENT_PORT=9529 \
DD_SERVICE=go-profiling-demo \
DD_ENV=demo \
DD_VERSION=v0.0.1 \
./go-profiling-demo
```

### 验证运行状态

```shell
$ curl 'http://localhost:8080/movies?q=batman'
[{"title":"Batman Begins","vote_average":4.6358799822491275,"release_date":"2005-05-31"},{"title":"Batman: Mystery of the Batwoman","vote_average":3.9549411967914105,"release_date":"2003-10-21"},{"title":"Batman Beyond: Return of the Joker","vote_average":1.8787282761678148,"release_date":"2000-10-31"},{"title":"Batman \u0026 Mr. Freeze: SubZero","vote_average":2.647401476348437,"release_date":"1998-03-17"},{"title":"Batman \u0026 Robin (film)","vote_average":2.7898866857094324,"release_date":"1997-06-12"},{"title":"Batman Forever","vote_average":2.4202373224443887,"release_date":"1995-06-09"},{"title":"Batman: Mask of the Phantasm","vote_average":4.107120385998093,"release_date":"1993-12-24"},{"title":"Batman Returns","vote_average":3.592054763414077,"release_date":"1992-06-16"},{"title":"Batman (1989 film)","vote_average":4.001755941748073,"release_date":"1989-06-19"},{"title":"Batman (1966 film)","vote_average":2.189099767926333,"release_date":"1966-07-30"},{"title":"Batman Fights Dracula","vote_average":3.863042623801402,"release_date":"1900-01-01"},{"title":"Batman Dracula","vote_average":4.483630276919295,"release_date":"1900-01-01"}]
```

如果安装了 `jq` 工具，可以对返回的json内容进行格式化

```shell
$ curl 'http://127.0.0.1:8080/movies?q=spider' | jq
[
  {
    "title": "Spider in the Web",
    "vote_average": 4.3551297815393175,
    "release_date": "2019-08-30"
  },
  {
    "title": "Spider-Man 3",
    "vote_average": 2.54672384115799,
    "release_date": "2007-04-16"
  },
  {
    "title": "Spider-Man 2",
    "vote_average": 2.7380715002602,
    "release_date": "2004-06-22"
  },
  {
    "title": "Spider (2002 film)",
    "vote_average": 2.1512631396751223,
    "release_date": "2002-12-13"
  },
  {
    "title": "Spider-Man (2002 film)",
    "vote_average": 2.666549403983728,
    "release_date": "2002-04-29"
  },
  {
    "title": "Kiss of the Spider Woman (film)",
    "vote_average": 2.3350488225969306,
    "release_date": "1985-05-13"
  },
  {
    "title": "Spider Baby",
    "vote_average": 0.33910029635005945,
    "release_date": "1900-01-01"
  },
  {
    "title": "Spiderweb (film)",
    "vote_average": 3.2936595576259915,
    "release_date": "1900-01-01"
  }
]
```

## Docker 下运行

```shell
$ docker build --build-arg DK_DATAWAY=<your-dataway-endpoint> -t go-profiling-demo .
$ docker run -d go-profiling-demo
```

> DK_DATAWAY可以从观测云空间 [集成 -> Datakit](https://console.guance.com/integration/datakit) 页面上复制，例如：
> docker build --build-arg DK_DATAWAY=https://openway.guance.com?token=tkn_f5b2989ba6ab44bc988cf7e2aa4a6de3 -t go-profiling-demo .
