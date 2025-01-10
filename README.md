# cmd-stream-examples-go
Contains several examples of using [cmd-stream-go](https://github.com/cmd-stream/cmd-stream-go)
(each package is one example):
- hello-world: demonstrates the standard usage of the cmd-stream-go library.
- hello-world_protobuf: demonstrates the standard usage of the cmd-stream-go 
  library with the Protobuf serializer.
- keepalive: demonstrates how the client can keep a connection alive when there  
  are no commands to send. 
- reconnect: demonstrates how the client can reconnect to the server when the  
  connection is lost.
- server-streaming: an example where the command sends back multiple results.  
- versioning: demonstrates how the server can handle different versions of the  
  same command, for example, to support older clients.  
- rpc - demonstrates how to implement RPC with cmd-stream-go.
- tls - cmd-stream-go + TLS.

More information can be found in the corresponding `main.go` files.