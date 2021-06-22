// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"log"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
)

func getVirtualMachinesClient() armcompute.VirtualMachinesClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachinesClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// The operation to create or update a virtual machine. Please note some properties can be set only during virtual machine creation.
func CreateVirtualMachine(ctx context.Context, virtualMachineName string, virtualMachineParameters armcompute.VirtualMachine) (string, error) {
	client := getVirtualMachinesClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		virtualMachineName,
		virtualMachineParameters,
		nil,
	)

	if err != nil {
		return "", err
	}

	resp, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return "", err
	}

	// As a workaround due to issue https://github.com/Azure/azure-sdk-for-go/issues/14730
	if resp.VirtualMachine.ID == nil {
		return poller.RawResponse.Request.URL.Path, nil
	}
	return *resp.VirtualMachine.ID, nil
}

// Retrieves information about the run-time state of a virtual machine.
func InstanceVirtualMachineView(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	_, err := client.InstanceView(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all available virtual machine sizes to which the specified virtual machine can be resized.
func ListVirtualMachineAvailableSizes(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	_, err := client.ListAvailableSizes(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves information about the model view or the instance view of a virtual machine.
func GetVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	_, err := client.Get(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Lists all of the virtual machines in the specified resource group. Use the nextLink property in the response to get the next page of virtual machines.
func ListVirtualMachine(ctx context.Context) error {
	client := getVirtualMachinesClient()
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

// Lists all of the virtual machines in the specified subscription. Use the nextLink property in the response to get the next page of virtual
// machines.
func ListAllVirtualMachine(ctx context.Context) error {
	client := getVirtualMachinesClient()
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

// Gets all the virtual machines under the specified subscription for the specified location.
func ListVirtualMachineByLocation(ctx context.Context) error {
	client := getVirtualMachinesClient()
	pager := client.ListByLocation(config.Location(), nil)
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

//  Run command on the VM.
func RunCommandOnVirtualMachine(ctx context.Context, virtualMachineName string, runCommandInputParameters armcompute.RunCommandInput) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginRunCommand(
		ctx,
		config.GroupName(),
		virtualMachineName,
		runCommandInputParameters,
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

// The operation to restart a virtual machine.
func RestartVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginRestart(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// The operation to power off (stop) a virtual machine. The virtual machine can be restarted with the same provisioned resources. You are
// still charged for this virtual machine.
func VirtualMachinePowerOff(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginPowerOff(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// The operation to start a virtual machine.
func StartVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginStart(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

//  The operation to reapply a virtual machine's state.
func ReapplyVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginReapply(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// Shuts down the virtual machine, moves it to a new node, and powers it back on.
func RedeployVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginRedeploy(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// The operation to update a virtual machine.
func UpdateVirtualMachineTags(ctx context.Context, virtualMachineName string, virtualMachineUpdateParameters armcompute.VirtualMachineUpdate) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		virtualMachineName,
		virtualMachineUpdateParameters,
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

// Sets the OS state of the virtual machine to generalized. It is recommended to sysprep the virtual machine before performing this operation.
// For Windows, please refer to Create a managed image of a generalized VM in Azure [https://docs.microsoft.com/en-us/azure/virtual-machines/windows/capture-image-resource].
// For Linux, please refer to How to create an image of a virtual machine or VHD [https://docs.microsoft.com/en-us/azure/virtual-machines/linux/capture-image].
func GenerializeVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	_, err := client.Generalize(
		ctx,
		config.GroupName(),
		virtualMachineName,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// Shuts down the virtual machine and releases the compute resources. You are not billed for the compute resources that this virtual machine
// uses.
func DeallocateVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginDeallocate(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// The operation to simulate the eviction of spot virtual machine.
func SimulateEvictionVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	_, err := client.SimulateEviction(
		ctx,
		config.GroupName(),
		virtualMachineName,
		nil)

	if err != nil {
		return err
	}

	return nil
}

// The operation to perform maintenance on a virtual machine.
func PerformMaintenanceVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginPerformMaintenance(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// Converts virtual machine disks from blob-based to managed disks. Virtual machine must be stop-deallocated before invoking
// this operation.
func ConvertVirtualMachineToManagedDisk(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginConvertToManagedDisks(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// Reimages the virtual machine which has an ephemeral OS disk back to its initial state.
func ReimageVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	poller, err := client.BeginReimage(
		ctx,
		config.GroupName(),
		virtualMachineName,
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

// The operation to delete a virtual machine.
func DeleteVirtualMachine(ctx context.Context, virtualMachineName string) error {
	client := getVirtualMachinesClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualMachineName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
