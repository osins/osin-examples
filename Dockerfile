FROM golang:alpine AS builder
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS="linux"
RUN go env -w CGO_LDFLAGS="-w -s"
RUN go env

COPY ./src /tmp/sso
RUN cd /tmp/sso && go mod vendor && go build -a -installsuffix cgo -o serve main.go
RUN ls /tmp/sso

FROM scratch
COPY --from=builder /tmp/sso/serve /

WORKDIR /
EXPOSE 80

CMD ["/serve"]