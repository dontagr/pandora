package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/agent/config"
	"github.com/dontagr/pandora/internal/agent/service/transport"
)

type Commander struct {
	zlog *zap.SugaredLogger
}

func NewCommander(
	commands []CobraCommand,
	zlog *zap.SugaredLogger,
	lc fx.Lifecycle,
	shutdowner fx.Shutdowner,
	client *transport.HTTPManager,
	config *config.Config,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pandora [command]",
		Version: config.Version,
		Short:   "Pandora - client app to save your secrets.",
		Long: `The Pandora app will help you store your data safely without fear of it being stolen.

It's quite easy to start using it, just sign up and store you data`,
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				zlog.Errorf("handler help failed: %w", err)
			}

			if client.User.Login != "" {
				bluePrint(cmd, fmt.Sprintf("Вы авторизованы как %s", client.User.Login))
			} else {
				redPrint(cmd, "Остановись. Сначало нужно авторизоваться")
			}
		},
	}

	cmd.SetVersionTemplate(`{{with .DisplayName}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}} compilation in ` + config.BildDT + "\n")

	done := make(map[string]bool)
	for _, c1 := range commands {
		for _, c2 := range commands {
			if c1.FullName() == c2.FullName() || done[c1.FullName()] {
				continue
			}
			fullParent := strings.Join(strings.Split(c1.FullName(), " ")[:len(strings.Split(c1.FullName(), " "))-1], " ")
			if fullParent == c2.FullName() {
				c2.Command().AddCommand(c1.Command())
				done[c1.FullName()] = true
				continue
			}
			if fullParent == "root" {
				cmd.AddCommand(c1.Command())
				done[c1.FullName()] = true
				continue
			}
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := cmd.Execute(); err != nil {
				//				zlog.Errorf("Whoops. There was an error while executing the CLI '%s'", err)

				_ = shutdowner.Shutdown(fx.ExitCode(1))
			}
			_ = shutdowner.Shutdown()
			return nil
		},
	})

	return cmd
}

func redPrint(cmd *cobra.Command, msg string) {
	color.Set(color.FgHiRed)
	cmd.Println("")
	cmd.Println(msg)
	color.Unset()
}

func bluePrint(cmd *cobra.Command, msg string) {
	color.Set(color.FgHiBlue)
	cmd.Println("")
	cmd.Println(msg)
	color.Unset()
}

func magentaPrint(cmd *cobra.Command, msg string) {
	color.Set(color.FgHiMagenta)
	cmd.Println("")
	cmd.Println(msg)
	color.Unset()
}

func cyanPrint(cmd *cobra.Command, msg string) {
	color.Set(color.FgHiCyan)
	cmd.Println("")
	cmd.Println(msg)
	color.Unset()
}
