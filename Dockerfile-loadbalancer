FROM ubuntu:16.04
RUN apt-get update && apt-get install -y golang
COPY backend /go/src/github.com/IAmRDhar/scaling-web-app/backend
COPY *.pem /
ENV GOPATH /go
EXPOSE 2000 2001
RUN go install github.com/IAmRDhar/scaling-web-app/backend/loadbalancer
WORKDIR /go
ENTRYPOINT ["./bin/loadbalancer"]