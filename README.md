# cmd-stream-examples-go
Demonstrates several examples of using the [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go)
with the [mus-stream-go](https://github.com/mus-format/mus-stream-go) 
serializer. It consists of the following packages (each package is one example):
- standard - demonstrates the standard use of the cmd-stream-go.
- multi-result - demonstrates an example where the command returns multiple 
  results.
- max-cmd-size - shows how you can set the maximum command size supported by 
  the server.
- versioning - shows how you can use different versions of the same command, for
  example, to support old clients.
- keepalive - demonstrates the keepalive feature.
- reconnect - demonstrates the reconnect feature.
- rpc - shows how the cmd-stream-go can be used to implement RPC.
- tls - cmd-stream-go + TLS.

You can find more information in the corresponding `main.go` files.