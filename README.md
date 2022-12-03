## gRPC-Go gateway best practices

gRPC-Gateway 是一个 protoc 插件。它读取 gRPC 服务定义并生成一个反向代理服务器，该服务器将 RESTful JSON API 转换为 gRPC。

**1. 安装插件**

```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
```

**2. 生成 gRPC 桩代码**

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_opt=paths=source_relative go_wallet_manage_svr.proto
```

**3. 添加 gRPC-Gateway 注释**

```text
...
import "google/api/annotations.proto";

...

option (google.api.http) = {
  post: "/v1/create_wallet"
  body: "*"
};
```

**4. 生成 gRPC-Gateway 桩代码**

```shell
protoc -I . --grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative go_wallet_manage_svr.proto
```

**5. 测试 gRPC-Gateway**

使用 cURL 发送 HTTP 请求：

```shell
curl -X POST -k http://127.0.0.1:8090/v1/create_wallet

curl -X POST -k http://127.0.0.1:8090/v1/import_wallet -d '{"private_key": "0x12345"}'
```

得到响应结果：

```json
{"address":"02efe390-2edb-4e26-a72a-ba5eac5f5b30"}
```