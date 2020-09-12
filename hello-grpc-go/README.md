## grpc golang demo
### Setup
```sh
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```
### Generate
```sh
sh proto2go.sh
```
### Run
```sh
sh test_server.sh
```

```sh
sh test_client.sh
```