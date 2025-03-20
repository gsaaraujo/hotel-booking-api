FROM mcr.microsoft.com/devcontainers/base:ubuntu

RUN apt-get update && \
  apt-get clean autoclean

RUN curl -OL https://go.dev/dl/go1.24.1.linux-amd64.tar.gz && \
  tar -C /usr/local -xzvf go1.24.1.linux-amd64.tar.gz && \
  rm go1.24.1.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOBIN=/usr/local/bin

RUN go install github.com/jackc/tern/v2@latest

RUN wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg && \
  echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list && \
  sudo apt update && sudo apt install terraform

ENV API_ENV=LOCAL