#!/bin/bash

# Set the path to the directory containing the .proto file
PROTO_DIR="./proto"

# Set the name of the .proto file
PROTO_FILE="example.proto"

# Set the path to the directory where you want to save the generated files
OUTPUT_DIR="./proto"

# Generate the Go language protobuf file
protoc --go_out=$OUTPUT_DIR --go_opt=paths=source_relative -I=$PROTO_DIR $PROTO_DIR/$PROTO_FILE
