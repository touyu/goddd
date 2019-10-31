package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"goddd/godddcore"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "goddd",
	Short: "Goddd is code generator for layered architecture",
	Long: "",
	Args: cobra.MinimumNArgs(1),
	Run: run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	name := args[0]
	godddcore.Run(name)
}
