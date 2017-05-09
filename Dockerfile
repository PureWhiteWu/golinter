FROM golang:onbuild
MAINTAINER Pure White daniel48@126.com

RUN echo "deb http://http.debian.net/debian jessie-backports main" > /etc/apt/sources.list.d/jessie-backports.list\
    && curl -sL https://deb.nodesource.com/setup_6.x | bash - \
    && apt-get install -t jessie-backports  openjdk-8-jre-headless ca-certificates-java\
    && apt-get install -y python python-dev python-pip nodejs\
    && apt-get clean && java -version && pip install flake8

RUN go get github.com/golinter/golinter
EXPOSE 48722
WORKDIR $GOPATH/src/github.com/golinter/golinter/linters/javascript
RUN npm install
WORKDIR ../..
CMD ["go", "run", "server.go", "dispatch.go"]
