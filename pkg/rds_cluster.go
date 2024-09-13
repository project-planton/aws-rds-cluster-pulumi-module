package pkg

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func rdsCluster(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider, createdSecurityGroup *ec2.SecurityGroup) (*rds.Cluster, error) {
	clusterArgs := &rds.ClusterArgs{
		ClusterIdentifier:                pulumi.String(locals.AwsAuroraPostgres.Metadata.Id),
		DatabaseName:                     pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.DatabaseName),
		PreferredMaintenanceWindow:       pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.MaintenanceWindow),
		NetworkType:                      pulumi.String("IPV4"),
		IamDatabaseAuthenticationEnabled: pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.IamDatabaseAuthenticationEnabled),
		Tags:                             pulumi.ToStringMap(locals.Labels),
		Engine:                           pulumi.String("aurora-postgresql"),
		BackupRetentionPeriod:            pulumi.Int(5),
		EngineVersion:                    pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.EngineVersion),
		AllowMajorVersionUpgrade:         pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.AllowMajorVersionUpgrade),
		EngineMode:                       pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.EngineMode),
		Port:                             pulumi.Int(5432),
		PreferredBackupWindow:            pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.BackupWindow),
		CopyTagsToSnapshot:               pulumi.Bool(false),
		ApplyImmediately:                 pulumi.Bool(true),
		EnabledCloudwatchLogsExports:     pulumi.ToStringArray(locals.AwsAuroraPostgres.Spec.RdsCluster.EnabledCloudwatchLogsExports),
		DeletionProtection:               pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.DeletionProtection),
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.DatabasePort > 0 {
		clusterArgs.Port = pulumi.Int(locals.AwsAuroraPostgres.Spec.RdsCluster.DatabasePort)
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.RetentionPeriod > 0 {
		clusterArgs.BackupRetentionPeriod = pulumi.Int(locals.AwsAuroraPostgres.Spec.RdsCluster.RetentionPeriod)
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.ManageMasterUserPassword {
		clusterArgs.ManageMasterUserPassword = pulumi.Bool(true)
		clusterArgs.MasterUserSecretKmsKeyId = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.MasterUserSecretKmsKeyId)
	} else {
		clusterArgs.MasterUsername = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.MasterUser)
		clusterArgs.MasterPassword = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.MasterPassword)
	}

	clusterArgs.SkipFinalSnapshot = pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.SkipFinalSnapshot)
	if !locals.AwsAuroraPostgres.Spec.RdsCluster.SkipFinalSnapshot {
		clusterArgs.FinalSnapshotIdentifier = pulumi.Sprintf("%s-final-snapshot", locals.AwsAuroraPostgres.Metadata.Id)
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.EngineMode != "serverless" {
		clusterArgs.StorageEncrypted = pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.StorageEncrypted)
		if locals.AwsAuroraPostgres.Spec.RdsCluster.StorageEncrypted {
			clusterArgs.KmsKeyId = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.StorageKmsKeyArn)
		}
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration != nil {
		maxCapacity := locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration.MaxCapacity
		if maxCapacity == 0 {
			maxCapacity = 16
		}

		minCapacity := locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration.MinCapacity
		if minCapacity == 0 {
			minCapacity = 2
		}

		secondsUntilAutoPause := locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration.SecondsUntilAutoPause
		if secondsUntilAutoPause == 0 {
			secondsUntilAutoPause = 300
		}

		timeoutAction := locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration.TimeoutAction
		if timeoutAction == "" {
			timeoutAction = "RollbackCapacityChange"
		}

		clusterArgs.ScalingConfiguration = &rds.ClusterScalingConfigurationArgs{
			AutoPause:             pulumi.Bool(locals.AwsAuroraPostgres.Spec.RdsCluster.ScalingConfiguration.AutoPause),
			MaxCapacity:           pulumi.Int(maxCapacity),
			MinCapacity:           pulumi.Int(minCapacity),
			SecondsUntilAutoPause: pulumi.Int(secondsUntilAutoPause),
			TimeoutAction:         pulumi.String(timeoutAction),
		}
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.Serverlessv2ScalingConfiguration != nil {
		clusterArgs.Serverlessv2ScalingConfiguration = &rds.ClusterServerlessv2ScalingConfigurationArgs{
			MaxCapacity: pulumi.Float64(locals.AwsAuroraPostgres.Spec.RdsCluster.Serverlessv2ScalingConfiguration.MaxCapacity),
			MinCapacity: pulumi.Float64(locals.AwsAuroraPostgres.Spec.RdsCluster.Serverlessv2ScalingConfiguration.MinCapacity),
		}
	}

	vpcSecurityGroupIds := pulumi.ToStringArray(locals.AwsAuroraPostgres.Spec.RdsCluster.AssociateSecurityGroupIds)
	vpcSecurityGroupIds = append(vpcSecurityGroupIds, createdSecurityGroup.ID())

	clusterArgs.VpcSecurityGroupIds = vpcSecurityGroupIds

	if len(locals.AwsAuroraPostgres.Spec.RdsCluster.SubnetIds) > 0 && locals.AwsAuroraPostgres.Spec.RdsCluster.DbSubnetGroupName == "" {
		createdSubnetGroup, err := subnetGroup(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create subnet group")
		}
		clusterArgs.DbSubnetGroupName = createdSubnetGroup.Name
	}
	if locals.AwsAuroraPostgres.Spec.RdsCluster.DbSubnetGroupName != "" {
		clusterArgs.DbSubnetGroupName = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.DbSubnetGroupName)
	}

	if locals.AwsAuroraPostgres.Spec.RdsCluster.ClusterParameterGroupName == "" {
		createdClusterParameterGroup, err := clusterParameterGroup(ctx, locals, awsProvider)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create cluster parameter group")
		}
		clusterArgs.DbClusterParameterGroupName = createdClusterParameterGroup.Name
	} else {
		clusterArgs.DbClusterParameterGroupName = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.ClusterParameterGroupName)
	}

	clusterArgs.SnapshotIdentifier = pulumi.String(locals.AwsAuroraPostgres.Spec.RdsCluster.SnapshotIdentifier)

	// Create RDS Cluster
	createdRdsCluster, err := rds.NewCluster(ctx, locals.AwsAuroraPostgres.Metadata.Id, clusterArgs, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create regional rds cluster")
	}

	return createdRdsCluster, nil
}
