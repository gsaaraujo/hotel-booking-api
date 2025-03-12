FROM ubuntu:24.04

RUN apt-get update && apt-get install -y \
  git \
  curl \
  && apt-get clean autoclean

RUN curl -OL https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && \
  tar -C /usr/local -xzvf go1.24.1.linux-amd64.tar.gz && \
  rm go1.24.1.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOBIN=/usr/local/bin

RUN go install github.com/jackc/tern/v2@latest

ENV PATH=$PATH:$GOBIN
