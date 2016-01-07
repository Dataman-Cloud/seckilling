FROM golang:1.5.1

ENV GO15VENDOREXPERIMENT 1

#RUN GO15VENDOREXPERIMENT=1 go get -u github.com/FiloSottile/gvt

WORKDIR /go/src/github.com/Dataman-Cloud/seckilling/gate

COPY . /go/src/github.com/Dataman-Cloud/seckilling/gate

