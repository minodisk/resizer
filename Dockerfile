FROM golang:1.6

RUN apt-get update -y && apt-get install git

RUN \
  wget https://github.com/Masterminds/glide/releases/download/v0.11.0/glide-v0.11.0-linux-amd64.tar.gz && \
  tar xvf glide-v0.11.0-linux-amd64.tar.gz && \
  mv linux-amd64/glide /usr/bin/

COPY . /go/src/github.com/go-microservices/resizer

WORKDIR /go/src/github.com/go-microservices/resizer

RUN go build -v

EXPOSE 3000

CMD ./resizer
