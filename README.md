# AWS CloudFront Pulumi Module

**Note:** This module is not completely implemented because the API resource specification is currently empty.

## Introduction

The AWS CloudFront Pulumi Module is designed to provide a unified and standardized way to deploy AWS CloudFront distributions using a Kubernetes-like API resource model. By leveraging Pulumi and our unified APIs, developers can define their infrastructure in simple YAML files, abstracting the complexity of AWS interactions and streamlining the deployment process.

## Key Features

- **Kubernetes-Like API Resource Model**: Mimics Kubernetes API resource structures, making it intuitive for developers familiar with Kubernetes to define and manage AWS CloudFront resources.
  
- **Unified API Structure**: Every resource follows a standard structure with `apiVersion`, `kind`, `metadata`, `spec`, and `status`, ensuring consistency across different resources and cloud providers.
  
- **Pulumi Integration**: Utilizes Pulumi for infrastructure provisioning, allowing for infrastructure as code with the benefits of real programming languages and Pulumi's state management.
  
- **Credential Management**: Securely manages AWS credentials through the `aws_credential_id` field, ensuring that deployments are authenticated and authorized.

- **Extensible and Modular**: Designed to be extended and customized, enabling developers to build upon the module for more complex use cases.

## Architecture

The module operates by accepting an API resource definition as input. It then uses Pulumi to interpret this definition and interact with AWS to create the specified resources. The key components involved are:

- **API Resource Definition**: A YAML file that includes all necessary information to define a CloudFront distribution, following the standard API structure.
  
- **Pulumi Module**: Written in Go, the module reads the API resource and uses Pulumi's AWS SDK to provision resources.
  
- **AWS Provider Initialization**: The module initializes the AWS provider within Pulumi using the credentials specified by `aws_credential_id`.
  
- **Status Reporting**: Outputs from the Pulumi deployment, such as the CloudFront distribution ID, are captured and stored in `status.stackOutputs` for easy reference.

## Usage

Refer to the example section for usage instructions.

## Limitations

- **Incomplete Implementation**: As noted, the module currently lacks a complete implementation of the AWS CloudFront resource creation due to an empty `spec`. Future updates will include full support for defining and deploying CloudFront distributions.

## Contributing

We welcome contributions to enhance the functionality of this module. Please submit pull requests or open issues to help improve the module and its documentation.