# Teamspeak 3 ServerQuery

[![GoDoc Widget]][GoDoc] [![Go Report Card Widget]][Go Report Card]

[GoDoc]: https://godoc.org/github.com/paralin/ts3-go
[GoDoc Widget]: https://godoc.org/github.com/paralin/ts3-go?status.svg
[Go Report Card Widget]: https://goreportcard.com/badge/github.com/paralin/ts3-go
[Go Report Card]: https://goreportcard.com/report/github.com/paralin/ts3-go

## Introduction

**go-ts3** is a Go client for the **ServerQuery** API in TeamSpeak 3.

The [ServerQuery API Specification](http://media.teamspeak.com/ts3_literature/TeamSpeak%203%20Server%20Query%20Manual.pdf) has the relevant information about the supported APIs.

## API Structures

API structures can be encoded into "ServerQuery" syntax, which looks like:

```
serverlist
clientlist –uid –away –groups
clientdbfind pattern=FPMPSC6MXqXq751dX7BKV0JniSo= –uid
clientkick reasonid=5 reasonmsg=Go\saway! clid=1|clid=2|clid=3
channelmove cid=16 cpid=1 order=0
sendtextmessage targetmode=2 target=12 msg=Hello\sWorld!endtextmessage 
```

This package can marshal (encode) and unmarshal (decode) the ServerQuery protocol intelligently, supporting anonymous types and nested / ignored fields.

```golang
// TargetMode specifies which kind of target to use.
type TargetMode int

// SendTextMessage sends text messages to channels or users.
type SendTextMessage struct {
	// TargetMode is the target mode of the command.
	TargetMode TargetMode `serverquery:"targetmode"`
}
```

## Getting Started

The following code snippit is approximately how one uses this library. A full example can be seen under the example directory.

```golang
import "github.com/paralin/ts3-go/serverquery"

client, err := serverquery.Dial("localhost:10011")
if err != nil { panic(err) }
client.Use(9877)
client.Login("username", "password")
clientList, err := client.GetClientList()
```

All calls are thread-safe.

