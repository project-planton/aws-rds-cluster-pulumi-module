package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	RdsInstanceEndpoint = "rds-instance-endpoint"
	RdsInstanceId       = "rds-instance-id"
	RdsInstanceArn      = "rds-instance-arn"
	RdsInstanceAddress  = "rds-instance-address"
	RdsSubnetGroup      = "rds-subnet-group"
	RdsSecurityGroup    = "rds-security-group"
	RdsParameterGroup   = "rds-parameter-group"
	RdsOptionsGroup     = "rds-options-group"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *awsaurorapostgres.AwsAuroraPostgresStackInput) *awsaurorapostgres.AwsAuroraPostgresStackOutputs {
	return &awsaurorapostgres.AwsAuroraPostgresStackOutputs{}
}
