package cli

import "github.com/spf13/cobra"

type CobraCommand interface {
	Command() *cobra.Command
	FullName() string
}

type GenericCommand struct {
	cmd      *cobra.Command
	fullName string
}

func (g GenericCommand) Command() *cobra.Command {
	return g.cmd
}

func (g GenericCommand) FullName() string {
	return g.fullName
}
