FROM golang:1.17 AS build-env
WORKDIR /app
ADD . /app/
RUN go get -d -v ./...
RUN go build -o /go/bin/app

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/app /
CMD ["/app"]