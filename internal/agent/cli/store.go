package cli

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/dontagr/pandora/internal/agent/service/store"
	"github.com/dontagr/pandora/internal/agent/service/transport"
	"github.com/dontagr/pandora/internal/models"
)

func NewStoreCmd() GenericCommand {
	cmd := cobra.Command{
		Use:   "store [command]",
		Short: "Action for working with storage",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return GenericCommand{cmd: &cmd, fullName: "root store"}
}

func NewStoreListCmd(client *transport.HTTPManager) GenericCommand {
	cmd := cobra.Command{
		Use:   "list",
		Short: "Action for list stored entities",
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			respoce, err := client.NewRequest(http.MethodGet, nil, models.UrlTest, true)
			if err != nil {
				redPrint(cmd, err.Error())
			}

			if respoce.Status == http.StatusUnauthorized {
				redPrint(cmd, "Authorization required, please log in")
				return
			}
		},
	}
	return GenericCommand{cmd: &cmd, fullName: "root store list"}
}

func NewStoreSaveCmd(client *transport.HTTPManager, service *store.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "save",
		Short: "Action for save stored entity",
		Long: `
Example json data for following types:
- auth type - requires input data in JSON format
======
{"login":"test","password":"123456"}
======
- text type - requires input data in txt format
======
"some string here"
======
- card type - requires input data in JSON format
======
{"number":"1234123412341234","date":"11/23","cvv":"123"}
======

For binary type you need to specify the path to the file using the flag -f
Files up to 1 MB are allowed.

`,
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			reqStoreSave, err := service.GetRequestStoreSave(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			reqStoreSave, err = service.Encrypt(reqStoreSave)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqStoreSave)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreSave, true)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			fmt.Println(respoce)

			cmd.Println("Try to save from store")
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of saving data, one of [auth|text|binary|card]")
	cmd.Flags().StringP("data", "d", "", "JSON format according to type")
	cmd.Flags().StringP("file", "f", "", "File path")
	cmd.Flags().StringP("label", "l", "", "Label")
	cmd.Flags().StringP("meta", "m", "", "Meta")

	return GenericCommand{cmd: &cmd, fullName: "root store save"}
}

func NewStoreLoadCmd() GenericCommand {
	cmd := cobra.Command{
		Use:   "load",
		Short: "Action for load stored entity",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Try to load from store")
		},
	}
	return GenericCommand{cmd: &cmd, fullName: "root store load"}
}

func NewStoreDeleteCmd() GenericCommand {
	cmd := cobra.Command{
		Use:   "delete",
		Short: "Action for delete stored entity",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Try to delete stored entity")
		},
	}
	return GenericCommand{cmd: &cmd, fullName: "root store delete"}
}
