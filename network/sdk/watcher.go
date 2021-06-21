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
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func getNetworkWatchersClient() armnetwork.NetworkWatchersClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armnetwork.NewNetworkWatchersClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create NetworkWatchers
func CreateNetworkWatcher(ctx context.Context, networkWatcherName string) error {
	client := getNetworkWatchersClient()
	_, err := client.CreateOrUpdate(
		ctx,
		config.GroupName(),
		networkWatcherName,
		armnetwork.NetworkWatcher{
			Resource: armnetwork.Resource{
				Location: to.StringPtr(config.Location()),
			},
		},
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

// Deletes the specified network watcher resource.
func DeleteNetworkWatcher(ctx context.Context, networkWatcherName string) error {
	client := getNetworkWatchersClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), networkWatcherName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets all network watchers by resource group.
func ListNetworkWatcher(ctx context.Context) error {
	client := getNetworkWatchersClient()
	_, err := client.List(ctx, config.GroupName(), nil)

	if err != nil {
		return err
	}

	return nil
}

// Gets all network watchers by subscription.
func ListAllNetworkWatcher(ctx context.Context) error {
	client := getNetworkWatchersClient()
	_, err := client.ListAll(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

//  Gets the specified network watcher by resource group.
func GetNetworkWatcher(ctx context.Context, networkWatcherName string) error {
	client := getNetworkWatchersClient()
	_, err := client.Get(ctx, config.GroupName(), networkWatcherName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets Network Configuration Diagnostic data to help customers understand and debug network behavior. It provides
// detailed information on what security rules were applied to a specified traffic flow and
// the result of evaluating these rules. Customers must provide details of a flow like source, destination, protocol, etc. The API returns whether traffic
// was allowed or denied, the rules evaluated for
// the specified flow and the evaluation results.
func GetNetworkConfigurationDiagnostic(ctx context.Context, networkWatcherName string, networkConfigurationDiagnosticParameters armnetwork.NetworkConfigurationDiagnosticParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginGetNetworkConfigurationDiagnostic(
		ctx,
		config.GroupName(),
		networkWatcherName,
		networkConfigurationDiagnosticParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the configured and effective security group rules on the specified VM.
func GetNetworkVMSecurityRules(ctx context.Context, networkWatcherName string, securityGroupViewParametersParameters armnetwork.SecurityGroupViewParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginGetVMSecurityRules(
		ctx,
		config.GroupName(),
		networkWatcherName,
		securityGroupViewParametersParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Verifies the possibility of establishing a direct TCP connection from a virtual machine to a given endpoint including another
// VM or an arbitrary remote server.
func CheckNetworkConnectivity(ctx context.Context, networkWatcherName string, connectivityParametersParameters armnetwork.ConnectivityParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginCheckConnectivity(
		ctx,
		config.GroupName(),
		networkWatcherName,
		connectivityParametersParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Gets the current network topology by resource group.
func GetNetworkTopology(ctx context.Context, networkWatcherName string, topologyParametersParameters armnetwork.TopologyParameters) error {
	client := getNetworkWatchersClient()
	_, err := client.GetTopology(
		ctx,
		config.GroupName(),
		networkWatcherName,
		topologyParametersParameters,
		nil)

	if err != nil {
		return err
	}

	return nil
}

//Verify IP flow from the specified VM to a location given the currently configured NSG rules.
func VerifyNetworkWatcherIPFlow(ctx context.Context, networkWatcherName string, verificationIPFlowParametersParameters armnetwork.VerificationIPFlowParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginVerifyIPFlow(
		ctx,
		config.GroupName(),
		networkWatcherName,
		verificationIPFlowParametersParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

//  Initiate troubleshooting on a specified resource.
func GetNetworkWatcherTroubleshooting(ctx context.Context, networkWatcherName string, troubleShootingParParameters armnetwork.TroubleshootingParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginGetTroubleshooting(
		ctx,
		config.GroupName(),
		networkWatcherName,
		troubleShootingParParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

//   Get the last completed troubleshooting result on a specified resource.
func GetNetworkWatcherTroubleshootingResult(ctx context.Context, networkWatcherName string, queryTroubleshootingParameters armnetwork.QueryTroubleshootingParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginGetTroubleshootingResult(
		ctx,
		config.GroupName(),
		networkWatcherName,
		queryTroubleshootingParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

//  Queries status of flow log and traffic analytics (optional) on a specified resource.
func GetNetworkWatcherFlowLogStatus(ctx context.Context, networkWatcherName string, flowLogStatusParametersParameters armnetwork.FlowLogStatusParameters) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginGetFlowLogStatus(
		ctx,
		config.GroupName(),
		networkWatcherName,
		flowLogStatusParametersParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

//  Configures flow log and traffic analytics (optional) on a specified resource.
func SetNetworkWatcherFlowLogConfiguration(ctx context.Context, networkWatcherName string, flowLogInformationParameters armnetwork.FlowLogInformation) error {
	client := getNetworkWatchersClient()
	poller, err := client.BeginSetFlowLogConfiguration(
		ctx,
		config.GroupName(),
		networkWatcherName,
		flowLogInformationParameters,
		nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Updates a network watcher tags.
func UpdateNetworkWatcherTags(ctx context.Context, networkWatcherName string, tagsObjectParameters armnetwork.TagsObject) error {
	client := getNetworkWatchersClient()
	_, err := client.UpdateTags(
		ctx,
		config.GroupName(),
		networkWatcherName,
		tagsObjectParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
