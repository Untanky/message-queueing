# Outward API

## Design Decisions

There are two common patterns when implementing APIs that come to mind for
this project: (i) RPC (specifically gRPC) and (ii) RESTful HTTP API. Many
points can be made for and against each API pattern.

gRPC is particularly efficient when using streams. When it comes to Unary-RPC
calls, gRPC is much slower than an HTTP call, yet when it comes to streaming
data, it is much more efficient. With Streamed gRPC, overhead for connection
establishing and headers is only occurred once. 

REST over HTTP is a much wider implemented approach, especially for user-facing
applications. HTTP can much more easily be load-balanced and distributed among
many nodes. gRPC's streaming is a great strength but the continuous connection
is much more resource intensive and harder to account for in traffic.

**_Decision Nov-2023:_** RESTful API
