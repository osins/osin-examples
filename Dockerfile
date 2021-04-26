FROM golang:alpine AS builder
# FROM golang:1.16.3-buster AS builder
ENV http_proxy=http://192.168.3.32:1080
ENV https_proxy=http://192.168.3.32:1080
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS="linux"
RUN go env -w CGO_LDFLAGS="-w -s"
RUN go env

COPY ./src /app
RUN apk --no-cache add tzdata && cd /app && go build -a -installsuffix cgo -o serve main.go
RUN ls /app

# FROM scratch
FROM alpine
# FROM golang:1.16.3-buster
COPY --from=builder /app/serve /app/
COPY --from=builder /app/template /app/template
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
ENV TZ=Asia/Shanghai

WORKDIR /
EXPOSE 80

CMD ["/app/serve"]