package features

type app struct {
	container Container
}

func newApp(container Container) app {
	return app{container: container}
}
