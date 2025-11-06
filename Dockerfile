FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .

ENV TODO_PORT=7540
ENV TODO_DBFILE="/data/scheduler.db"
ENV TODO_PASSWORD=flisthdo

EXPOSE 7540

CMD ["go", "run", "./pkg/server/main.go"]