# PINGU [![GoDoc](https://godoc.org/github.com/dbadoy/pingu?status.svg)](https://godoc.org/github.com/dbadoy/pingu)

This is a simple and tiny heartbeat library with UDP. <br>
If you need build a server with heartbeat, It could make it your building easier. <br>
<b>(Only supports communication between internal networks)</b>

## Install
```
$ go get -u github.com/dbadoy/pingu
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



