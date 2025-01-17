FROM golang:1.22.1

WORKDIR /usr/src/app

COPY . .
RUN go mod tidy

EXPOSE 9007
