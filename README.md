# cmd-stream-examples-go
Contains several examples of using [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go)
with the [mus-stream-go](https://github.com/mus-format/mus-stream-go) 
serializer (each package is one example):
- standard - demonstrates the standard use of cmd-stream-go.
- rpc - shows how cmd-stream-go can be used to implement RPC.
- multi_result - demonstrates an example where the command returns multiple 
  results.
- keepalive - demonstrates the keepalive feature.
- reconnect - demonstrates the reconnect feature.
- max_cmd_size - shows how to set the maximum command size supported by the 
  server.
- data_versioning - shows how to use different versions of the same command, for
  example, to support old clients.
- tls - cmd-stream-go + TLS.

More information can be found in the corresponding `main.go` files.