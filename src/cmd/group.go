package cmd

import (
	"fmt"

	"hehan.net/my/stockcmd/sina"

	"github.com/pkg/errors"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/store"
)

var GroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Management of groups",
}

var GroupListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list all groups",
	Long:    `List all groups.`,
	Example: ` group list`,
	RunE:    listGroupCmdF,
}

var GroupCreateCmd = &cobra.Command{
	Use:     "create [group]",
	Short:   "Create a group",
	Long:    `Create a group.`,
	Example: ` group create mygroup`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    createGroupCmdF,
}

var GroupDeleteCmd = &cobra.Command{
	Use:     "delete [group]",
	Short:   "Delete a group",
	Long:    `Delete a group.`,
	Example: ` group delete mygroup`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    deleteGroupCmdF,
}

var AddGroupStockCmd = &cobra.Command{
	Use:   "add [group] [stock]",
	Short: "Add stock to group",
	Long:  `Add stock to group, the [stock] can be a hint and then command would prompt suggestions.`,
	Args:  cobra.MinimumNArgs(2),
	RunE:  addGroupStockCmdF,
}

var RemoveGroupStockCmd = &cobra.Command{
	Use:   "remove [group]",
	Short: "Remove stock from group",
	Long:  `Remove stock from group, the stocks list would prompt suggestions for selection`,
	Args:  cobra.MinimumNArgs(1),
	RunE:  removeGroupStockCmdF,
}

func parsePromptError(err error) error {
	switch err {
	case promptui.ErrInterrupt:
		return nil
	case promptui.ErrAbort:
		return nil
	default:
		return errors.Wrap(err, "prompt failed")
	}
}

func listGroupCmdF(cmd *cobra.Command, args []string) error {
	groups := store.ListGroup()
	for _, g := range groups {
		fmt.Println(g)
	}
	return nil
}

func createGroupCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	g := store.GetGroup(gName)
	if g != nil {
		return errors.Errorf("Group [%s] already exist", gName)
	}

	store.AddGroup(gName)
	return nil
}

func deleteGroupCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	g := store.GetGroup(gName)
	if g == nil {
		return errors.Errorf("Group [%s] not exist", gName)
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Are you sure to delete group [%s]", gName),
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return parsePromptError(err)
	}
	if result == "y" {
		store.DeleteGroup(gName)
		fmt.Printf("Group [%s] deleted\n", gName)
	}
	return nil
}

func addGroupStockCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	stockHint := args[1]
	g := store.GetGroup(gName)
	if g == nil {
		return errors.Errorf("Group [%s] not exist", gName)
	}

	suggestions := sina.Suggest(stockHint)
	if len(suggestions) == 0 {
		return errors.Errorf("[%s] does not match any stock", stockHint)
	}

	prompt := promptui.Select{
		Label: "Select the matched stock to add",
		Items: suggestions,
		Size:  10,
		Templates: &promptui.SelectTemplates{
			Active:   fmt.Sprintf(`%s {{ index . "name" | underline }}`, promptui.IconSelect),
			Inactive: `  {{ index . "name" }}`,
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ index . "name" | faint }}`, promptui.IconGood),
		},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return parsePromptError(err)
	}

	retMap := suggestions[idx]
	g.AddStock(retMap["code"], retMap["name"])
	fmt.Printf("Stock [%s] added to group [%s]", retMap["name"], gName)
	return nil
}

func removeGroupStockCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	g := store.GetGroup(gName)
	if g == nil {
		return errors.Errorf("Group [%s] not exist", gName)
	}

	stocks := make([]map[string]string, 0, 32)
	for code, name := range g.Codes {
		stocks = append(stocks, map[string]string{
			"code": code,
			"name": name,
		})
	}

	prompt := promptui.Select{
		Label: "Select the stock to remove",
		Items: stocks,
		Size:  10,
		Templates: &promptui.SelectTemplates{
			Active:   fmt.Sprintf(`%s {{ index . "name" | underline }}`, promptui.IconSelect),
			Inactive: `  {{ index . "name" }}`,
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ index . "name" | faint }}`, promptui.IconGood),
		},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return parsePromptError(err)
	}

	retMap := stocks[idx]
	g.RemoveStock(retMap["code"])
	fmt.Printf("Stock [%s] removed from group [%s]", retMap["name"], gName)
	return nil
}

func init() {
	GroupCmd.AddCommand(
		GroupListCmd,
		GroupCreateCmd,
		GroupDeleteCmd,
		AddGroupStockCmd,
		RemoveGroupStockCmd,
	)
	rootCmd.AddCommand(GroupCmd)
}
