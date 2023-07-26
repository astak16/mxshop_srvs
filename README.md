## user_srvs proto 生成

```bash
protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. user.proto
```