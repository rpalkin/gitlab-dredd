FROM golang:1.13 as builder
WORKDIR /src
COPY . .
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /gitlab-dredd .

FROM alpine:3.10
COPY --from=builder /gitlab-dredd /go/bin/gitlab-dredd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/gitlab-dredd"]
