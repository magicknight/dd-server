package classification

import (
	"github.com/shmy/dd-server/model"
)

var M = Model{&model.Model{model.Db.C("classifications")}}

type Model struct {
	*model.Model
}
