FROM golang:onbuild
MAINTAINER Pure White daniel48@126.com

# fuck gfw
RUN curl http://mirrors.163.com/.help/sources.list.jessie > /etc/apt/sources.list\
    && mkdir ~/.pip && curl http://repo.joway.wang/pip.conf > ~/.pip/pip.conf

# add support for linters
RUN curl -sL https://deb.nodesource.com/setup_6.x | bash - \
    && echo "deb http://httpredir.debian.org/debian/ jessie main contrib" >> /etc/apt/sources.list\
    && apt-get update && apt-get install -y java-package\
    && wget --no-check-certificate --no-cookies --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u66-b17/jdk-8u66-linux-x64.tar.gz\
    && make-jpkg jdk-8u66-linux-x64.tar.gz && dpkg -i oracle-java8-jdk_8u66_amd64.deb && rm -rf jdk-8u66-linux-x64.tar.gz oracle-java8-jdk_8u66_amd64.deb \
    && echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | sudo /usr/bin/debconf-set-selections \
    && apt-get install -y oracle-java8-installer python python-dev python-pip nodejs build-essential\
    && apt-get clean && java -version && pip install flake8

# RUN go install github.com/golinter/golinter
EXPOSE 48722
WORKDIR src/github.com/golinter/golinter/linters/javascript
RUN echo '\n#alias for cnpm\nalias cnpm="npm --registry=https://registry.npm.taobao.org \
      	  --cache=$HOME/.npm/.cache/cnpm \
      	  --disturl=https://npm.taobao.org/dist \
      	  --userconfig=$HOME/.cnpmrc"' >> xshrc && source xshrc && npm install
WORKDIR ../..
CMD ["go", "run", "server.go", "dispatch.go"]
