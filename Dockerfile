# syntax=docker/dockerfile:1

FROM golang:1.17

RUN mkdir /goapp
ADD . /goapp
WORKDIR /goapp

ENV AUTH_TOKEN=
ENV APP_TOKEN=
ENV CHANNEL_ID=
ENV CLIENT_ID=
ENV TOKEN=

RUN export GO111MODULE=on
RUN cd /goapp && git clone https://github.com/ShivaniKumar1/gobot.git

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o main .

CMD ["/goapp/main"]
