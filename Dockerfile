FROM debian:7.4

MAINTAINER Joern Weissenborn <joern.weissenborn@gmail.com>

WORKDIR zeromq-4.0.4
RUN apt-get update -y
RUN apt-get install -y curl git mercurial file build-essential libtool autoconf pkg-config net-tools

RUN curl -o /tmp/go.tar.gz https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
RUN tar -C /usr/local -zxvf /tmp/go.tar.gz
RUN rm /tmp/go.tar.gz
RUN /usr/local/go/bin/go version

ENV GOROOT /usr/local/go
ENV GOPATH /var/local/gopath
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

RUN curl -o /tmp/zeromq.tar.gz http://download.zeromq.org/zeromq-4.0.4.tar.gz
RUN tar -C /tmp -zxvf /tmp/zeromq.tar.gz
RUN rm /tmp/zeromq.tar.gz
WORKDIR /tmp/zeromq-4.0.4
RUN ./autogen.sh && ./configure && make && make install
RUN ldconfig

RUN mkdir -p $GOPATH/src
RUN mkdir -p $GOPATH/bin
RUN mkdir -p $GOPATH/pkg

RUN go get github.com/nu7hatch/gouuid
RUN go get github.com/pebbe/zmq4
RUN go get gopkg.in/yaml.v2
RUN go get gopkg.in/vmihailenco/msgpack.v2
RUN go get github.com/ugorji/go/codec

COPY . /var/local/gopath/src/github.com/joernweissenborn/aursir4go/
