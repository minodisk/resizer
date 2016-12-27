FROM golang:1.7.3-alpine

WORKDIR /go/src/github.com/go-microservices/resizer

RUN apk --update add \
      git && \
    go get -u \
      github.com/Masterminds/glide
COPY glide.yaml glide.yaml
COPY glide.lock glide.lock
RUN glide install
COPY . .

EXPOSE 3000

CMD go run ./main.go
