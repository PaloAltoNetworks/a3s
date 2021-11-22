//go:generate go-bindata -pkg help -o bindata.go ../../docs

package help

import "fmt"

// Load loads the documentation asset
func Load(name string) string {

	doc, err := Asset(fmt.Sprintf("../../docs/%s.md", name))
	if err != nil {
		panic(err)
	}

	return "\n" + string(doc)
}
