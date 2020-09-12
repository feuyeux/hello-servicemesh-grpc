## grpc python demo
### Setup
```sh
sudo python3 -m pip install grpcio --ignore-installed
sudo python3 -m pip install grpcio-tools
```

#### generate
```sh
sh proto2py.sh
```

#### test
```sh
sh test_server.sh
```

```sh
sh test_client.sh
```