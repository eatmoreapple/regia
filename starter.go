package regia

import "fmt"

var title = formatColor("[REGIA URL INFO]", 96)

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
		m := formatColor(method, 97)
		for _, n := range nodes {
			handleCount := formatColor(fmt.Sprintf("%d handlers", len(n.group)), colorBlue)
			path := formatColor(n.path, colorYellow)
			fmt.Printf("%-15s   %-18s   %-18s   %s\n", title, m, handleCount, path)
		}
	}
}
