FROM scratch 
MAINTAINER G Manoj <manoj@ionosnetworks.com>
RUN mkdir -p /cpcxchng
WORKDIR /cpcxchng
COPY /server /cpcxchng/
EXPOSE 3000
CMD ["/server"]
