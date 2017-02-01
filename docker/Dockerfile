FROM golang:latest
MAINTAINER G Manoj <manoj@ionosnetworks.com>
RUN mkdir -p /cpcxchng
WORKDIR /cpcxchng
ADD . /cpcxchng
EXPOSE 3000
CMD ["go","run","servr/server.go"]
