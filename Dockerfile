#
# Build container
#
FROM golang:1.11 as build
WORKDIR /go/src/solarwinds/k8slessons
RUN go get -u github.com/golang/dep/cmd/dep
ADD Gopkg.toml .
ADD Gopkg.lock .
RUN dep ensure --vendor-only
ADD . .
RUN GOOS=linux go build -a -o /bin/dumbstore

#
# Run container
#
FROM alpine:latest
# https://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/ca-certificates.crt
COPY --from=build /bin/dumbstore dumbstore
ENTRYPOINT [ "./dumbstore"]
