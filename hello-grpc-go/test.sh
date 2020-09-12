#!/bin/bash
cd "$(
  cd "$(dirname "$0")" >/dev/null 2>&1
  pwd -P
)/" || exit

export GOPATH=$GOPATH:${PWD}
go run server/proto_server.go
go run client/proto_client.go