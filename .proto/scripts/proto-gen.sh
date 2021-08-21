#!/bin/bash

printLine() {
  echo "\e[32m$1\e[0m"
}

printLine '===============================================';
printLine "Generate protobuf code";
printLine '===============================================';

cd /

#printLine 'Generate protobuf code';
protoc -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. news-ms/app/interfaces/grpc/proto/v1/news/*.proto && protoc -I/usr/local/include -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:.  news-ms/app/interfaces/grpc/proto/v1/news/*.proto && protoc -I/usr/local/include -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:.  news-ms/app/interfaces/grpc/proto/v1/news/*.proto
protoc -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. news-ms/app/interfaces/grpc/proto/v1/access_control/access_control.proto && protoc -I/usr/local/include -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:.  news-ms/app/interfaces/grpc/proto/v1/access_control/access_control.proto && protoc -I/usr/local/include -I. -I$GOPATH/pkg -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:.  news-ms/app/interfaces/grpc/proto/v1/access_control/access_control.proto
echo 'Done!'