FROM golang:onbuild
MAINTAINER Pure White daniel48@126.com

RUN echo "deb http://http.debian.net/debian jessie-backports main" > /etc/apt/sources.list.d/jessie-backports.list\
    && curl -sL https://deb.nodesource.com/setup_6.x | bash - \
    && apt-get install -y -t jessie-backports openjdk-8-jre-headless ca-certificates-java\
    && apt-get install -y python python-dev python-pip nodejs\
    && apt-get clean && java -version && pip install flake8

EXPOSE 48722
WORKDIR /go/src/app/linters/javascript
RUN npm install
WORKDIR /go/src/app
CMD ["go", "run", "server.go", "dispatch.go"]
