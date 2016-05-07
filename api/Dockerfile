FROM golang:1.4

WORKDIR /go/src/github.com/josemrobles/data-ingestion-poc/api
EXPOSE 1337

RUN apt-get update &&\
    apt-get install -y supervisor ;\
    mkdir -p /var/log/supervisor 

COPY supervisor.conf /etc/supervisor/conf.d/supervisor.conf
ADD . /go/src/github.com/josemrobles/data-ingestion-poc/api
CMD "./entrypoint.sh"
