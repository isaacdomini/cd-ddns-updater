FROM golang:1.17.7 AS GO_BUILD
ENV CGO_ENABLED 0
COPY . /app
WORKDIR /app
RUN go build -o server

FROM alpine:3.15
RUN mkdir -p /data
WORKDIR /app
COPY --from=GO_BUILD /app/server /app/server
EXPOSE 8080
CMD ["./server"]
