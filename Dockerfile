## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /iw

## Deploy
FROM gcr.io/distroless/base-debian10

ENV TZ="Europe/Berlin"

WORKDIR /tmp/

COPY --from=build /iw /iw

COPY html/ /tmp/html/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/iw"]