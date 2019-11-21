FROM golang:1.13-alpine3.10 as builder
WORKDIR /go/src/github.com/leominov/gitlab-dredd
COPY . .
RUN go build -o gitlab-dredd ./

FROM alpine:3.10
COPY --from=builder /go/src/github.com/leominov/gitlab-dredd/gitlab-dredd /usr/local/bin/gitlab-dredd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["gitlab-dredd"]
