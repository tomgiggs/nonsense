
# 简介
纯golang编写的IM服务端，使用protobuf协议进行通讯，支持tcp直连，rpc调用，websocket调用，可动态扩缩容，项目目标是实现超大规模用户同时在线

# 安装与部署
## 编译
安装go开发工具
cd cmd/access
go build

## 安装环境
为了方便部署，使用docker部署用到的组件
### 部署consul:
docker run --name consul1 -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600 consul:latest agent -server -bootstrap-expect 1 -ui -bind=0.0.0.0 -client=0.0.0.0
### 部署redis:
docker run --name redisDemo -p 36379:6379 -d redis
### 部署mysql
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123 -d mysql:5.7

### 创建mysql库表，并刷入数据
 先刷sql/create.sql，再刷update.sql
### 启动
复制一份internal/config/access-example.toml 配置文件，并做相应修改
然后执行：
./access -conf ../../internal/config/access-example.toml启动服务即可，后面为配置文件路径


# 开发
## 生成proto对应代码
修改proto文件后需要重新生成对应的代码，使用scripts/gen_pb.bat即可生成对应代码，需要提前安装好protoc工具
```
grpc-生成go代码命令：
cd pkg
protoc --proto_path=./proto --plugin=protoc-gen-go.exe --go_out=proto/ --go_opt=paths=source_relative proto/api.proto
protoc --proto_path=./proto --plugin=protoc-gen-go.exe --go_out=proto/ --go_opt=paths=source_relative proto/message.proto
cd proto
protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  service.proto

proto生成js代码：
protoc.exe --js_out=import_style=commonjs,binary:. message.proto
pbjs -t json-module -w commonjs -o scripts/proto.js pkg/proto/*.proto
```
## goland远程调试
1. 安装dlv工具

go get -u github.com/go-delve/delve/cmd/dlv
2. 然后使用dlv启动服务

dlv --listen=:9004 --headless=true --api-version=2 exec ./access -- start --env debug
3. goland 中添加go remote调试

run/debug configurations--->add--->go remote--->host--->port

## 功能测试
基本测试代码在test目录下，入口文件在entrance.go

## 数据存储
发送视频/音频/图片/文件等先存储到云服务上，然后发送链接即可，数据库当前为mysql，后续将调整为Cassandra/MongoDB

## TODO
- 用户消息加密/加密算法调整
- 服务发现接口
- 抽象数据存储接口
- 优化扩展性，方便添加新功能
- 使用orm替代直接写SQL
- 通用错误码定义与使用/报错信息透传
- 优化代码，提高健壮性
- 性能测试

# 参考
项目由 https://github.com/alberliu/gim 魔改而来
