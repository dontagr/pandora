package cli

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
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

func NewStoreListCmd(client *transport.HTTPManager, service *store.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "list",
		Short: "Action for list stored entities",
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			reqStoreList, err := service.GetRequestStoreList(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqStoreList)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreList, true)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
				storeList, err := service.RecoveryStoreList(respoce.Body)
				if err != nil {
					redPrint(cmd, err.Error())
					return
				}

				bluePrint(cmd, "Get listing was successfully, pls look at him:")
				color.Set(color.FgHiBlue)
				for i, secretLite := range storeList.List {
					cmd.Printf("\n%d: {label: \"%s\", dt: \"%s\", ver: \"%d\"}", i, secretLite.Label, secretLite.DT, secretLite.Version)
				}
				color.Unset()
			} else {
				redPrint(cmd, fmt.Sprintf("Failed load listing of secret by type=%s: %v", reqStoreList.Kind, string(respoce.Body)))
				return
			}
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of saving data, one of [auth|text|binary|card]")

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

			err = service.Encrypt(reqStoreSave)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqStoreSave)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			synced := false
			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreSave, true)
			if err != nil {
				redPrint(cmd, fmt.Sprintf("Failed sync with remote server: %v", err.Error()))
			} else {
				if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
					bluePrint(cmd, "Save and Sync secret was successfully")
					synced = true
				} else {
					redPrint(cmd, fmt.Sprintf("Sync secret failed: %v", string(respoce.Body)))
				}
			}

			err = service.SaveLocalChanges(reqStoreSave, synced, client.User.Login)
			if err != nil {
				redPrint(cmd, fmt.Sprintf("Failed to save local: %v", string(respoce.Body)))
			} else {
				bluePrint(cmd, "Local BD was updated successfully")
			}
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of saving data, one of [auth|text|binary|card]")
	cmd.Flags().StringP("data", "d", "", "JSON format according to type")
	cmd.Flags().StringP("file", "f", "", "File path")
	cmd.Flags().StringP("label", "l", "", "Label")
	cmd.Flags().StringP("meta", "m", "", "Meta something elso in JSON format")

	return GenericCommand{cmd: &cmd, fullName: "root store save"}
}

func NewStoreLoadCmd(client *transport.HTTPManager, service *store.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "load",
		Short: "Action for load stored entity",
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			reqStoreLoad, err := service.GetRequestStoreLoad(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqStoreLoad)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreLoad, true)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
				storeLoad, err := service.RecoveryStoreLoad(respoce.Body)
				if err != nil {
					redPrint(cmd, err.Error())
					return
				}

				bluePrint(cmd, "Load from server secret was successfully")
				bluePrint(cmd, fmt.Sprintf("You secret: %v", storeLoad.ReveryData))
				bluePrint(cmd, fmt.Sprintf("You meta: %v", storeLoad.Meta))
			} else if respoce.Status == http.StatusNotFound {
				redPrint(cmd, fmt.Sprintf("Secret with type=%s and label=%s not found", reqStoreLoad.Kind, reqStoreLoad.Label))
				return
			} else {
				redPrint(cmd, fmt.Sprintf("Sync secret failed: %v", string(respoce.Body)))
				return
			}
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of saving data, one of [auth|text|binary|card]")
	cmd.Flags().StringP("label", "l", "", "Label")

	return GenericCommand{cmd: &cmd, fullName: "root store load"}
}

func NewStoreDeleteCmd(client *transport.HTTPManager, service *store.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "delete",
		Short: "Action for delete stored entity",
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			reqStoreLoad, err := service.GetRequestStoreLoad(cmd)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			req, err := client.PreparationReq(reqStoreLoad)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreDelete, true)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
				bluePrint(cmd, "Secret was successfully delete")
			} else if respoce.Status == http.StatusNotFound {
				redPrint(cmd, fmt.Sprintf("Secret with type=%s and label=%s not found", reqStoreLoad.Kind, reqStoreLoad.Label))
				return
			} else {
				redPrint(cmd, fmt.Sprintf("Failed delet secret: %v", string(respoce.Body)))
			}
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of saving data, one of [auth|text|binary|card]")
	cmd.Flags().StringP("label", "l", "", "Label")

	return GenericCommand{cmd: &cmd, fullName: "root store delete"}
}

func NewStoreSyncCmd(client *transport.HTTPManager, service *store.Service) GenericCommand {
	cmd := cobra.Command{
		Use:   "sync",
		Short: "Action for sync local changes with remote server",
		Run: func(cmd *cobra.Command, args []string) {
			if client.User.Token == "" {
				redPrint(cmd, "Authorization required, please log in")
				return
			}

			reqSyncList, err := service.GetRequestStoreSync(client.User.Login)
			if err != nil {
				redPrint(cmd, err.Error())
				return
			}

			for _, node := range reqSyncList.List {
				req, err := client.PreparationReq(node)
				if err != nil {
					redPrint(cmd, err.Error())
					return
				}

				respoce, err := client.NewRequest(http.MethodPost, req, models.UrlStoreSync, true)
				if err != nil {
					redPrint(cmd, err.Error())
					return
				}

				if respoce.Status >= http.StatusOK && respoce.Status < http.StatusMultipleChoices {
					bluePrint(cmd, fmt.Sprintf("Secret type=%s and lable=%s was successfully sync", node.Kind, node.Label))

					syncResp, err := service.ResponceSync(respoce.Body)
					if err != nil {
						redPrint(cmd, err.Error())
						return
					}

					err = service.UpdateLocalChanges(client.User.Login, node.Kind, node.Label, syncResp.Version)
					if err != nil {
						redPrint(cmd, err.Error())
						return
					}
				} else {
					redPrint(cmd, fmt.Sprintf("Failed sync secret: %v", string(respoce.Body)))
				}
			}
		},
	}

	return GenericCommand{cmd: &cmd, fullName: "root store sync"}
}
