// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package compute

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	network "github.com/Azure-Samples/azure-sdk-for-go-samples/network/sdk"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestVirtualMachineScaleSet(t *testing.T) {
	groupName := config.GenerateGroupName("compute")
	config.SetGroupName(groupName)

	virtualMachineScaleSetName := config.AppendRandomSuffix("virtualmachinescaleset")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subNetName := config.AppendRandomSuffix("subnet")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	loadBalancingRuleName := config.AppendRandomSuffix("loadbalancingrule")
	outBoundRuleName := config.AppendRandomSuffix("outboundrule")
	probeName := "probe"
	frontendIpConfigurationName := config.AppendRandomSuffix("frontendipconfiguration")
	backendAddressPoolName := config.AppendRandomSuffix("backendaddresspool")

	ctx, cancel := context.WithTimeout(context.Background(), 8000*time.Second)
	defer cancel()
	defer resources.Cleanup(ctx)

	_, err := resources.CreateGroup(ctx, groupName)
	if err != nil {
		t.Fatalf("failed to create group: %+v", err)
	}

	virtualNetworkParameters := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = network.CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.1.0/24"),
		},
	}
	subnetId, err := network.CreateSubnet(ctx, virtualNetworkName, subNetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},

		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			IdleTimeoutInMinutes:     to.Int32Ptr(10),
			PublicIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
		SKU: &armnetwork.PublicIPAddressSKU{
			Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr(),
		},
	}
	publicIpAddressId, err := network.CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	loadBalancerUrl := "/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/loadBalancers/" + loadBalancerName
	loadBalancerParameters := armnetwork.LoadBalancer{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			BackendAddressPools: []*armnetwork.BackendAddressPool{
				{
					Name: &backendAddressPoolName,
				},
			},
			FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
				{
					Name: &frontendIpConfigurationName,
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							Resource: armnetwork.Resource{
								ID: &publicIpAddressId,
							},
						},
					},
				},
			},
			LoadBalancingRules: []*armnetwork.LoadBalancingRule{
				{
					Name: &loadBalancingRuleName,
					Properties: &armnetwork.LoadBalancingRulePropertiesFormat{
						BackendAddressPool: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/backendAddressPools/" + backendAddressPoolName),
						},
						BackendPort:         to.Int32Ptr(80),
						DisableOutboundSnat: to.BoolPtr(true),
						EnableFloatingIP:    to.BoolPtr(true),
						EnableTCPReset:      new(bool),
						FrontendIPConfiguration: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
						},
						FrontendPort:         to.Int32Ptr(80),
						IdleTimeoutInMinutes: to.Int32Ptr(15),
						LoadDistribution:     armnetwork.LoadDistributionDefault.ToPtr(),
						Probe: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/probes/" + probeName),
						},
						Protocol: armnetwork.TransportProtocolTCP.ToPtr(),
					},
				},
			},
			OutboundRules: []*armnetwork.OutboundRule{
				{
					Name: &outBoundRuleName,
					Properties: &armnetwork.OutboundRulePropertiesFormat{
						BackendAddressPool: &armnetwork.SubResource{
							ID: to.StringPtr(loadBalancerUrl + "/backendAddressPools/" + backendAddressPoolName),
						},
						FrontendIPConfigurations: []*armnetwork.SubResource{
							{
								ID: to.StringPtr(loadBalancerUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
							},
						},
						Protocol: armnetwork.LoadBalancerOutboundRuleProtocolAll.ToPtr(),
					},
				},
			},
			Probes: []*armnetwork.Probe{
				{
					Name: &probeName,
					Properties: &armnetwork.ProbePropertiesFormat{
						IntervalInSeconds: to.Int32Ptr(15),
						NumberOfProbes:    to.Int32Ptr(2),
						Port:              to.Int32Ptr(80),
						Protocol:          armnetwork.ProbeProtocolHTTP.ToPtr(),
						RequestPath:       to.StringPtr("healthcheck.aspx"),
					},
				},
			},
		},
		SKU: &armnetwork.LoadBalancerSKU{
			Name: armnetwork.LoadBalancerSKUNameStandard.ToPtr(),
		},
	}

	loadBalancerId, err := network.CreateLoadBalancer(ctx, loadBalancerName, loadBalancerParameters)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}

	probeUri := loadBalancerId + "/probes/" + probeName
	backedPoolsUri := loadBalancerId + "/backendAddressPools/" + backendAddressPoolName
	virtualMachineScaleSetParameters := armcompute.VirtualMachineScaleSet{
		Resource: armcompute.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armcompute.VirtualMachineScaleSetProperties{
			Overprovision: to.BoolPtr(true),
			UpgradePolicy: &armcompute.UpgradePolicy{
				Mode: armcompute.UpgradeModeManual.ToPtr(),
			},
			AutomaticRepairsPolicy: &armcompute.AutomaticRepairsPolicy{
				Enabled:     to.BoolPtr(true),
				GracePeriod: to.StringPtr("PT30M"),
			},
			VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
				NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
					NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{{
						Name: to.StringPtr("testPC"),
						Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
							EnableIPForwarding: to.BoolPtr(true),
							IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{{
								Name: to.StringPtr("testPC"),
								Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
									Subnet: &armcompute.APIEntityReference{
										ID: &subnetId,
									},
									LoadBalancerBackendAddressPools: []*armcompute.SubResource{
										{ID: &backedPoolsUri},
									},
								},
							}},
							Primary: to.BoolPtr(true),
						},
					}},
					HealthProbe: &armcompute.APIEntityReference{ID: &probeUri},
				},

				OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
					AdminPassword:      to.StringPtr("Aa!1()-xyz"),
					AdminUsername:      to.StringPtr("testuser"),
					ComputerNamePrefix: to.StringPtr("testPC"),
				},
				StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
					ImageReference: &armcompute.ImageReference{
						Offer:     to.StringPtr("WindowsServer"),
						Publisher: to.StringPtr("MicrosoftWindowsServer"),
						SKU:       to.StringPtr("2016-Datacenter"),
						Version:   to.StringPtr("latest"),
					},
					OSDisk: &armcompute.VirtualMachineScaleSetOSDisk{
						Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
						CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
						DiskSizeGB:   to.Int32Ptr(512),
						ManagedDisk: &armcompute.VirtualMachineScaleSetManagedDiskParameters{
							StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(),
						},
					},
				},
			},
		},
		SKU: &armcompute.SKU{
			Capacity: to.Int64Ptr(2),
			Name:     to.StringPtr("Standard_D1_v2"),
			Tier:     to.StringPtr("Standard"),
		},
	}

	err = CreateVirtualMachineScaleSet(ctx, virtualMachineScaleSetName, virtualMachineScaleSetParameters)
	if err != nil {
		t.Fatalf("failed to create virtual machine scale set: % +v", err)
	}
	t.Logf("created virtual machine scale set")

	orchestrationServiceStateInputParameters := armcompute.OrchestrationServiceStateInput{
		Action:      armcompute.OrchestrationServiceStateActionSuspend.ToPtr(),
		ServiceName: armcompute.OrchestrationServiceNamesAutomaticRepairs.ToPtr(),
	}
	err = SetVirtualMachineScaleSetOrchestrationServiceState(ctx, virtualMachineScaleSetName, orchestrationServiceStateInputParameters)
	if err != nil {
		t.Fatalf("failed to change ServiceState property for a given service: % +v", err)
	}
	t.Logf("changed ServiceState property for a given service")

	err = VirtualMachineScaleSetReimage(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to reimage one or more virtual machines in a VM scale set: % +v", err)
	}
	t.Logf("reimaged one or more virtual machines in a VM scale set")

	err = ReimageAllVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to reimage all the disks ( including data disks ) in the virtual machines in a VM scale set: % +v", err)
	}
	t.Logf("reimaged all the disks ( including data disks ) in the virtual machines in a VM scale set")

	// Do not test from feedback
	//  "Operation 'performMaintenance' is not allowed on VM since the Subscription of this VM is not eligible."
	// err = VirtualMachineScaleSetPerformMaintenance(ctx, virtualMachineScaleSetName)
	// if err != nil {
	// 	t.Fatalf("failed to perform maintenance on one or more virtual machines in a VM scale set: % +v", err)
	// }
	// t.Logf("performed maintenance on one or more virtual machines in a VM scale set")

	instanceId := 0
	for i := 0; i < 4; i++ {
		instanceId = i
		err = GetVirtualMachineScaleSetVmInstanceView(ctx, virtualMachineScaleSetName, strconv.Itoa(instanceId))
		if err != nil {
			if instanceId >= 3 {
				t.Fatalf("failed to redeploy a virtual machines in a VM scale set: %+v", err)
			}
			continue
		}
		break
	}
	t.Logf("got the status of a virtual machine from a VM scale set by instanceid: %+v", instanceId)

	vmInstanceIDs := armcompute.VirtualMachineScaleSetVMInstanceRequiredIDs{
		InstanceIDs: []*string{to.StringPtr(strconv.Itoa(instanceId))},
	}
	err = UpdateVirtualMachineScaleSetInstance(ctx, virtualMachineScaleSetName, vmInstanceIDs)
	if err != nil {
		t.Fatalf("failed to update virtual machine scale instance: %+v", err)
	}
	t.Logf("updated virtual machine scale instance")

	vmInstanceIDs = armcompute.VirtualMachineScaleSetVMInstanceRequiredIDs{
		InstanceIDs: []*string{to.StringPtr(strconv.Itoa(instanceId))},
	}
	err = DeleteVirtualMachineScaleSetInstance(ctx, virtualMachineScaleSetName, vmInstanceIDs)
	if err != nil {
		t.Fatalf("failed to update virtual machine scale instance: %+v", err)
	}
	t.Logf("updated virtual machine scale instance")

	err = GetVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to get virtual machine scale set: %+v", err)
	}
	t.Logf("got virtual machine scale set")

	err = GetVirtualMachineScaleSetOSUpgradeHistory(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to get list of OS upgrades on a VM scale set instance: %+v", err)
	}
	t.Logf("got list of OS upgrades on a VM scale set instance")

	err = GetVirtualMachineScaleSetInstanceView(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to get the status of a VM scale set instance: %+v", err)
	}
	t.Logf("got the status of a VM scale set instance")

	err = ListVirtualMachineScaleSet(ctx)
	if err != nil {
		t.Fatalf("failed to list virtual machine scale set: %+v", err)
	}
	t.Logf("listed virtual machine scale set")

	err = ListAllVirtualMachineScaleSet(ctx)
	if err != nil {
		t.Fatalf("failed to list all virtual machine scale set: %+v", err)
	}
	t.Logf("listed all virtual machine scale set")

	err = ListVirtualMachineScaleSetSKU(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to list SKUs available for VM scale set: %+v", err)
	}
	t.Logf("listed SKUs available for VM scale set")

	virtualMachineScaleSetUpdateParameters := armcompute.VirtualMachineScaleSetUpdate{
		Properties: &armcompute.VirtualMachineScaleSetUpdateProperties{
			UpgradePolicy: &armcompute.UpgradePolicy{
				Mode: armcompute.UpgradeModeManual.ToPtr(),
			},
		},
		SKU: &armcompute.SKU{
			Capacity: to.Int64Ptr(2),
			Name:     to.StringPtr("Standard_D1_v2"),
			Tier:     to.StringPtr("Standard"),
		},
	}
	err = UpdateVirtualMachineScaleSet(ctx, virtualMachineScaleSetName, virtualMachineScaleSetUpdateParameters)
	if err != nil {
		t.Fatalf("failed to update VM scale set: %+v", err)
	}
	t.Logf("updated VM scale set")

	err = RestartVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to restart one or more virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("restarted one or more virtual machines in a VM scale set")

	err = VirtualMachineScaleSetPowerOff(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to power off one or more virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("powered off one or more virtual machines in a VM scale set")

	err = StartVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to start one or more virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("start one or more virtual machines in a VM scale set")

	err = RedeployVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to redeploy virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("redeplied virtual machines in a VM scale set")

	err = DeallocateVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to deallocte virtual machines in a VM scale set: %+v", err)
	}
	t.Logf("deallocted virtual machines in a VM scale set")

	err = DeleteVirtualMachineScaleSet(ctx, virtualMachineScaleSetName)
	if err != nil {
		t.Fatalf("failed to delete virtual machine scale set: %+v", err)
	}
	t.Logf("deleted virtual machine scale set")
}
