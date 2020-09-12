#!/bin/bash
cd "$(
  cd "$(dirname "$0")" >/dev/null 2>&1
  pwd -P
)/" || exit
protoDir=$(pwd)/proto
## py pb path
pyProtoDir=pb
if [ ! -d "${pyProtoDir}" ]; then
  mkdir -p ${pyProtoDir}
fi

## https://developers.google.com/protocol-buffers/docs/reference/python-generated
## *_pb2.py which contains our generated request and response classes
## *_pb2_grpc.py which contains our generated client and server classes.
python3 -m grpc.tools.protoc \
  -I ${protoDir} \
  --python_out=${pyProtoDir} \
  --grpc_python_out=${pyProtoDir} \
  ${protoDir}/landing.proto
