// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getApplicationGatewaysClient() armnetwork.ApplicationGatewaysClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewApplicationGatewaysClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Creates or updates the specified application gateway
func CreateApplicationGateway(ctx context.Context, applicationGatewayName string, applicationGatewayParameters armnetwork.ApplicationGateway) error {
	client := getApplicationGatewaysClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		applicationGatewayParameters,
		nil,
	)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets Ssl predefined policy with the specified policy name.
func GetApplicationGatewaySSLPredefinedPolicy(ctx context.Context, predefinedPolicyName string) error {
	client := getApplicationGatewaysClient()
	_, err := client.GetSSLPredefinedPolicy(ctx, predefinedPolicyName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all SSL predefined policies for configuring Ssl policy.
func ListApplicationGatewayAvailableSSLPredefinedPolicie(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	pager := client.ListAvailableSSLPredefinedPolicies(nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Lists available Ssl options for configuring Ssl policy
func ListApplicationGatewayAvailableSSLOptions(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	_, err := client.ListAvailableSSLOptions(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets the specified application gateway.
func GetApplicationGateway(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewaysClient()
	_, err := client.Get(ctx, config.GroupName(), applicationGatewayName, nil)
	if err != nil {
		return err
	}
	return nil
}

//  Lists all application gateways in a resource group.
func ListApplicationGateway(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	pager := client.List(config.GroupName(), nil)

	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Lists all available server variables
func ListApplicationGatewayAvailableServerVariables(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	_, err := client.ListAvailableServerVariables(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all available response headers
func ListApplicationGatewayAvailableResponseHeaders(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	_, err := client.ListAvailableResponseHeaders(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all available request headers
func ListApplicationGatewayAvailableRequestHeaders(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	_, err := client.ListAvailableRequestHeaders(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all available web application firewall rule sets
func ListApplicationGatewayAvailableWafRuleSets(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	_, err := client.ListAvailableWafRuleSets(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets all the application gateways in a subscription.
func ListAllApplicationGateway(ctx context.Context) error {
	client := getApplicationGatewaysClient()
	pager := client.ListAll(nil)
	for pager.NextPage(ctx) {
		if pager.Err() != nil {
			return pager.Err()
		}
	}

	if pager.Err() != nil {
		return pager.Err()
	}
	return nil
}

// Gets the backend health for given combination of backend pool and http setting of the specified application gateway in a
// resource group
func GetApplicationGatewayBackendHealthOnDemand(ctx context.Context, applicationGatewayName string, probeRequestParameters armnetwork.ApplicationGatewayOnDemandProbe) error {
	client := getApplicationGatewaysClient()
	poller, err := client.BeginBackendHealthOnDemand(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		probeRequestParameters,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the backend health of the specified application gateway in a resource group
func GetApplicationGatewayBackendHealth(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewaysClient()
	poller, err := client.BeginBackendHealth(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Starts the specified application gateway
func StartApplicationGateway(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewaysClient()
	poller, err := client.BeginStart(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Stops the specified application gateway in a resource group
func StopApplicationGateway(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewaysClient()
	poller, err := client.BeginStop(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Updates the specified application gateway tags.
func UpdateApplicationGatewayTags(ctx context.Context, applicationGatewayName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getApplicationGatewaysClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		applicationGatewayName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes the specified application gateway.
func DeleteApplicationGateway(ctx context.Context, applicationGatewayName string) error {
	client := getApplicationGatewaysClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), applicationGatewayName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
