FROM golang:1.20-alpine as builder
COPY . /tmp/myService
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
WORKDIR /tmp/myService/cmd/aliddns2
RUN --mount=type=cache,target=/root/.cache/go-build go build -o aliddns2 main.go

FROM alpine
WORKDIR /usr/local/bin

EXPOSE 8080

COPY --from=builder /tmp/myService/cmd/aliddns2/config.toml /usr/local/bin/
COPY --from=builder /tmp/myService/cmd/aliddns2/aliddns2 /usr/local/bin/

ENTRYPOINT ["aliddns2"]