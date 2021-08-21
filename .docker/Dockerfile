# build stage
FROM golang:alpine as build
RUN mkdir /app
WORKDIR /app
COPY ./app .
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main cmd/grpc/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cli cmd/cli/worker/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o seeds cmd/cli/seeds/seeds.go
# final stage
FROM scratch
COPY --from=build /app/main /
COPY --from=build /app/cli /
COPY --from=build /app/seeds /
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
ENTRYPOINT ["/main"]