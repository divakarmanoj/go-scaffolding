package cmd

import (
	"fmt"
	"github.com/divakarmanoj/go-scaffolding/generator"
	"github.com/spf13/cobra"
	"os"
)

// interactiveCmd represents the interactive command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "This command will start an interactive command line interface to generate the model",
	Long:  `This command will start an interactive command line interface to generate the model.`,
	Run: func(cmd *cobra.Command, args []string) {
		outputDir, _ := cmd.Flags().GetString("output-dir")
		config, err := generator.InteractiveConfigGeneration("", false)
		if err != nil {
			fmt.Printf("Prompt failed: %v", err)
			os.Exit(1)
		}
		config.Generate(outputDir)
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
