FROM golang:1.8.1

WORKDIR /go/src/github.com/minodisk/resizer

RUN go get -u \
      github.com/golang/dep/...
COPY . .
<<<<<<< HEAD
RUN dep ensure
=======
>>>>>>> 63e7b418dc82d04870bfced6809d0eec5167f27c

CMD resizer -help
