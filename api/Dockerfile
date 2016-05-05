FROM golang:1.4

WORKDIR /go/src/github.com/josemrobles/data-ingestion-poc/app
EXPOSE 8080

RUN apt-get update &&\
    apt-get install -y supervisor ;\
    mkdir -p /var/log/supervisor 

COPY supervisor.conf /etc/supervisor/conf.d/supervisor.conf
ADD . /go/src/github.com/josemrobles/data-ingestion-poc/app
CMD "./entrypoint.sh"
