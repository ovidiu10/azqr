// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azqr

import (
	"github.com/cmendible/azqr/internal/scanners"
	"github.com/cmendible/azqr/internal/scanners/oai"
	"github.com/spf13/cobra"
)

func init() {
	scanCmd.AddCommand(oaiCmd)
}

var oaiCmd = &cobra.Command{
	Use:   "oai",
	Short: "Scan Azure OpenAI Services",
	Long:  "Scan Azure OpenAI Services",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		serviceScanners := []scanners.IAzureScanner{
			&oai.OpenAIScanner{},
		}

		scan(cmd, serviceScanners)
	},
}
