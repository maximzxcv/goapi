FROM golang:1.16.0-alpine AS build_base


ENV APP_HOME /src
WORKDIR ${APP_HOME}
 

COPY go.mod .
COPY go.sum .
COPY bal bal
COPY cmd cmd

 RUN go mod download 
 RUN go mod vendor

 RUN go build ./cmd/api  download 

EXPOSE 8080

# Run the binary program produced by `go install`
ENTRYPOINT ["./api"]
 