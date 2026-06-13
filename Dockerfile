FROM golang:1.24-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY icecast_exporter.go .
RUN CGO_ENABLED=0 go build -o /icecast_exporter .

# Final stage
FROM alpine

COPY --from=build /icecast_exporter /go/bin/icecast_exporter

EXPOSE 9146
USER nobody
ENTRYPOINT ["/go/bin/icecast_exporter"]
