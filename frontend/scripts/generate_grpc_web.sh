#!/bin/bash

PROTO_DIR=../api/proto
OUT_DIR=./src/grpc

protoc -I=$PROTO_DIR users.proto \
  --js_out=import_style=commonjs:$OUT_DIR \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:$OUT_DIR 