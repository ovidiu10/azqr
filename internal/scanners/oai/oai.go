// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package oai

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"github.com/cmendible/azqr/internal/scanners"
)

// OpenAIScanner - Scanner for Azure OpenAI Services
type OpenAIScanner struct {
	config              *scanners.ScannerConfig
	diagnosticsSettings scanners.DiagnosticsSettings
	accountsClient      *armcognitiveservices.AccountsClient
	listAccountsFunc    func(resourceGroupName string) ([]*armcognitiveservices.Account, error)
}

// Init - Initializes the OpenAIScanner
func (a *OpenAIScanner) Init(config *scanners.ScannerConfig) error {
	a.config = config
	var err error
	a.accountsClient, err = armcognitiveservices.NewAccountsClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	if err != nil {
		return err
	}
	a.diagnosticsSettings = scanners.DiagnosticsSettings{}
	err = a.diagnosticsSettings.Init(config)
	if err != nil {
		return err
	}
	return nil
}

// Scan - Scans all Azure OpenAI Services in a Resource Group
func (a *OpenAIScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	log.Printf("Scanning Azure OpenAI Services in Resource Group %s", resourceGroupName)

	accounts, err := a.listAccounts(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := a.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, account := range accounts {
		if account.Kind == nil || *account.Kind != "OpenAI" {
			continue
		}

		rr := engine.EvaluateRules(rules, account, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID: a.config.SubscriptionID,
			ResourceGroup:  resourceGroupName,
			ServiceName:    *account.Name,
			Type:           *account.Type,
			Location:       *account.Location,
			Rules:          rr,
		})
	}
	return results, nil
}

func (a *OpenAIScanner) listAccounts(resourceGroupName string) ([]*armcognitiveservices.Account, error) {
	if a.listAccountsFunc == nil {
		pager := a.accountsClient.NewListByResourceGroupPager(resourceGroupName, nil)

		accounts := make([]*armcognitiveservices.Account, 0)
		for pager.More() {
			resp, err := pager.NextPage(a.config.Ctx)
			if err != nil {
				return nil, err
			}
			accounts = append(accounts, resp.Value...)
		}
		return accounts, nil
	}

	return a.listAccountsFunc(resourceGroupName)
}
