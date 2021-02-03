package main

import (
	"context"
	"log"
	"os"

	"github.com/and07/parser/pkg/parser"
	"github.com/spf13/cobra"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {

	ctx, cancel := parser.New(context.Background())
	defer cancel()
	var err error

	var rootCmd = &cobra.Command{
		Use:   "parser",
		Short: "parser",
		Long:  `parser`,
		Run: func(cmd *cobra.Command, args []string) {
			if isInputFromPipe() {
				if ctx, err = parser.Conf(ctx, os.Stdin); err != nil {
					log.Printf("ERROR parser.Conf %s", err)
				}
			} else {
				rule := cmd.Flag("rule")
				rulePath := rule.Value.String()
				if rulePath == "" {
					rulePath = "./rule.json"
				}
				if ctx, err = parser.RuleConfig(ctx, rulePath); err != nil {
					log.Printf("ERROR parser.Rule %s", err)
				}
			}

			if ctx, err = parser.Run(ctx); err != nil {
				log.Printf("ERROR parser.Run %s", err)
			}

			csv := cmd.Flag("csv")
			path := csv.Value.String()
			log.Printf("PATH %s", path)
			if path != "" {
				if err = parser.ExportCSV(ctx, path); err != nil {
					log.Printf("ERROR parser.ExportCSV %s", err)
				}
			} else {
				parser.Output(ctx, os.Stdout)
			}

		},
	}

	rootCmd.Flags().StringP("rule", "r", "", "file rule for parsing (required)")
	//rootCmd.MarkFlagRequired("rule")
	rootCmd.Flags().StringP("csv", "c", "", "Output to csv file")
	rootCmd.Flags().BoolP("debug", "d", false, "debug turn off/on")

	if err := rootCmd.Execute(); err != nil {
		log.Printf("ERROR %s", err)
	}

}
