package help

import (
	"embed"
	"fmt"
)

//go:embed docs/*.md
var f embed.FS

// Load loads the documentation asset
func Load(name string) string {

	doc, err := f.ReadFile(fmt.Sprintf("docs/%s.md", name))
	if err != nil {
		panic(err)
	}

	return "\n" + string(doc)
}
