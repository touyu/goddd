package godddcore

type templateData struct {
	Name string
	CurrentDir string
}

func Run(name string) error {
	return generateOutput(name)
}
