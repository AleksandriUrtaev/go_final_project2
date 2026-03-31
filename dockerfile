FROM ubuntu:latest

WORKDIR /app

COPY scheduler .

COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

EXPOSE 7540

CMD ["./scheduler"]