package main

import (
	"go.uber.org/fx"
	"koladilip.github.io/event-server/api"
	"koladilip.github.io/event-server/config"
	"koladilip.github.io/event-server/destination"
	"koladilip.github.io/event-server/source"
	"koladilip.github.io/event-server/store"
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
