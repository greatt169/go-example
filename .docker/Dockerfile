# build stage
FROM golang:alpine as build
RUN mkdir /app
WORKDIR /app
COPY ./app .
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o geo cmd/main.go
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o cli cmd/cli/worker/main.go
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o seeds cmd/cli/seeds/seeds.go
# final stage

RUN mkdir /psql-certs
RUN wget "https://storage.yandexcloud.net/cloud-certs/CA.pem" -O /psql-certs/psql.crt && chmod 0600 /psql-certs/psql.crt

FROM golang:alpine
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/geo /
COPY --from=build /psql-certs /

ENTRYPOINT ["/geo"]