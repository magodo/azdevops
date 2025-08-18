//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pipeline

import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

const SerivceAzureDevOps cloud.ServiceName = "azuredevops"

var AzureChina = cloud.Configuration{
	ActiveDirectoryAuthorityHost: "https://login.chinacloudapi.cn/",
	Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
		SerivceAzureDevOps: {
			// TODO: Confirm the scope is correct
			Audience: "499b84ac-1321-427f-aa17-267ca6975798",
		},
	},
}

var AzureGovernment = cloud.Configuration{
	ActiveDirectoryAuthorityHost: "https://login.microsoftonline.us/",
	Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
		SerivceAzureDevOps: {
			// TODO: Confirm the scope is correct
			Audience: "499b84ac-1321-427f-aa17-267ca6975798",
		},
	},
}

var AzurePublic = cloud.Configuration{
	ActiveDirectoryAuthorityHost: "https://login.microsoftonline.com/",
	Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
		SerivceAzureDevOps: {
			Audience: "499b84ac-1321-427f-aa17-267ca6975798",
		},
	},
}
