FROM golang:1.17-alpine AS builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build -o server ./api

FROM alpine
WORKDIR /app
COPY --from=builder /build/server /app/
COPY --from=builder /build/public /app/
WORKDIR /app
CMD ["./server"]