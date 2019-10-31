package godddcore

type templateData struct {
	Name string
}

func Run(name string) error {
	return generateOutput(name)
}
