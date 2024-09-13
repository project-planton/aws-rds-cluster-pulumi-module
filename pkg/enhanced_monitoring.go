package pkg

import (
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func enhancedMonitoring(ctx *pulumi.Context, locals *Locals, awsProvider *aws.Provider) (*iam.Role, error) {
	// Define the IAM policy document for enhanced monitoring
	enhancedMonitoringPolicy, err := iam.GetPolicyDocument(ctx, &iam.GetPolicyDocumentArgs{
		Statements: []iam.GetPolicyDocumentStatement{
			{
				Actions: []string{
					"sts:AssumeRole",
				},
				Effect: to.StringPtr("Allow"),
				Principals: []iam.GetPolicyDocumentStatementPrincipal{
					{
						Type:        "Service",
						Identifiers: []string{"monitoring.rds.amazonaws.com"},
					},
				},
			},
		},
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Errorf("failed to get iam policy document")
	}

	// Create IAM Role for Enhanced Monitoring
	enhancedMonitoringRole, err := iam.NewRole(ctx, "enhanced-monitoring-role", &iam.RoleArgs{
		Name:             pulumi.String(locals.AwsAuroraPostgres.Metadata.Id),
		AssumeRolePolicy: pulumi.String(enhancedMonitoringPolicy.Json),
		Tags:             pulumi.ToStringMap(locals.Labels),
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Errorf("failed to create enhanced monitoring role")
	}

	// Attach Amazon's managed policy for RDS enhanced monitoring
	_, err = iam.NewRolePolicyAttachment(ctx, "enhanced-monitoring-policy-attachment", &iam.RolePolicyAttachmentArgs{
		Role:      enhancedMonitoringRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"),
	}, pulumi.Provider(awsProvider))
	if err != nil {
		return nil, errors.Errorf("failed to create enhanced monitoring policy attachment")
	}

	return enhancedMonitoringRole, nil
}
