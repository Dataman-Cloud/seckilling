FROM golang:1.5

ENV GO15VENDOREXPERIMENT 1

#RUN GO15VENDOREXPERIMENT=1 go get -u github.com/FiloSottile/gvt

ENV DOCKER_VERSION 1.6.2


#Download docker
RUN set -ex; \
    curl https://get.docker.com/builds/Linux/x86_64/docker-${DOCKER_VERSION} -o /usr/local/bin/docker-${DOCKER_VERSION}; \
    chmod +x /usr/local/bin/docker-${DOCKER_VERSION}

RUN ln -s /usr/local/bin/docker-${DOCKER_VERSION} /usr/local/bin/docker

WORKDIR /go/src/github.com/Dataman-Cloud/seckilling/queue

COPY . /go/src/github.com/Dataman-Cloud/seckilling/queue

