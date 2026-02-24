// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package oai

import (
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"github.com/cmendible/azqr/internal/scanners"
)

// GetRules - Returns the rules for the OpenAIScanner
func (a *OpenAIScanner) GetRules() map[string]scanners.AzureRule {
	return map[string]scanners.AzureRule{
		"DiagnosticSettings": {
			Id:          "oai-001",
			Category:    "Monitoring and Logging",
			Subcategory: "Diagnostic Logs",
			Description: "Azure OpenAI should have diagnostic settings enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				service := target.(*armcognitiveservices.Account)
				hasDiagnostics, err := a.diagnosticsSettings.HasDiagnostics(*service.ID)
				if err != nil {
					log.Fatalf("Error checking diagnostic settings for service %s: %s", *service.Name, err)
				}

				return !hasDiagnostics, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/ai-services/diagnostic-logging",
		},
		"SLA": {
			Id:          "oai-002",
			Category:    "High Availability and Resiliency",
			Subcategory: "SLA",
			Description: "Azure OpenAI should have a SLA",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				return false, "99.9%"
			},
			Url: "https://www.azure.cn/en-us/support/sla/cognitive-services/",
		},
		"Private": {
			Id:          "oai-003",
			Category:    "Security",
			Subcategory: "Networking",
			Description: "Azure OpenAI should have private endpoints enabled",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcognitiveservices.Account)
				pe := len(i.Properties.PrivateEndpointConnections) > 0
				return !pe, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/ai-services/cognitive-services-virtual-networks",
		},
		"SKU": {
			Id:          "oai-004",
			Category:    "High Availability and Resiliency",
			Subcategory: "SKU",
			Description: "Azure OpenAI SKU",
			Severity:    "High",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				i := target.(*armcognitiveservices.Account)
				return false, string(*i.SKU.Name)
			},
			Url: "https://learn.microsoft.com/en-us/azure/ai-services/openai/overview",
		},
		"CAF": {
			Id:          "oai-005",
			Category:    "Governance",
			Subcategory: "Naming Convention (CAF)",
			Description: "Azure OpenAI Name should comply with naming conventions",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcognitiveservices.Account)
				caf := strings.HasPrefix(*c.Name, "oai")
				return !caf, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/cloud-adoption-framework/ready/azure-best-practices/resource-abbreviations",
		},
		"oai-006": {
			Id:          "oai-006",
			Category:    "Governance",
			Subcategory: "Use tags to organize your resources",
			Description: "Azure OpenAI should have tags",
			Severity:    "Low",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcognitiveservices.Account)
				return c.Tags == nil || len(c.Tags) == 0, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/azure-resource-manager/management/tag-resources?tabs=json",
		},
		"oai-007": {
			Id:          "oai-007",
			Category:    "Throttling",
			Subcategory: "Dynamic Throttling",
			Description: "Azure OpenAI should have dynamic throttling enabled",
			Severity:    "Medium",
			Eval: func(target interface{}, scanContext *scanners.ScanContext) (bool, string) {
				c := target.(*armcognitiveservices.Account)
				enabled := c.Properties.DynamicThrottlingEnabled != nil && *c.Properties.DynamicThrottlingEnabled
				return !enabled, ""
			},
			Url: "https://learn.microsoft.com/en-us/azure/ai-services/openai/how-to/quota",
		},
	}
}
