FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*


RUN wget https://golang.org/dl/go1.24.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.5.linux-amd64.tar.gz && \
    rm go1.24.5.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /app

COPY . .

ENV TODO_PORT=7540
ENV TODO_PASSWORD=flisthdo

EXPOSE 7540

CMD ["go", "run", "./pkg/server/main.go"]


