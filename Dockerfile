FROM golang:latest
ADD . /go/src/github.com/STNS/STNS
WORKDIR /go/src/github.com/STNS/STNS
RUN go get github.com/tools/godep && godep restore
CMD go test ./... && GOARCH=amd64 go build -o binary/stns.bin
