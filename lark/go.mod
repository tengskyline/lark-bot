module lark

go 1.24

replace github.com/tengskyline/lark-bot/conf => ../conf

require (
	github.com/larksuite/oapi-sdk-go/v3 v3.4.19
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/tengskyline/lark-bot/conf v0.0.0
	github.com/tengskyline/lark-bot/qwencli v0.0.0
)

require (
	github.com/devinyf/dashscopego v0.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	go.uber.org/mock v0.5.2 // indirect
	golang.org/x/net v0.41.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tengskyline/lark-bot/qwencli => ../qwencli
