module github.com/tengskyline/lark-bot

go 1.24

require (
	github.com/tengskyline/lark-bot/conf v0.0.0
    github.com/tengskyline/lark-bot/lark v0.0.0
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/larksuite/oapi-sdk-go/v3 v3.4.18 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tengskyline/lark-bot/conf => ./conf

replace github.com/tengskyline/lark-bot/lark => ./lark
