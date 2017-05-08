FROM golang:onbuild
MAINTAINER Pure White daniel48@126.com

# add support for linters
RUN apt-get install software-properties-common && add-apt-repository -y ppa:webupd8team/java && apt-get update \
    && echo oracle-java8-installer shared/accepted-oracle-license-v1-1 select true | sudo /usr/bin/debconf-set-selections \
    && apt-get install -y oracle-java8-installer python python-dev python-pip\
    && apt-get clean && java -version
RUN pip install flake8
RUN go install github.com/golinter/golinter
EXPOSE 48722
WORKDIR src/github.com/golinter/golinter
CMD ["go", "run", "server.go", "dispatch.go"]
