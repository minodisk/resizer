FROM golang:1.8.1

WORKDIR /go/src/github.com/minodisk/resizer

RUN mkdir -p /secret && \
    go get -u \
      github.com/golang/dep/...
COPY . .
RUN go build .

CMD echo $GOOGLE_AUTH_JSON > /secret/google-auth.json && \
    resizer
