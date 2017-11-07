# usersvc

## Building the proto files

Install protobuf (via brew)
Install the micro fork of protoc-gen-go, the protobuf compiler for Go.
```
go get github.com/micro/protobuf/{proto,protoc-gen-go}
```

After modifying the proto definition, recompile it
```
protoc -I$GOPATH/src --go_out=plugins=micro:$GOPATH/src \
	$GOPATH/src/github.com/baddayduck/services/usersvc/proto/account/account.proto
```