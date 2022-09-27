# protoc-gen-go-gw-map

Plugin used to generate a map with the path and method of the gRPC Gateway request as key and the gRPC method full name (HTTP path of gRPC requests) as value.

It is highly based of the code of [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go/tree/b92717ecb630d4a4824b372bf98c729d87311a4d/cmd/protoc-gen-go).
I just got the code, stripped out any unnecessary parts, and built this plugin.

## TODO

[ ] A way to make it easy to match routes with parameters.

# Usage

Install it using the `go install` command:
> go install github.com/Hellysonrp/protoc-gen-go-gw-map

Some usage examples:
> protoc --plugin protoc-gen-go-gw-map --go-gw-map_out=output example.proto

> protoc --plugin protoc-gen-go-gw-map --go-gw-map_out=paths=source_relative:output example.proto

If you have problems with `protoc` not finding the plugin in the `PATH`, I recommend passing the absolute path to the plugin:
> protoc --plugin ${HOME}/go/bin/protoc-gen-go-gw-map --go-gw-map_out=output example.proto

> protoc --plugin ${HOME}/go/bin/protoc-gen-go-gw-map --go-gw-map_out=paths=source_relative:output example.proto
