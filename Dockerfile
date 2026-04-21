FROM golang:1.24-alpine AS build

WORKDIR /src

RUN apk add --no-cache git

COPY icecast_exporter.go .

RUN go mod init github.com/DonkeyMobile/kdgm-icecast_exporter \
  && go mod tidy \
  && CGO_ENABLED=0 go build -o /icecast_exporter .

# Final stage
FROM alpine

COPY --from=build /icecast_exporter /icecast_exporter

EXPOSE 9146
USER nobody
ENTRYPOINT ["/icecast_exporter"]
