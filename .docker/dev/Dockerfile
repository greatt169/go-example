# build stage
FROM golang:alpine as build

RUN mkdir /app
WORKDIR /app
COPY ./app .

RUN go clean --modcache
RUN go get github.com/githubnemo/CompileDaemon
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -cover ./...
RUN go build /app/cmd/main.go

EXPOSE 8080
ENTRYPOINT CompileDaemon --build="go build /app/cmd/main.go" --command=./main