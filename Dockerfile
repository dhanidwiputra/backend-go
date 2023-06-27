FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /final-project-backend

COPY . .

RUN go mod tidy

RUN go build -o binary

ENTRYPOINT ["/final-project-backend/binary"]