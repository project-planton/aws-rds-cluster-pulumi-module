package main

import (
	"github.com/plantoncloud/aws-aurora-postgres-pulumi-module/pkg"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/apiresource"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/awscredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/pulumibackendcredential"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/pulumibackendcredential/enums/pulumibackendtype"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/progress/progressstatus"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/resourcemanager/v1/environment"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		s := pkg.ResourceStack{
			Input: &awsaurorapostgres.AwsAuroraPostgresStackInput{
				ApiResource: &awsaurorapostgres.AwsAuroraPostgres{
					ApiVersion: "code2cloud.planton.cloud/v1",
					Kind:       "AwsAuroraPostgres",
					Metadata: &apiresource.ApiResourceMetadata{
						Name: "demo",
						Id:   "pgaur-planton-cloud-aws-module-test-demo",
					},
					Spec: &awsaurorapostgres.AwsAuroraPostgresSpec{
						EnvironmentInfo: &environment.ApiResourceEnvironmentInfo{
							EnvId: os.Getenv("ENV_ID"),
						},
						StackJobSettings: &stackjob.StackJobSettings{
							PulumiBackendCredentialId: os.Getenv("PULUMI_BACKEND_CREDENTIAL_ID"),
							StackJobRunnerId:          os.Getenv("STACK_JOB_RUNNER_ID"),
						},
						RdsCluster: &awsaurorapostgres.AwsAuroraPostgresRdsCluster{
							EngineMode:     "global",
							EngineVersion:  "13.11",
							ClusterFamily:  "aurora-postgresql13",
							MasterUser:     "postgres",
							MasterPassword: "password",
						},
					},
				},
				AwsCredential: &awscredential.AwsCredential{
					Spec: &awscredential.AwsCredentialSpec{
						AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
						SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
						Region:          os.Getenv("AWS_REGION"),
					},
				},
				PulumiBackendCredential: &pulumibackendcredential.PulumiBackendCredential{
					Spec: &pulumibackendcredential.PulumiBackendCredentialSpec{
						HttpBackend: &pulumibackendcredential.PulumiBackendCredentialHttpBackendSpec{
							AccessToken: os.Getenv("PULUMI_ACCESS_TOKEN"),
							ApiUrl:      os.Getenv("PULUMI_API_URL"),
						},
						PulumiBackendType:  pulumibackendtype.PulumiBackendType_http,
						PulumiOrganization: os.Getenv("PULUMI_ORGANIZATION"),
					},
				},
				StackJob: &stackjob.StackJob{
					Metadata: &apiresource.ApiResourceMetadata{
						Id: "pgaur-stack-job",
					},
					Spec: &stackjob.StackJobSpec{
						EnvId:           "planton-cloud-aws-module-test",
						ResourceId:      "pgaur-planton-cloud-aws-module-test-demo",
						PulumiStackName: "pgaur-planton-cloud-aws-module-test-demo",
					},
					Status: &stackjob.StackJobStatus{
						PulumiOperations: &stackjob.StackJobStatusPulumiOperationsStatus{
							Apply: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: true,
							},
							ApplyPreview: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
							Destroy: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
							DestroyPreview: &progressstatus.StackJobProgressPulumiOperationStatus{
								IsRequired: false,
							},
						},
					},
				},
			},
		}
		err := s.Resources(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
