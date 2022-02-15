FROM golang:1.17.7 AS GO_BUILD
ENV CGO_ENABLED 0
COPY . /app
WORKDIR /app
RUN go build -o server

FROM ubuntu:latest
RUN apt-get update && apt-get install -y curl unzip
WORKDIR /
RUN curl -L https://github.com/mkorthof/freenom-script/archive/refs/heads/master.zip -o freenom.zip
RUN unzip freenom.zip -d freenom
RUN sed -i 's/#   source "\/home\/${LOGNAME}\/.secret\/.freenom"/source "\/data\/.freenom"/g' /freenom/freenom-script-master/freenom.conf
RUN mkdir -p /data
WORKDIR /app
COPY --from=GO_BUILD /app/server /app/server
EXPOSE 8080
CMD ["./server"]
