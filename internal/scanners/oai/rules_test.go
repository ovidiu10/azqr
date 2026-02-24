// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package oai

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cmendible/azqr/internal/scanners"
)

func TestOpenAIScanner_Rules(t *testing.T) {
	type fields struct {
		rule                string
		target              interface{}
		scanContext         *scanners.ScanContext
		diagnosticsSettings scanners.DiagnosticsSettings
	}
	type want struct {
		broken bool
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "OpenAIScanner DiagnosticSettings",
			fields: fields{
				rule: "DiagnosticSettings",
				target: &armcognitiveservices.Account{
					ID: to.StringPtr("test"),
				},
				scanContext: &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{
					HasDiagnosticsFunc: func(resourceId string) (bool, error) {
						return true, nil
					},
				},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "OpenAIScanner SLA",
			fields: fields{
				rule:                "SLA",
				target:              &armcognitiveservices.Account{},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "99.9%",
			},
		},
		{
			name: "OpenAIScanner Private Endpoint",
			fields: fields{
				rule: "Private",
				target: &armcognitiveservices.Account{
					Properties: &armcognitiveservices.AccountProperties{
						PrivateEndpointConnections: []*armcognitiveservices.PrivateEndpointConnection{
							{
								ID: to.StringPtr("test"),
							},
						},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "OpenAIScanner no Private Endpoint",
			fields: fields{
				rule: "Private",
				target: &armcognitiveservices.Account{
					Properties: &armcognitiveservices.AccountProperties{
						PrivateEndpointConnections: []*armcognitiveservices.PrivateEndpointConnection{},
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "OpenAIScanner SKU",
			fields: fields{
				rule: "SKU",
				target: &armcognitiveservices.Account{
					SKU: &armcognitiveservices.SKU{
						Name: to.StringPtr("S0"),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "S0",
			},
		},
		{
			name: "OpenAIScanner CAF",
			fields: fields{
				rule: "CAF",
				target: &armcognitiveservices.Account{
					Name: to.StringPtr("oai-test"),
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "OpenAIScanner CAF broken",
			fields: fields{
				rule: "CAF",
				target: &armcognitiveservices.Account{
					Name: to.StringPtr("myopenai"),
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "OpenAIScanner Tags",
			fields: fields{
				rule: "oai-006",
				target: &armcognitiveservices.Account{
					Tags: map[string]*string{
						"env": to.StringPtr("test"),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "OpenAIScanner Dynamic Throttling Enabled",
			fields: fields{
				rule: "oai-007",
				target: &armcognitiveservices.Account{
					Properties: &armcognitiveservices.AccountProperties{
						DynamicThrottlingEnabled: to.BoolPtr(true),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "OpenAIScanner Dynamic Throttling Disabled",
			fields: fields{
				rule: "oai-007",
				target: &armcognitiveservices.Account{
					Properties: &armcognitiveservices.AccountProperties{
						DynamicThrottlingEnabled: to.BoolPtr(false),
					},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
		{
			name: "OpenAIScanner Dynamic Throttling nil",
			fields: fields{
				rule: "oai-007",
				target: &armcognitiveservices.Account{
					Properties: &armcognitiveservices.AccountProperties{},
				},
				scanContext:         &scanners.ScanContext{},
				diagnosticsSettings: scanners.DiagnosticsSettings{},
			},
			want: want{
				broken: true,
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &OpenAIScanner{
				diagnosticsSettings: tt.fields.diagnosticsSettings,
			}
			rules := s.GetRules()
			b, w := rules[tt.fields.rule].Eval(tt.fields.target, tt.fields.scanContext)
			got := want{
				broken: b,
				result: w,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OpenAIScanner Rule.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
