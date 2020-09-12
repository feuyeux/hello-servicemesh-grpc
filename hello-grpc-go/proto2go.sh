#!/bin/bash
cd "$(
  cd "$(dirname "$0")" >/dev/null 2>&1
  pwd -P
)/" || exit
export PATH=$PATH:$GOPATH/bin
proto_path=$(pwd)/proto
go_proto_path=$(pwd)

if [ ! -d "${go_proto_path}" ]; then
  mkdir -p "${go_proto_path}"
else
  rm -f "${go_proto_path}/*"
fi

echo "protoc --proto_path=${proto_path} --go_out=plugins=grpc:${go_proto_path} ${proto_path}/landing.proto"
protoc --proto_path="${proto_path}" --go_out=plugins=grpc:"${go_proto_path}" "${proto_path}"/landing.proto
