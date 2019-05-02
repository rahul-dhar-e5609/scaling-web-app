FROM ubuntu:16.04
RUN apt-get update && apt-get install -y golang
COPY backend /go/src/github.com/IAmRDhar/scaling-web-app/backend
COPY *.pem /
COPY web /web
ENV GOPATH /go
ENV DATA_SERVICE_URL https://172.18.0:4000
EXPOSE 3000
RUN go install github.com/IAmRDhar/scaling-web-app/backend/web
WORKDIR /go
ENTRYPOINT ["./bin/web"]
