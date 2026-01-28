package aws_integrations

import (
	"bytes"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/mackerelio/mkr/mackerelclient"
	"github.com/stretchr/testify/assert"
)

func TestAWSIntegrationsApp_Run(t *testing.T) {
	testCases := []struct {
		id              string
		awsIntegrations []*mackerel.AWSIntegration
		expected        string
	}{
		{
			id: "default",
			awsIntegrations: []*mackerel.AWSIntegration{
				{
					ID:           "4mpaCTV4CPA",
					Name:         "mkr-integration-test-1",
					Memo:         "This is a memo",
					RoleArn:      "arn:aws:iam::123456789123:role/mkr-integration-test-1",
					ExternalID:   "wUCwqnEBFdSoSLr712Lk723fRewdWUT3P8vVxtjx",
					Region:       "ap-southeast-1",
					IncludedTags: "Name:develop-server,Environment:develop",
					ExcludedTags: "Name:staging-server,Environment:staging",
					Services: map[string]*mackerel.AWSIntegrationService{
						"ALB": {
							Enable:          true,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"APIGateway": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Batch": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Billing": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"CloudFront": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Connect": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"DynamoDB": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"EC2": {
							Enable:              true,
							Role:                nil,
							ExcludedMetrics:     []string{},
							RetireAutomatically: true,
						},
						"ECSCluster": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"EFS": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"ELB": {
							Enable:          true,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"ES": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"ElastiCache": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Firehose": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Kinesis": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Lambda": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"NLB": {
							Enable:          true,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"RDS": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Redshift": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"Route53": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"S3": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"SES": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"SQS": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"States": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
						"WAF": {
							Enable:          false,
							Role:            nil,
							ExcludedMetrics: []string{},
						},
					},
				},
			},
			expected: `[
    {
        "id": "4mpaCTV4CPA",
        "name": "mkr-integration-test-1",
        "memo": "This is a memo",
        "roleArn": "arn:aws:iam::123456789123:role/mkr-integration-test-1",
        "externalId": "wUCwqnEBFdSoSLr712Lk723fRewdWUT3P8vVxtjx",
        "region": "ap-southeast-1",
        "includedTags": "Name:develop-server,Environment:develop",
        "excludedTags": "Name:staging-server,Environment:staging",
        "services": {
            "ALB": {
                "enable": true,
                "role": null,
                "excludedMetrics": []
            },
            "APIGateway": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Batch": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Billing": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "CloudFront": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Connect": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "DynamoDB": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "EC2": {
                "enable": true,
                "role": null,
                "excludedMetrics": [],
                "retireAutomatically": true
            },
            "ECSCluster": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "EFS": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "ELB": {
                "enable": true,
                "role": null,
                "excludedMetrics": []
            },
            "ES": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "ElastiCache": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Firehose": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Kinesis": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Lambda": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "NLB": {
                "enable": true,
                "role": null,
                "excludedMetrics": []
            },
            "RDS": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Redshift": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "Route53": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "S3": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "SES": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "SQS": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "States": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            },
            "WAF": {
                "enable": false,
                "role": null,
                "excludedMetrics": []
            }
        }
    }
]
`,
		},
		{
			id:              "default",
			awsIntegrations: []*mackerel.AWSIntegration{},
			expected: `[]
`,
		},
	}

	for _, tc := range testCases {
		client := mackerelclient.NewMockClient(
			mackerelclient.MockFindAWSIntegrations(func() ([]*mackerel.AWSIntegration, error) {
				return tc.awsIntegrations, nil
			}),
		)
		t.Run(tc.id, func(t *testing.T) {
			out := new(bytes.Buffer)
			app := &awsIntegrationsApp{
				client:    client,
				outStream: out,
			}
			assert.NoError(t, app.run(t.Context()))
			assert.Equal(t, tc.expected, out.String())
		})
	}
}
