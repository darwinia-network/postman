package postman

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Get toml config from path
func Config(path string) (config *toml.Tree, err error) {
	doc, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config, err = toml.Load(string(doc))
	if err != nil {
		return nil, err
	}

	return
}
