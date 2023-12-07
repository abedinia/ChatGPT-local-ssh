FROM golang:1.21.4

WORKDIR /app

COPY . .

COPY ssh/id_rsa /root/.ssh/id_rsa
COPY ssh/id_rsa.pub /root/.ssh/id_rsa.pub
COPY ssh/known_hosts /root/.ssh/known_hosts

RUN chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub && \
    chmod 644 /root/.ssh/known_hosts

RUN go build main.go

CMD ["./main"]