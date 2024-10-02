package pkg

import (
	awsrdsclusterv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/provider/aws/awsrdscluster/v1"
	"github.com/pkg/errors"
	"github.com/plantoncloud/aws-rds-cluster-pulumi-module/pkg/outputs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input  *awsrdsclusterv1.AwsRdsClusterStackInput
	Labels map[string]string
}

func Resources(ctx *pulumi.Context, stackInput *awsrdsclusterv1.AwsRdsClusterStackInput) error {
	locals := initializeLocals(ctx, stackInput)
	awsCredential := stackInput.AwsCredential

	//create aws provider using the credentials from the input
	awsProvider, err := aws.NewProvider(ctx,
		"classic-provider",
		&aws.ProviderArgs{
			AccessKey: pulumi.String(awsCredential.AccessKeyId),
			SecretKey: pulumi.String(awsCredential.SecretAccessKey),
			Region:    pulumi.String(awsCredential.Region),
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

	ctx.Export(outputs.RdsClusterIdentifier, createdRdsCluster.ClusterIdentifier)
	ctx.Export(outputs.RdsClusterMasterEndpoint, createdRdsCluster.Endpoint)
	ctx.Export(outputs.RdsClusterReaderEndpoint, createdRdsCluster.ReaderEndpoint)

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
