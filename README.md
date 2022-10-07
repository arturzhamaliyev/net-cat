# net-cat

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

## Usage

```
$ go run cmd/main.go
Listening on the port :8989
$ go run cmd/main.go 2525
Listening on the port :2525
$ go run cmd/main.go 2525 localhost
[USAGE]: ./TCPChat $port
$
```

### Note

**./TCPChat** - cmd/main.go
