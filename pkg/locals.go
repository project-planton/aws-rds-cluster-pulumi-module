package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsrdscluster"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsRdsCluster *awsrdscluster.AwsRdsCluster
	Labels        map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsrdscluster.AwsRdsClusterStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsRdsCluster = stackInput.Target

	return locals
}
