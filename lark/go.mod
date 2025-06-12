module lark

go 1.24

replace github.com/tengskyline/lark-bot/conf => ../conf

require (
	github.com/larksuite/oapi-sdk-go/v3 v3.4.18
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/tengskyline/lark-bot/conf v0.0.0-00010101000000-000000000000
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
