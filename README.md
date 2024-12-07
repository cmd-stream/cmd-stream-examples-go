# cmd-stream-examples-go
Contains several examples of using [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go)
(each package is one example):
- standard - demonstrates the standard use of cmd-stream-go with 
  [mus-stream-go](https://github.com/mus-format/mus-stream-go) serializer.
- standard_protobuf - demonstrates the standard use of cmd-stream-go with 
  Protobuf serializer.
- rpc - shows how cmd-stream-go can be used to implement RPC.
- multi_result - demonstrates an example where the command returns multiple 
  results.
- keepalive - demonstrates the keepalive feature.
- reconnect - demonstrates the reconnect feature.
- max_cmd_size - shows how to set the maximum command size supported by the 
  server.
- cmd_versioning - shows how the server can handle different versions of the 
  same command, for example, to support old clients.
- tls - cmd-stream-go + TLS.

More information can be found in the corresponding `main.go` files.