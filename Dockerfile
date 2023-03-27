FROM golang:1.16-alpine

RUN mkdir /goapp
ADD . /goapp
WORKDIR /goapp
RUN go build -o main .

CMD ["/goapp/main"]