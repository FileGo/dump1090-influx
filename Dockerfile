FROM golang:alpine
RUN apk update && apk upgrade && \
    apk add --no-cache  git bash
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go build . 
CMD ["./dump1090-influx"]