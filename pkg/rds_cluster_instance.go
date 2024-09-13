package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func rdsClusterInstance(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider, createdRdsCluster *rds.Cluster) ([]*rds.ClusterInstance, error) {
	clusterInstanceArgs := &rds.ClusterInstanceArgs{
		ClusterIdentifier:          createdRdsCluster.ID(),
		InstanceClass:              pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.InstanceType),
		DbSubnetGroupName:          createdRdsCluster.DbSubnetGroupName,
		PubliclyAccessible:         pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.IsPubliclyAccessible),
		Tags:                       pulumi.ToStringMap(locals.Labels),
		Engine:                     createdRdsCluster.Engine,
		EngineVersion:              createdRdsCluster.EngineVersion,
		AutoMinorVersionUpgrade:    pulumi.Bool(true),
		MonitoringInterval:         pulumi.Int(locals.AwsAuroraPostgres.Spec.RdsCluster.RdsMonitoringInterval),
		ApplyImmediately:           pulumi.Bool(true),
		PreferredMaintenanceWindow: pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.MaintenanceWindow),
		PreferredBackupWindow:      pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.BackupWindow),
		CopyTagsToSnapshot:         pulumi.Bool(false),
		CaCertIdentifier:           pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.CaCertIdentifier),
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.EnhancedMonitoringRoleEnabled {
		enhancedMonitoringIamRole, err := enhancedMonitoring(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create enhanced monitoring iam role")
		}
		clusterInstanceArgs.MonitoringRoleArn = enhancedMonitoringIamRole.Arn
	}

	clusterInstanceArgs.PerformanceInsightsEnabled = pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.IsPerformanceInsightsEnabled)
	if locals.AwsAuroraPostgres.Spec.RdsCluster.IsPerformanceInsightsEnabled {
		clusterInstanceArgs.PerformanceInsightsKmsKeyId = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.PerformanceInsightsKmsKeyId)
	}

	var rdsClusterInstances []*rds.ClusterInstance
	for i := 0; i < int(locals.AwsAuroraPostgres.Spec.RdsCluster.ClusterSize); i++ {
		clusterInstanceIdentifier := fmt.Sprintf("%s-%d", locals.AwsAuroraPostgres.Metadata.Id, i)
		clusterInstanceArgs.Identifier = pulumi.String(clusterInstanceIdentifier)
		// Create RDS Cluster
		createdRdsClusterInstance, err := rds.NewClusterInstance(ctx, clusterInstanceIdentifier,
			clusterInstanceArgs,
			pulumi.Provider(awsProvider), pulumi.Parent(createdRdsCluster), pulumi.IgnoreChanges([]string{
				"engine_version",
			}))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create rds cluster instance")
		}
		rdsClusterInstances = append(rdsClusterInstances, createdRdsClusterInstance)
	}
	return rdsClusterInstances, nil
}
