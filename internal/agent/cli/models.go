// Package cli содержит реализацию моделей для командной строки.
package cli

import "github.com/spf13/cobra"

// CobraCommand представляет собой интерфейс для командной строки.
// Этот интерфейс обеспечивает доступ к команде и ее полному имени.
type CobraCommand interface {
	Command() *cobra.Command
	FullName() string
}

// GenericCommand представляет собой реализацию командной строки.
// Этот тип обеспечивает доступ к команде и ее полному имени.
type GenericCommand struct {
	cmd      *cobra.Command
	fullName string
}

// Command возвращает команду.
func (g GenericCommand) Command() *cobra.Command {
	return g.cmd
}

// FullName возвращает полное имя команды.
func (g GenericCommand) FullName() string {
	return g.fullName
}
