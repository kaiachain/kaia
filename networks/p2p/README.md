# p2p stack

## Interfaces and implementations

```go
// The top-level p2p server that commands all the other components.

type Config struct // server.go
type Server interface // server.go
type BaseServer struct implements Server // server_base.go
type SingleChannelServer struct extends BaseServer // server_base.go
type MultiChannelServer struct extends BaseServer // server_multi.go

// Tools to establish connections to other nodes.

type NodeDialer interface // server_util.go
type TCPDialer struct implements NodeDialer // server_util.go

type dialer interface // dial.go
type dialstate struct implements dialer // dial.go

// Represents connections to other nodes.

type transport interface // transport.go
type rlpxTransport struct implements transport // transport.go

type MsgReadWriter interface // message.go
type conn struct implements MsgReadWriter // server_util.go

type Peer interface // peer.go
type singleChannelPeer struct implements Peer // peer.go
type multiChannelPeer struct implements Peer // peer.go
```
