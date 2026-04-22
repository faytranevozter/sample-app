# builder
FROM golang:1.25 AS builder

# change workdir
WORKDIR /home

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o build-app app/main.go

# final image
FROM alpine

## Copy the pre-built binary file from the previous stage
COPY --from=builder /home/build-app .
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /

ENV ZONEINFO=/zoneinfo.zip

EXPOSE 5050

ENTRYPOINT ["/build-app"]
