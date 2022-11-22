module github.com/koddr/example-go-grpc-server

go 1.13

replace github.com/keploy/go-sdk => ../../go-sdk

replace go.keploy.io/server => ../../myfolder/keploy

require (
	github.com/keploy/go-sdk v0.6.5
	go.mongodb.org/mongo-driver v1.8.3
	golang.org/x/net v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)
