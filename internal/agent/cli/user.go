package cli

import (
	"net/http"

	"github.com/spf13/cobra"

	"github.com/dontagr/pandora/internal/agent/service/transport"
	userServ "github.com/dontagr/pandora/internal/agent/service/user"
	"github.com/dontagr/pandora/internal/agent/store/user"
	"github.com/dontagr/pandora/internal/models"
)

func NewUserCmd() GenericCommand {
	cmd := cobra.Command{
		Use:   "user",
		Short: "Action for users function",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}
		},
	}
	return GenericCommand{cmd: &cmd, fullName: "root user"}
}

func NewUserSignUPCmd(client *transport.HTTPManager, userStore *user.User, service *userServ.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "sign-up [-l login][-p password]",
		Short: "Action for new user registration",
		Long: `
The password must meet the following conditions:
- at least 6 characters
- contains numbers
- at least one letter in the entry and at the top of the registration
`,
		Run: func(cmd *cobra.Command, args []string) {
			reqUser, err := service.GetLoginPassword(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqUser)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlUserSignUP, false)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
				bluePrint(cmd, "Registration was successfully")
			} else {
				if respoce.Status == http.StatusConflict {
					redPrint(cmd, "Registration was failed because 'Логин уже занят'")
				} else {
					redPrint(cmd, "Registration was failed")
				}
				return
			}

			err = userStore.SaveUserAuth(reqUser.Login, respoce.Meta.Authorization)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}
		},
	}
	cmd.Flags().StringP("login", "l", "", "User login")
	cmd.Flags().StringP("password", "p", "", "User password")

	return GenericCommand{cmd: &cmd, fullName: "root user sign-up"}
}

func NewUserLoginCmd(client *transport.HTTPManager, userStore *user.User, service *userServ.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "login [-l login][-p password]",
		Short: "Action for user login",
		Run: func(cmd *cobra.Command, args []string) {
			reqUser, err := service.GetLoginPassword(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqUser)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlUserLogin, false)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
				bluePrint(cmd, "Login was successfully")
			} else if respoce.Status == http.StatusUnauthorized {
				redPrint(cmd, "Неверная пара логин/пароль")
				return
			} else {
				redPrint(cmd, "Login was failed")
				return
			}

			err = userStore.SaveUserAuth(reqUser.Login, respoce.Meta.Authorization)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}
		},
	}
	cmd.Flags().StringP("login", "l", "", "User login")
	cmd.Flags().StringP("password", "p", "", "User password")

	return GenericCommand{cmd: &cmd, fullName: "root user login"}
}
