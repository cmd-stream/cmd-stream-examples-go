# cmd-stream-examples-go
This repository contains several examples of using cmd-stream-go (each package 
is one example):
- echo: A minimal example.
- hello-world: Shows the basic usage of cmd-stream-go.
- hello-world_protobuf: Demonstrates the basic usage of cmd-stream-go with the 
  Protobuf serializer.
- keepalive: Shows how the client can keep a connection alive when there are no 
  commands to send.
- reconnect: Demonstrates how the client can reconnect to the server after a 
  lost connection.
- server-streaming: An example where the command sends back multiple results.
- versioning: Shows how the server can handle different versions of the same 
  command, such as supporting older clients.
- rpc: Demonstrates how to implement RPC using cmd-stream-go.
- tls: Shows how to use cmd-stream-go with TLS.

More details can be found in the corresponding `..._test.go` files.