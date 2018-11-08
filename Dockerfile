#
# Build container
#
FROM golang:1.11 as build
WORKDIR /go/src/solarwinds/k8slessons
ADD . .
RUN GOOS=linux go build -a -o /bin/dumbstore

#
# Run container
#
FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/ca-certificates.crt
COPY --from=build /bin/dumbstore /dumbstore
ENTRYPOINT [ "./dumbstore"]
