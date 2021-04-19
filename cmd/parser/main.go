package main

import (
	"context"
	"log"

	"github.com/and07/parser/pkg/parser"
	"github.com/spf13/cobra"
)

func main() {

	var rulePath, path string

	var rootCmd = &cobra.Command{
		Use:   "parser",
		Short: "parser",
		Long:  `parser`,
		Run: func(cmd *cobra.Command, args []string) {

			rule := cmd.Flag("rule")
			rulePath = rule.Value.String()

			csv := cmd.Flag("csv")
			path = csv.Value.String()
			log.Printf("PATH %s", path)
		},
	}

	rootCmd.Flags().StringP("rule", "r", "", "file rule for parsing (required)")
	//rootCmd.MarkFlagRequired("rule")
	rootCmd.Flags().StringP("csv", "c", "", "Output to csv file")
	rootCmd.Flags().BoolP("debug", "d", false, "debug turn off/on")

	if err := rootCmd.Execute(); err != nil {
		log.Printf("ERROR %s", err)
	}

	ctx, cancel := parser.New(context.Background(), parser.WithWriter(path), parser.WithConfigs(rulePath))
	defer cancel()

	if _, err := parser.Run(ctx); err != nil {
		log.Printf("ERROR parser.Run %s", err)
	}
}
