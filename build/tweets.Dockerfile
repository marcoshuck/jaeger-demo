FROM golang:1.18 AS builder

WORKDIR /build

# Copying go module files and downloading dependencies before copying the codebase allow us to make better use of cache
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o app ./cmd/tweets

WORKDIR /dist
RUN cp /build/app .

FROM alpine

COPY --chown=0:0 --from=builder /dist /

USER 65534
CMD ["/app"]

