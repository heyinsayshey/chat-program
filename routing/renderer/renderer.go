package renderer

import (
	"sync"

	"github.com/unrolled/render"
)

var renderer *render.Render
var once sync.Once

func InitRenderer() {
	GetInstance()
}

func GetInstance() *render.Render {
	once.Do(func() {
		renderer = render.New(render.Options{
			Charset: "UTF-8",
		})
	})

	return renderer
}

func Fin() {
}
