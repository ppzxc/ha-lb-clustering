FROM golang:1.16.2

USER root
RUN mkdir /client
COPY ./client /client
COPY ./entrypoint.sh /client/entrypoint.sh
WORKDIR /client

RUN go build -o ./client ./
RUN ls /client
RUN chmod +x ./client

EXPOSE 5141/udp

ENTRYPOINT bash entrypoint.sh
