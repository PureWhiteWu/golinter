FROM golang:onbuild
MAINTAINER Pure White daniel48@126.com

# fuck gfw
RUN echo "deb http://mirrors.aliyun.com/debian/ jessie main non-free contrib" > /etc/apt/sources.list\
    && echo "deb http://mirrors.aliyun.com/debian/ jessie-proposed-updates main non-free contrib" >> /etc/apt/sources.list\
    && echo "deb-src http://mirrors.aliyun.com/debian/ jessie main non-free contrib" >> /etc/apt/sources.list\
    && echo "deb-src http://mirrors.aliyun.com/debian/ jessie-proposed-updates main non-free contrib" >> /etc/apt/sources.list\
    && mkdir ~/.pip && curl http://repo.joway.wang/pip.conf > ~/.pip/pip.conf

# add support for linters
RUN echo "deb http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main" > /etc/apt/sources.list.d/java-8-debian.list\
    && echo "deb-src http://ppa.launchpad.net/webupd8team/java/ubuntu trusty main" >> /etc/apt/sources.list.d/java-8-debian.list\
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys EEA14886\
    && echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | /usr/bin/debconf-set-selections\
    && curl -sL https://deb.nodesource.com/setup_6.x | bash - \
    && apt-get install -y oracle-java8-installer python python-dev python-pip nodejs\
    && apt-get clean && java -version && pip install flake8

RUN go get github.com/golinter/golinter
EXPOSE 48722
WORKDIR $GOPATH/src/github.com/golinter/golinter/linters/javascript
RUN npm install
WORKDIR ../..
CMD ["go", "run", "server.go", "dispatch.go"]
