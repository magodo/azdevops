//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pipeline

import (
	"errors"
	"reflect"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// NewPipeline creates a pipeline from connection options. Policies from ClientOptions are
// placed after policies from PipelineOptions.
func NewPipeline(module, version string, cred azcore.TokenCredential, plOpts runtime.PipelineOptions, cliOpts *policy.ClientOptions) (runtime.Pipeline, error) {
	if cliOpts == nil {
		cliOpts = &policy.ClientOptions{}
	}
	conf, err := getConfiguration(cliOpts)
	if err != nil {
		return runtime.Pipeline{}, err
	}
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{conf.Audience}, &policy.BearerTokenOptions{
		InsecureAllowCredentialWithHTTP: cliOpts.InsecureAllowCredentialWithHTTP,
	})
	// we don't want to modify the underlying array in plOpts.PerRetry
	perRetry := make([]policy.Policy, len(plOpts.PerRetry), len(plOpts.PerRetry)+1)
	copy(perRetry, plOpts.PerRetry)
	perRetry = append(perRetry, authPolicy)
	plOpts.PerRetry = perRetry
	if plOpts.APIVersion.Name == "" {
		plOpts.APIVersion.Name = "api-version"
	}
	return runtime.NewPipeline(module, version, plOpts, cliOpts), nil
}

func getConfiguration(o *policy.ClientOptions) (cloud.ServiceConfiguration, error) {
	c := AzurePublic
	if !reflect.ValueOf(o.Cloud).IsZero() {
		c = o.Cloud
	}
	if conf, ok := c.Services[SerivceAzureDevOps]; ok && conf.Audience != "" {
		return conf, nil
	} else {
		return conf, errors.New("provided Cloud field is missing Azure DevOps configuration")
	}
}
