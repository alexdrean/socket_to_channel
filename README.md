# socket-to-channel
## Usage
```go
import (
    "github.com/alexdrean/socket_to_channel"
    "net"
)
...
addr := net.Dial("tcp", "192.168.1.1:8080")
receive, transmit := socket_to_channel.DialToChannel(addr, 100 * time.Millisecond) // retry every 100ms
transmit("Hello world!")
for s := receive {
    fmt.Printf("Received: %s\n", s)
}
```