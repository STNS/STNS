FROM golang:latest as builder
ADD . /opt/stns
WORKDIR /opt/stns/
RUN GOOS=linux CGO_ENABLED=0 make build

FROM scratch
COPY --from=builder /opt/stns/tmp/bin/stns /bin/stns
COPY --from=builder /opt/stns/tmp/bin/mod_stns_etcd.so /usr/local/stns/modules.d
COPY --from=builder /opt/stns/tmp/bin/mod_stns_dynamodb.so /usr/local/stns/modules.d
COPY misc/docker.conf /etc/stns/server/stns.conf
EXPOSE 1104
CMD ["/bin/stns"]
