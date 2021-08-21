# build stage
FROM golang as build

RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

# NOTE: for now, this docker image always builds the current HEAD version of
# gRPC.  After gRPC's beta release, the Dockerfile versions will be updated to
# build a specific version.

# Get the source from GitHub
RUN go get google.golang.org/grpc
# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go

RUN go get github.com/githubnemo/CompileDaemon
#RUN go get github.com/grpc-ecosystem/grpc-gateway/

RUN mkdir /app

WORKDIR /app

RUN go mod init composer


RUN go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.6

RUN go install \
        github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
        github.com/golang/protobuf/protoc-gen-go
#RUN protoc -I/app -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. application/proto/services.proto && protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:.  application/proto/services.proto && protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:.  application/proto/services.proto

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /app/cmd/rest/main.go
RUN mv /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway\@v1.14.6 /go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway