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

func getVirtualMachineScaleSetVmsClient() armcompute.VirtualMachineScaleSetVMsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetVMsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Gets the status of a virtual machine from a VM scale set.
func GetVirtualMachineScaleSetVmInstanceView(ctx context.Context, vmScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	_, err := client.GetInstanceView(ctx, config.GroupName(), vmScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}
	return nil
}

// Gets the status of a virtual machine from a VM scale set.
func ListVirtualMachineScaleSetVm(ctx context.Context, vmScaleSetName string) error {
	client := getVirtualMachineScaleSetVmsClient()
	pager := client.List(config.GroupName(), vmScaleSetName, nil)

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

// Gets a virtual machine from a VM scale set
func GetVirtualMachineScaleSetVm(ctx context.Context, vmScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	_, err := client.Get(ctx, config.GroupName(), vmScaleSetName, instanceId, nil)
	if err != nil {
		return err
	}
	return nil
}

// Updates a virtual machine of a VM scale set.
func UpdateVirtualMachineScaleSetVm(ctx context.Context, vmScaleSetName string, instanceId string,
	virtualMachineScaleSetVMParameters armcompute.VirtualMachineScaleSetVM) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		vmScaleSetName,
		instanceId,
		virtualMachineScaleSetVMParameters,
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

// Power off (stop) a virtual machine in a VM scale set. Note that resources are still attached and you are getting charged for the resources.
// Instead, use deallocate to release resources and avoid
// charges.
func VirtualMachineScaleSetVmPowerOff(ctx context.Context, vmScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginPowerOff(ctx, config.GroupName(), vmScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Starts a virtual machine in a VM scale set
func StartVirtualMachineScaleSetVm(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginStart(ctx, config.GroupName(), virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Restarts a virtual machine in a VM scale set
func RestartVirtualMachineScaleSetVm(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginRestart(ctx, config.GroupName(), virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Run command on a virtual machine in a VM scale set
func RunCommandOnVirtualMachineScaleSetVm(ctx context.Context, virtualMachineName string, instanceId string, runCommandInputParameters armcompute.RunCommandInput) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginRunCommand(
		ctx,
		config.GroupName(),
		virtualMachineName,
		instanceId,
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

// Deallocates a specific virtual machine in a VM scale set. Shuts down the virtual machine and releases the compute resources it uses.
// You are not billed for the compute resources of this virtual
// machine once it is deallocated
func DeallocateVirtualMachineScaleSetVm(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginDeallocate(ctx, config.GroupName(), virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// huts down the virtual machine in the virtual machine scale set, moves it to a new node, and powers it back on.
func RedeployVirtualMachineScaleSetVm(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginRedeploy(ctx, config.GroupName(), virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Reimages (upgrade the operating system) a specific virtual machine in a VM scale set.
func VirtualMachineScaleSetVmReimage(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginReimage(ctx, config.GroupName(), virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Allows you to re-image all the disks ( including data disks ) in the a VM scale set instance. This operation is only supported for
// managed disks.
func ReimageAllVirtualMachineScaleSetVm(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginReimageAll(ctx, config.GroupName(),
		virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// maintenance on a virtual machine in a VM scale set
func VirtualMachineScaleSetVmPerformMaintenance(ctx context.Context, virtualMachineScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	poller, err := client.BeginPerformMaintenance(ctx, config.GroupName(),
		virtualMachineScaleSetName, instanceId, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

//  Deletes a virtual machine from a VM scale set.
func DeleteVirtualMachineScaleSetVm(ctx context.Context, vmScaleSetName string, instanceId string) error {
	client := getVirtualMachineScaleSetVmsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), vmScaleSetName, instanceId, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}
