FROM golang:1.21.5-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod download
ENV CGO_ENABLED=0
RUN go build -o /go/bin/app .

FROM gcr.io/distroless/base-debian11
COPY --from=builder /go/bin/app /
CMD ["/app"]
