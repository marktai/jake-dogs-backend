FROM golang:alpine

RUN mkdir /app
WORKDIR /app
ENV GOPATH /app
ADD . /app/
RUN /bin/sh -c "cd /app/src/main && go get"
RUN go build -o /app/bin/main src/main/main.go
CMD ["/app/bin/main"]
