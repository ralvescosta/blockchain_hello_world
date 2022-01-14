package cmd

type CliCmd struct{}

func (CliCmd) Run() {
	container := NewCliContainer()

	container.cli.Run()
}
