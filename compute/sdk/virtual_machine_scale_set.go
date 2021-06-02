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
	"github.com/Azure/azure-sdk-for-go/sdk/arm/compute/2020-09-30/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func getVirtualMachineScaleSetsClient() armcompute.VirtualMachineScaleSetsClient {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	client := armcompute.NewVirtualMachineScaleSetsClient(armcore.NewDefaultConnection(cred, nil), config.SubscriptionID())
	return *client
}

// Create or update a VM scale set.
func CreateVirtualMachineScaleSet(ctx context.Context, vmScaleSetName string, virtualMachineScaleSetParameters armcompute.VirtualMachineScaleSet) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginCreateOrUpdate(
		ctx,
		config.GroupName(),
		vmScaleSetName,
		virtualMachineScaleSetParameters,
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

// Deletes a VM scale set.
func DeleteVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	resp, err := client.BeginDelete(ctx, config.GroupName(), virtualMachineScaleSetName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Upgrades one or more virtual machines to the latest SKU set in the VM scale set model.
func UpdateVirtualMachineScaleSetInstance(ctx context.Context, virtualMachineScaleSetName string, vmInstanceIDs armcompute.VirtualMachineScaleSetVMInstanceRequiredIDs) error {
	client := getVirtualMachineScaleSetsClient()
	_, err := client.BeginUpdateInstances(
		ctx,
		config.GroupName(),
		virtualMachineScaleSetName,
		vmInstanceIDs,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Deletes virtual machines in a VM scale set.
func DeleteVirtualMachineScaleSetInstance(ctx context.Context, virtualMachineScaleSetName string, vmInstanceIDs armcompute.VirtualMachineScaleSetVMInstanceRequiredIDs) error {
	client := getVirtualMachineScaleSetsClient()
	_, err := client.BeginDeleteInstances(
		ctx,
		config.GroupName(),
		virtualMachineScaleSetName,
		vmInstanceIDs,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

// Get - Display information about a virtual machine scale set.
func GetVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	_, err := client.Get(ctx, config.GroupName(), virtualMachineScaleSetName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets list of OS upgrades on a VM scale set instance.
func GetVirtualMachineScaleSetOSUpgradeHistory(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	pager := client.GetOSUpgradeHistory(config.GroupName(), virtualMachineScaleSetName, nil)
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

// Gets the status of a VM scale set instance.
func GetVirtualMachineScaleSetInstanceView(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	_, err := client.GetInstanceView(ctx, config.GroupName(), virtualMachineScaleSetName, nil)
	if err != nil {
		return err
	}
	return nil
}

// Gets a list of all VM scale sets under a resource group.
func ListVirtualMachineScaleSet(ctx context.Context) error {
	client := getVirtualMachineScaleSetsClient()
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

// Gets a list of all VM Scale Sets in the subscription, regardless of the associated resource group. Use nextLink property in the response to
// get the next page of VM Scale Sets. Do this till nextLink is
// null to fetch all the VM Scale Sets.
func ListAllVirtualMachineScaleSet(ctx context.Context) error {
	client := getVirtualMachineScaleSetsClient()
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

// Gets a list of SKUs available for your VM scale set, including the minimum and maximum VM instances allowed for each SKU.
func ListVirtualMachineScaleSetSKU(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	pager := client.ListSKUs(config.GroupName(), virtualMachineScaleSetName, nil)

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

// Restarts one or more virtual machines in a VM scale set.
func RestartVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginRestart(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Power off (stop) one or more virtual machines in a VM scale set. Note that resources are still attached and you are getting charged for
// the resources. Instead, use deallocate to release resources and
// avoid charges.
func VirtualMachineScaleSetPowerOff(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginPowerOff(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Starts one or more virtual machines in a VM scale set.
func StartVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginStart(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Shuts down all the virtual machines in the virtual machine scale set, moves them to a new node, and powers them back on.
func RedeployVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginRedeploy(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Deallocates specific virtual machines in a VM scale set. Shuts down the virtual machines and releases the compute resources. You are
// not billed for the compute resources that this virtual machine
// scale set deallocates.
func DeallocateVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginDeallocate(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Changes ServiceState property for a given service
func SetVirtualMachineScaleSetOrchestrationServiceState(ctx context.Context, virtualMachineScaleSetName string,
	orchestrationServiceStateInputParameters armcompute.OrchestrationServiceStateInput) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginSetOrchestrationServiceState(ctx, config.GroupName(),
		virtualMachineScaleSetName, orchestrationServiceStateInputParameters, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Reimages (upgrade the operating system) one or more virtual machines in a VM scale set which don't have a ephemeral OS disk, for virtual
// machines who have a ephemeral OS disk the virtual machine is
// reset to initial state.
func VirtualMachineScaleSetReimage(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginReimage(ctx, config.GroupName(), virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Reimages all the disks ( including data disks ) in the virtual machines in a VM scale set. This operation is only supported for managed
// disks.
func ReimageAllVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginReimageAll(ctx, config.GroupName(),
		virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Perform maintenance on one or more virtual machines in a VM scale set. Operation on instances which are not eligible for perform
// maintenance will be failed. Please refer to best practices for more
// details: https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/virtual-machine-scale-sets-maintenance-notifications
func VirtualMachineScaleSetPerformMaintenance(ctx context.Context, virtualMachineScaleSetName string) error {
	client := getVirtualMachineScaleSetsClient()
	poller, err := client.BeginPerformMaintenance(ctx, config.GroupName(),
		virtualMachineScaleSetName, nil)

	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

// Update a VM scale set.
func UpdateVirtualMachineScaleSet(ctx context.Context, virtualMachineScaleSetName string, virtualMachineScaleSetUpdateParameters armcompute.VirtualMachineScaleSetUpdate) error {
	client := getVirtualMachineScaleSetsClient()
	_, err := client.BeginUpdate(
		ctx,
		config.GroupName(),
		virtualMachineScaleSetName,
		virtualMachineScaleSetUpdateParameters,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}
