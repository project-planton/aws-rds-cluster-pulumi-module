package main

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-rds-cluster-pulumi-module/pkg"
	"github.com/plantoncloud/planton/apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrdscluster"
	_ "github.com/plantoncloud/planton/apis/zzgo/cloud/planton/apis/iac/v1/stackjob"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsrdscluster.AwsRdsClusterStackInput{}
		//stackInput := &awsrdscluster.AwsRdsClusterStackInput{
		//	Target: &awsrdscluster.AwsRdsCluster{
		//		ApiVersion: "code2cloud.planton.cloud/v1",
		//		Kind:       "AwsRdsCluster",
		//		Metadata: &apiresource.ApiResourceMetadata{
		//			Name: "demo",
		//			Id:   "aurpg-planton-cloud-aws-module-test-demo",
		//		},
		//		Spec: &awsrdscluster.AwsRdsClusterSpec{
		//			EnvironmentInfo: &environment.ApiResourceEnvironmentInfo{
		//				EnvId: os.Getenv("ENV_ID"),
		//			},
		//			AwsCredentialId: "N/A",
		//			StackJobSettings: &stackjob.StackJobSettings{
		//				PulumiBackendCredentialId: os.Getenv("PULUMI_BACKEND_CREDENTIAL_ID"),
		//				StackJobRunnerId:          os.Getenv("STACK_JOB_RUNNER_ID"),
		//			},
		//			EngineMode:     "provisioned",
		//			EngineVersion:  "13.11",
		//			ClusterFamily:  "aurora-postgresql13",
		//			MasterUser:     "postgres",
		//			MasterPassword: "password",
		//			ClusterSize:    1,
		//			InstanceType:   "db.r5.large",
		//			AutoScaling: &awsrdscluster.AwsRdsClusterAutoScaling{
		//				IsEnabled: true,
		//			},
		//			//EnhancedMonitoringRoleEnabled: true,
		//			//RdsMonitoringInterval:         1,
		//			//EnhancedMonitoringAttributes:  []string{"postgressql", "monitoring"},
		//		},
		//	},
		//	AwsCredential: &awscredential.AwsCredentialSpec{
		//		AccessKeyId:     os.Getenv("AWS_ACCESS_KEY_ID"),
		//		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		//		Region:          os.Getenv("AWS_REGION"),
		//	},
		//	Pulumi: &stackjob.StackPulumiInput{
		//		Operation: pulumioperationtype.PulumiOperationType_up,
		//		Preview:   true,
		//		Project: &stackjob.PulumiProject{
		//			Name: "planton-cloud-aws-module-test",
		//		},
		//		StackName: "aurpg-planton-cloud-aws-module-test-demo",
		//		Backend: &pulumibackendcredential.PulumiBackendCredentialSpec{
		//			Http: &pulumibackendcredential.PulumiHttpBackend{
		//				AccessToken: os.Getenv("PULUMI_ACCESS_TOKEN"),
		//				ApiUrl:      os.Getenv("PULUMI_API_URL"),
		//			},
		//			PulumiBackendType:  pulumibackendcredential.PulumiBackendType_http,
		//			PulumiOrganization: os.Getenv("PULUMI_ORGANIZATION"),
		//		},
		//	},
		//}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return pkg.Resources(ctx, stackInput)
	})
}
