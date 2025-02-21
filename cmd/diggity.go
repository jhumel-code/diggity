package cmd

import (
	"os"

	"github.com/carbonetes/diggity/internal/cli"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/spf13/cobra"
)

var (
	diggity = &cobra.Command{
		Use:   "diggity",
		Args:  cobra.MaximumNArgs(1),
		Short: "BOM diggity SBOM Analyzer",
		Long:  `BOM Diggity is an open-source tool developed to streamline the critical process of generating comprehensive Software Bill of Materials (SBOMs) for Container Images and File Systems across various supported ecosystems.`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 && !flagHasArg() {
				_ = cmd.Help()
				os.Exit(0)
			}
			err := ValidateOutputArg(string(*Arguments.Output))
			if err != nil {
				log.Error(err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) > 0 {
				if flagHasArg() {
					log.Println(`"diggity [-d, --dir]" or diggity "[-t, --tar]" does not support with argument image`)
					os.Exit(127)
				}
				Arguments.Image = &args[0]
			} else if image != nil {
				Arguments.Image = image
			} else if flagHasArg() {
				//continue to sbom.Start
			} else {
				if len(args) == 0 || len(*Arguments.Image) == 0 {
					log.Printf(`"diggity [-i, --image]" is required or at least 1 argument "diggity [image]"`)
				}
				os.Exit(127)
			}
			if !*Arguments.Quiet {
				ui.Enable()
			}
			cli.Start(Arguments)
		},
	}
)
