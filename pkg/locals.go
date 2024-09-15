package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/aws/awsaurorapostgres"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsAuroraPostgres *awsaurorapostgres.AwsAuroraPostgres
	Labels            map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsaurorapostgres.AwsAuroraPostgresStackInput) *Locals {
	locals := &Locals{}

	//assign value for the locals variable to make it available across the project
	locals.AwsAuroraPostgres = stackInput.Target

	return locals
}
