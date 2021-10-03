// Copyright 2021 eatMoreApple.  All rights reserved.
// Use of this source code is governed by a GPL style
// license that can be found in the LICENSE file.

package regia

import (
	"fmt"
	"github.com/eatmoreapple/regia/internal"
)

var Banner = `
██████╗ ███████╗ ██████╗ ██╗ █████╗ 
██╔══██╗██╔════╝██╔════╝ ██║██╔══██╗
██████╔╝█████╗  ██║  ███╗██║███████║
██╔══██╗██╔══╝  ██║   ██║██║██╔══██║
██║  ██║███████╗╚██████╔╝██║██║  ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
`

const urlInfo = "[URL INFO]"

// Starter will be called while engine is running
type Starter interface {
	Start(engine *Engine)
}

type BannerStarter struct{ Banner string }

func (b BannerStarter) Start(engine *Engine) { fmt.Println(b.Banner) }

type UrlInfoStarter struct{}

func (u UrlInfoStarter) Start(engine *Engine) {
	for method, nodes := range engine.methodsTree {
		m := internal.FormatColor(97, method)
		for _, n := range nodes {
			handleCount := internal.BlueString(fmt.Sprintf("%d handlers", len(n.group)))
			path := internal.YellowString(n.path)
			fmt.Printf("%-15s   %-18s   %-18s   %s\n", urlInfo, m, handleCount, path)
		}
	}
}
