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
	engine := regia.Default()
	engine.GET("/", func(context *regia.Context) {
		context.JSON(regia.Map{"hello": "world"})
	})
	engine.POST("/", func(context *regia.Context) {
		var form struct {
			Name    string                `form:"name"`
			Hobbies []string              `form:"hobbies"`
			Avatar  *multipart.FileHeader `file:"avatar"`
		}
		// Parse the form.
		if err := context.BindMultipartForm(&form); err != nil {
			context.JSON(regia.Map{"err": err.Error()})
			return
		}
		// save upload file
		path, err := context.FileStorage.Save(form.Avatar)
		if err != nil {
			context.JSON(regia.Map{"err": err.Error()})
			return
		}
		context.JSON(regia.Map{"name": form.Name, "hobbyies": form.Hobbies, "avatar": path})
	})
	engine.Run(":8000")
}
```

```shell
$ go run main.go
# open your brower and visit `localhost:8000/`
```



#### Bind Form Data

```go
```



