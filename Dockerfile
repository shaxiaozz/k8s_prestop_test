FROM golang:1.22-alpine as builder
WORKDIR /data/k8s_prestop_test-code
ENV GOPROXY=https://goproxy.cn
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache upx ca-certificates tzdata
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o k8s_prestop_test

FROM golang:1.21-alpine as runner
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /data/k8s_prestop_test-code/k8s_prestop_test /k8s_prestop_test
EXPOSE 9090
CMD ["/k8s_prestop_test"]
