FROM golang:alpine3.13 AS builder

WORKDIR /src
COPY . .
RUN go build -o /out/downloader1C .

FROM alpine AS app
COPY --from=builder /out/downloader1C /usr/local/bin/downloader1C
