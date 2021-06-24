package regia

import "fmt"

var title = formatColor(96, "[REGIA URL INFO]")

var Banner = `
██████╗ ███████╗ ██████╗ ██╗ █████╗ 
██╔══██╗██╔════╝██╔════╝ ██║██╔══██╗
██████╔╝█████╗  ██║  ███╗██║███████║
██╔══██╗██╔══╝  ██║   ██║██║██╔══██║
██║  ██║███████╗╚██████╔╝██║██║  ██║
╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝╚═╝  ╚═╝
`

// Starter will be called while engine is running
type Starter interface {
	Start(engine *Engine)
}

type BannerStarter struct{ Banner string }

func (b BannerStarter) Start(engine *Engine) { fmt.Println(b.Banner) }

type UrlInfoStarter struct{}

func (u UrlInfoStarter) Start(engine *Engine) {
	for method, nodes := range engine.GetMethodTree() {
		m := formatColor(97, method)
		for _, n := range nodes {
			handleCount := formatColor(colorBlue, fmt.Sprintf("%d handlers", len(n.group)))
			path := formatColor(colorYellow, n.path)
			fmt.Printf("%-15s   %-18s   %-18s   %s\n", title, m, handleCount, path)
		}
	}
}
