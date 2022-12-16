# builder
FROM golang:1.19-alpine AS builder

WORKDIR /src

# download module
COPY go.mod /src
COPY go.sum /src
RUN --mount=target=/go/pkg,type=cache go mod download

# build applicaiton
COPY . /src

ENV GOOS=linux
ENV CGO_ENABLED=0 
# ENV GOARCH=amd64
# ENV GOAMD64=v4
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -o app .

# ---
FROM alpine:latest

WORKDIR /opt

RUN apk --no-cache add ca-certificates
COPY --from=builder /src/app /opt/app

EXPOSE 8080
ENTRYPOINT ["/opt/app"]
