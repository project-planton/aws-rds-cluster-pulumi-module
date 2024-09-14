package pkg

import (
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *awsaurorapostgres.AwsAuroraPostgresStackInput
	Labels map[string]string
}

func Resources(ctx *pulumi.Context, stackInput *awsaurorapostgres.AwsAuroraPostgresStackInput) error {
	v, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(stackInput.ApiResource.Spec),
	)
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(stackInput.ApiResource.Spec); err != nil {
		return errors.Errorf("%s", err)
	}

	locals := initializeLocals(ctx, stackInput)
	awsCredential := stackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"classic-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.Spec.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.Spec.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Spec.Region),
		})
	if err != nil {
		return errors.Wrap(err, "failed to create aws provider")
	}

	createdSecurityGroup, err := securityGroup(ctx, locals, awsProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create default security group")
	}

	// Create RDS Cluster
	createdRdsCluster, err := rdsCluster(ctx, locals, awsProvider, createdSecurityGroup)
	if err != nil {
		return errors.Wrap(err, "failed to create rds cluster")
	}

	// Create RDS Cluster Instance
	_, err = rdsClusterInstance(ctx, locals, awsProvider, createdRdsCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create rds cluster instances")
	}

	// Create RDS Cluster Instance
	err = appAutoscaling(ctx, locals, awsProvider, createdRdsCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create auto scaling policy")
	}

	return nil
}
