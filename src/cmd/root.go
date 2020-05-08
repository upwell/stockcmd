package cmd

import (
	"fmt"
	"os"

	"hehan.net/my/stockcmd/logger"

	"hehan.net/my/stockcmd/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "stockcmd",
	Short:        "main entry point for stockcmd",
	SilenceUsage: true,
	//SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.json)")

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		logger.InitLogger()
		return nil
	}
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "verbose output")
}
