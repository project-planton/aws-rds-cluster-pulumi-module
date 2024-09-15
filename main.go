package main

import (
	"github.com/plantoncloud/aws-aurora-postgres-pulumi-module/pkg"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/apiresource"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/awscredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/pulumibackendcredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob"
	_ "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/pulumioperationtype"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/resourcemanager/v1/environment"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		stackInput := &awsaurorapostgres.AwsAuroraPostgresStackInput{
			Target: &awsaurorapostgres.AwsAuroraPostgres{
				ApiVersion: "code2cloud.planton.cloud/v1",
				Kind:       "AwsAuroraPostgres",
				Metadata: &apiresource.ApiResourceMetadata{
					Name: "demo",
					Id:   "aurpg-planton-cloud-aws-module-test-demo",
				},
				Spec: &awsaurorapostgres.AwsAuroraPostgresSpec{
					EnvironmentInfo: &environment.ApiResourceEnvironmentInfo{
						EnvId: os.Getenv("ENV_ID"),
					},
					AwsCredentialId: "N/A",
					StackJobSettings: &stackjob.StackJobSettings{
						PulumiBackendCredentialId: os.Getenv("PULUMI_BACKEND_CREDENTIAL_ID"),
						StackJobRunnerId:          os.Getenv("STACK_JOB_RUNNER_ID"),
					},
					RdsCluster: &awsaurorapostgres.AwsAuroraPostgresRdsCluster{
						EngineMode:     "provisioned",
						EngineVersion:  "13.11",
						ClusterFamily:  "aurora-postgresql13",
						MasterUser:     "postgres",
						MasterPassword: "password",
						ClusterSize:    1,
						InstanceType:   "db.r5.large",
						AutoScaling: &awsaurorapostgres.AwsAuroraPostgresAutoScaling{
							IsEnabled: true,
						},
						//EnhancedMonitoringRoleEnabled: true,
						//RdsMonitoringInterval:         1,
						//EnhancedMonitoringAttributes:  []string{"postgressql", "monitoring"},
					},
				},
			},
			AwsCredential: &awscredential.AwsCredentialSpec{
				AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Region:          os.Getenv("AWS_REGION"),
			},
			Pulumi: &stackjob.StackPulumiInput{
				Operation: pulumioperationtype.PulumiOperationType_up,
				Preview:   true,
				Project: &stackjob.PulumiProject{
					Name: "planton-cloud-aws-module-test",
				},
				StackName: "aurpg-planton-cloud-aws-module-test-demo",
				Backend: &pulumibackendcredential.PulumiBackendCredentialSpec{
					HttpBackend: &pulumibackendcredential.PulumiBackendCredentialHttpBackendSpec{
						AccessToken: os.Getenv("PULUMI_ACCESS_TOKEN"),
						ApiUrl:      os.Getenv("PULUMI_API_URL"),
					},
					PulumiBackendType:  pulumibackendcredential.PulumiBackendType_http,
					PulumiOrganization: os.Getenv("PULUMI_ORGANIZATION"),
				},
			},
		}
		return pkg.Resources(ctx, stackInput)
	})
}
