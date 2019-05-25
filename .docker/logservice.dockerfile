FROM ubuntu:16.04
RUN apt-get update && apt-get install -y golang
COPY backend /go/src/github.com/IAmRDhar/scaling-web-app/backend
COPY *.pem /
ENV GOPATH /go
EXPOSE 6000
RUN go install github.com/IAmRDhar/scaling-web-app/backend/logservice
WORKDIR /go
ENTRYPOINT ["./bin/logservice"]
