cd ../pkg
protoc --proto_path=./proto --plugin=protoc-gen-go.exe --go_out=proto/ --go_opt=paths=source_relative proto/api.proto
protoc --proto_path=./proto --plugin=protoc-gen-go.exe --go_out=proto/ --go_opt=paths=source_relative proto/message.proto
cd proto
protoc --go_out=. --go_opt=paths=source_relative  --go-grpc_out=. --go-grpc_opt=paths=source_relative  service.proto

cd ../../
pbjs -t json-module -w commonjs -o scripts/proto.js pkg/proto/*.proto

exit