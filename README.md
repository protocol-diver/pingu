# PINGU [![GoDoc](https://godoc.org/github.com/protocol-diver/pingu?status.svg)](https://godoc.org/github.com/protocol-diver/pingu)

This is a simple and tiny heartbeat library with UDP. <br>
Heartbeat is a message that checks if you are healthy or if other connected peers are healthy. The core feature is to periodically request a heartbeat and send an 'alive' message when a heartbeat request is received. Also need to surface a method that allows the user to check the health status. <br>
It could make it your building easier when you need build a server with heartbeat.<br>
If you need heartbeat communication with an external network among peers, register the public IP obtained through a process such as hole punching(<b>This library does not support communication across NAT</b>).

## Install
```
$ go get -u github.com/protocol-diver/pingu
```

## Usage

### Create your Pingu
```go
myPingu, err := pingu.NewPingu("127.0.0.1:4874", nil)

// You could preconfig like below,
pingu.Config{
  RecvBufferSize: 512, // default value : 256
  Verbose: true // It's notify that what's going on, default value : false
}
```

### Embed into your Server
```go
type Server struct {
  conn  *net.TCPConn
  pingu *pingu.Pingu
}
```

### Register other Pingus
```go
if err := myPingu.RegisterWithRawAddr("127.0.0.1:4875"); err != nil {
  return err
}
if err := myPingu.RegisterWithRawAddr("127.0.0.1:4876"); err != nil {
  return err
}
```

### Let work Pingu
```go
myPingu.Start()

ticker := time.Ticker(5 * time.Second)
// Second param is time limit of send ping -> recv pong.
cancel, err := myPingu.BroadcastPingWithTicker(*ticker, 3*time.Second)
if err != nil {
  return err
}
```

### Watch Pingu Working
```go
// It's returns map[string]bool.
// Mapping string ip:port to health status.
fmt.Println(myPingu.PingTable())

table := myPingu.PingTable()
// true: target's Pingu is healthy
// false: target's Pingu is unhealthy
fmt.Println(table["127.0.0.1:8552"])
```

### Controll the Pingu
```go
myPingu.Stop()
```
```go
// Continue previous BroadcastPingWithTicker if exist.
myPingu.Start()

// If you want stop the BroadcastPingWithTicker, close the 'cancel'.
cancel, _ := myPingu.BroadcastPingWithTicker(*ticker, 3*time.Second)
close(cancel)
```



