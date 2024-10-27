FROM golang:1.22-alpine
LABEL authors="aemustaf"
WORKDIR /app

COPY COPY . ./
RUN go mod download
RUN go build

EXPOSE 8080

ENTRYPOINT ["./to-do-list"]