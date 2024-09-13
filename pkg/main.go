package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	StackInput *awsaurorapostgres.AwsAuroraPostgresStackInput
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := initializeLocals(ctx, s.StackInput)
	if locals.AwsAuroraPostgres.Spec.RdsCluster == nil {
		return errors.Errorf("RDS Cluster stack failed: Ensure that the 'rds_cluster' field within 'spec' is properly defined before proceeding.")
	}

	awsCredential := s.StackInput.AwsCredential

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
