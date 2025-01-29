# cmd-stream-examples-go
This repository contains several examples of using cmd-stream-go (each package 
is one example):
- echo: A minimal example.
- hello-world: Shows the basic usage of cmd-stream-go.
- hello-world_protobuf: Demonstrates the basic usage of cmd-stream-go with the 
  Protobuf serializer.
- keepalive: Shows how the client can keep a connection alive when there are no 
  Commands to send.
- reconnect: Demonstrates how the client can reconnect to the server after
  losing the connection.
- server-streaming: An example where the Command sends back multiple Results.
- versioning: Shows how the server can handle different versions of the same 
  Command.
- rpc: Demonstrates how to implement RPC using cmd-stream-go.
- tls: Shows how to use cmd-stream-go with TLS.

More details can be found in the corresponding `..._test.go` files.