# regia

Regia is a web framework written with golang ! 

It is simple, helpful and with high performance. Build your own idea with it !

## Installation

Golang version 1.11 + required

```shell
go get github.com/eatmoreapple/regia
```



## Quick Start

```sh
$ touch main.go
# add all following codes into main.go
```

```go
package main

import "github.com/eatmoreapple/regia"

func main() {
	engine := regia.New()
	engine.GET("/string", func(context *regia.Context) {
		context.String("string")
	})
	engine.POST("/json", func(context *regia.Context) {
		var body struct {
			Name string `json:"name"`
			Age  uint8  `json:"age"`
		}
		context.BindJSON(&body)
		context.JSON(body)
	})
	engine.Run(":8000")
}
```

```shell
$ go run main.go
# open your brower and visit `localhost:8000/`
```









