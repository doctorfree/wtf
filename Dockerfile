FROM golang:1.19-alpine as build

ARG version=master

RUN apk add git make ncurses && \
    git clone https://github.com/doctorfree/wtf.git $GOPATH/src/github.com/doctorfree/wtf && \
    cd $GOPATH/src/github.com/doctorfree/wtf && \
    git checkout $version

ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
ENV GOSUMDB=off

WORKDIR $GOPATH/src/github.com/doctorfree/wtf

ENV PATH=$PATH:./bin

RUN make build

FROM alpine

COPY --from=build /go/src/github.com/doctorfree/wtf/bin/wtfutil /usr/local/bin/
RUN adduser -h /config -DG users -u 20000 wtf

USER wtf
ENTRYPOINT ["wtfutil"]
