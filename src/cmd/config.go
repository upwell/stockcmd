package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/store"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Config command to set/get configurations",
}

var ConfigSetCmd = &cobra.Command{
	Use:     "set [key] [value]",
	Short:   "set value to config [key] to [value]",
	Long:    "set value to config [key] to [value]",
	Example: `set redisAddr 127.0.0.1`,
	Args:    cobra.MinimumNArgs(2),
	RunE:    setCmdF,
}

var ConfigGetCmd = &cobra.Command{
	Use:     "get [key]",
	Short:   "get value of config [key]",
	Long:    "get value of config [key]",
	Example: `get redisAddr`,
	Args:    cobra.MinimumNArgs(1),
	RunE:    getCmdF,
}

func setCmdF(cmd *cobra.Command, args []string) error {
	key := args[0]
	val := args[1]

	store.RunningConfig.Set(key, val)
	println("Set configuration successfully")
	return nil
}

func getCmdF(cmd *cobra.Command, args []string) error {
	key := args[0]
	val, existed, _ := store.RunningConfig.GetString(key)
	if !existed {
		fmt.Printf("config [%s] not exist\n", key)
	} else {
		fmt.Printf("The value of config [%s] is [%s]\n", key, val)
	}
	return nil
}

func init() {
	ConfigCmd.AddCommand(ConfigGetCmd, ConfigSetCmd)
	rootCmd.AddCommand(ConfigCmd)
}
