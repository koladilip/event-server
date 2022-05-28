package main

import (
	"github.com/koladilip/event-server/api"
	"github.com/koladilip/event-server/config"
	"github.com/koladilip/event-server/destination"
	"github.com/koladilip/event-server/source"
	"github.com/koladilip/event-server/store"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Fx,
		store.Fx,
		source.Fx,
		destination.Fx,
		api.Fx,
	)
	app.Run()
}
