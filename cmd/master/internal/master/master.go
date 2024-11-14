package master

import (
	"os"

	model "github.com/mahaonan001/dsfs/cmd/master/internal/model"
	"gopkg.in/yaml.v3"
)

var Clients model.Clients

func init() {
	data, err := os.ReadFile("clients.txt")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &Clients)
	if err != nil {
		panic(err)
	}
}
