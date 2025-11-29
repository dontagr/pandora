package bootstrap

import (
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"github.com/dontagr/pandora/internal/agent/cli"
)

// CLI является набором параметров и функций для инициализации и настройки командной строки приложения.
// Он использует библиотеку fx для предоставления зависимостей и упрощения процесса настройки.
var CLI = fx.Options(
	fx.Provide(
		AsCommand(cli.NewUserCmd),
		AsCommand(cli.NewUserSignUPCmd),
		AsCommand(cli.NewUserLoginCmd),
		AsCommand(cli.NewStoreCmd),
		AsCommand(cli.NewStoreListCmd),
		AsCommand(cli.NewStoreLoadCmd),
		AsCommand(cli.NewStoreSaveCmd),
		AsCommand(cli.NewStoreDeleteCmd),
		AsCommand(cli.NewStoreSyncCmd),
		fx.Annotate(
			cli.NewCommander,
			fx.ParamTags(`group:"commands"`),
		),
	),
	fx.NopLogger, // Note, this currently removes all fx logs, even on error
	fx.Invoke(func(*cobra.Command) {}),
)

// AsCommand аннотирует функцию, переданную в параметре f, и преобразует ее в команду Cobra.
// Возвращаемое значение добавляется в группу "commands".
func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(cli.CobraCommand)),
		fx.ResultTags(`group:"commands"`),
	)
}
