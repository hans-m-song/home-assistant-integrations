FROM golang:1.22-alpine AS builder
RUN apk add --no-cache make
WORKDIR /go/src/app
COPY . .
RUN go mod download
ENV CGO_ENABLED=0
ARG BUILD_VERSION
ARG BUILD_TIME
ARG BUILD_COMMIT
RUN make build-binary \
  BINARY_NAME=/go/bin/app \
  BUILD_VERSION=${BUILD_VERSION} \
  BUILD_TIME=${BUILD_TIME} \
  BUILD_COMMIT=${BUILD_COMMIT}

FROM gcr.io/distroless/base-debian12
COPY --from=builder /go/bin/app /bin/app
CMD ["/bin/app"]
