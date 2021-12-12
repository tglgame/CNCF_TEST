# stage 1
FROM golang:1.17.2-alpine3.14 as builder
LABEL golang.httpserver.author=tglgame
WORKDIR /home/httpserver/
COPY server/httpserver.go .
RUN go mod init example.http || true
RUN go env -w GOPROXY="https://goproxy.cn,direct" 
RUN go mod tidy
RUN go build -o server && ls -al && pwd

# stage 2
FROM alpine
WORKDIR /home/
COPY --from=builder /home/httpserver/server .
EXPOSE 80
ENTRYPOINT ["./server"]
