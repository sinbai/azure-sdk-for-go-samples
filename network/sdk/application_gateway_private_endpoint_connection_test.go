// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

package network

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/Azure-Samples/azure-sdk-for-go-samples/internal/config"
	"github.com/Azure-Samples/azure-sdk-for-go-samples/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/arm/network/2020-07-01/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

func TestApplicationGatewayPrivateEndpointConnection(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)
	connectionName := config.AppendRandomSuffix("agpeconnection")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetAppgwName := config.AppendRandomSuffix("subnetappgw")
	subnetLoadBalancerName := config.AppendRandomSuffix("subnetloadbalancer")
	subnetEndpointName := config.AppendRandomSuffix("subnetendpoint")
	privateEndpointName := config.AppendRandomSuffix("privateendpoint")
	privateLinkServiceName := config.AppendRandomSuffix("privatelinkservice")
	loadBalancerName := config.AppendRandomSuffix("loadbalancer")
	ipConfigurationName := config.AppendRandomSuffix("ipconfig")
	applicationGatewayName := config.AppendRandomSuffix("applicationgateway")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	frontendPortName := config.AppendRandomSuffix("appgwfp")
	frontendPortName2 := config.AppendRandomSuffix("appgwfp80")
	sslCertificateName1 := config.AppendRandomSuffix("sslcert")
	gatewayIpConfiguration := config.AppendRandomSuffix("appgwipc")
	sslProfileName := config.AppendRandomSuffix("sslprofile")
	httpListenerName1 := config.AppendRandomSuffix("httplistener")
	httpListenerName2 := config.AppendRandomSuffix("httplistener")
	backendHttpSettingsCollectionName := config.AppendRandomSuffix("backendhttpsettingscollection")
	rewriteRuleSetName := config.AppendRandomSuffix("rewriteruleset")
	urlPathMapName := config.AppendRandomSuffix("pathmap")
	frontendIpConfigurationName := config.AppendRandomSuffix("frontendipconfiguration")
	backendAddressPoolName := config.AppendRandomSuffix("appgwpool")
	networkSecurityGroupName := config.AppendRandomSuffix("networksecuritygroup")

	ctx, cancel := context.WithTimeout(context.Background(), 5000*time.Second)
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
				AddressPrefixes: &[]*string{to.StringPtr("10.0.0.0/16")},
			},
		},
	}
	_, err = CreateVirtualNetwork(ctx, virtualNetworkName, virtualNetworkParameters)
	if err != nil {
		t.Fatalf("failed to create virtual network: % +v", err)
	}

	networkSecurityGroupParameters := armnetwork.NetworkSecurityGroup{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.NetworkSecurityGroupPropertiesFormat{
			SecurityRules: &[]*armnetwork.SecurityRule{{
				Name: to.StringPtr("rule1"),
				Properties: &armnetwork.SecurityRulePropertiesFormat{
					Access:                   armnetwork.SecurityRuleAccessAllow.ToPtr(),
					DestinationAddressPrefix: to.StringPtr("*"),
					DestinationPortRange:     to.StringPtr("65200-65535"),
					Direction:                armnetwork.SecurityRuleDirectionInbound.ToPtr(),
					Priority:                 to.Int32Ptr(100),
					Protocol:                 armnetwork.SecurityRuleProtocolAsterisk.ToPtr(),
					SourceAddressPrefix:      to.StringPtr("*"),
					SourcePortRange:          to.StringPtr("*"),
				},
			}},
		},
	}
	secruityGroupId, err := CreateNetworkSecurityGroup(ctx, networkSecurityGroupName, networkSecurityGroupParameters)
	if err != nil {
		t.Fatalf("failed to create network security group: % +v", err)
	}
	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/24"),
			NetworkSecurityGroup: &armnetwork.NetworkSecurityGroup{
				Resource: armnetwork.Resource{
					ID: &secruityGroupId,
				},
			},
		},
	}
	subnetAppgwId, err := CreateSubnet(ctx, virtualNetworkName, subnetAppgwName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                     to.StringPtr("10.0.1.0/24"),
			PrivateLinkServiceNetworkPolicies: to.StringPtr("Disabled"),
		},
	}
	subnetLoadBalancerId, err := CreateSubnet(ctx, virtualNetworkName, subnetLoadBalancerName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	subnetParameters = armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix:                  to.StringPtr("10.0.2.0/24"),
			PrivateEndpointNetworkPolicies: to.StringPtr("Disabled"),
		},
	}
	subnetEndpointId, err := CreateSubnet(ctx, virtualNetworkName, subnetEndpointName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	loadBalancerParameters := armnetwork.LoadBalancer{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.LoadBalancerPropertiesFormat{
			FrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{
				{
					Name: &ipConfigurationName,
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						Subnet: &armnetwork.Subnet{
							SubResource: armnetwork.SubResource{
								ID: &subnetLoadBalancerId,
							},
						},
					},
				},
			},
		},
		SKU: &armnetwork.LoadBalancerSKU{
			Name: armnetwork.LoadBalancerSKUNameStandard.ToPtr(),
		},
	}

	loadBalancerId, err := CreateLoadBalancer(ctx, loadBalancerName, loadBalancerParameters)
	if err != nil {
		t.Fatalf("failed to create load balancer: % +v", err)
	}

	privateLinkServiceParameters := armnetwork.PrivateLinkService{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PrivateLinkServiceProperties{
			AutoApproval: &armnetwork.PrivateLinkServicePropertiesAutoApproval{
				ResourceSet: armnetwork.ResourceSet{
					Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
				},
			},
			Fqdns: &[]*string{to.StringPtr("fqdn1"),
				to.StringPtr("fqdn2"),
				to.StringPtr("fqdn3")},
			IPConfigurations: &[]*armnetwork.PrivateLinkServiceIPConfiguration{{
				Name: &ipConfigurationName,
				Properties: &armnetwork.PrivateLinkServiceIPConfigurationProperties{
					PrivateIPAddress:          to.StringPtr("10.0.1.5"),
					PrivateIPAddressVersion:   armnetwork.IPVersionIPv4.ToPtr(),
					PrivateIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
					Subnet: &armnetwork.Subnet{
						SubResource: armnetwork.SubResource{
							ID: &subnetLoadBalancerId,
						},
					},
				},
			}},
			LoadBalancerFrontendIPConfigurations: &[]*armnetwork.FrontendIPConfiguration{{
				SubResource: armnetwork.SubResource{
					ID: to.StringPtr(loadBalancerId + "/frontendIPConfigurations/" + ipConfigurationName),
				},
			}},
			Visibility: &armnetwork.PrivateLinkServicePropertiesVisibility{
				ResourceSet: armnetwork.ResourceSet{
					Subscriptions: &[]*string{to.StringPtr(config.SubscriptionID())},
				},
			},
		},
	}
	privateLinkServiceId, err := CreatePrivateLinkService(ctx, privateLinkServiceName, privateLinkServiceParameters)
	if err != nil {
		t.Fatalf("failed to create private link service: % +v", err)
	}

	privateEndpointParameters := armnetwork.PrivateEndpoint{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.PrivateEndpointProperties{
			PrivateLinkServiceConnections: &[]*armnetwork.PrivateLinkServiceConnection{{
				Name: &privateLinkServiceName,
				Properties: &armnetwork.PrivateLinkServiceConnectionProperties{
					PrivateLinkServiceID: &privateLinkServiceId,
				},
			}},
			Subnet: &armnetwork.Subnet{
				SubResource: armnetwork.SubResource{
					ID: &subnetEndpointId,
				},
			},
		},
	}

	privateEndpointId, err := CreatePrivateEndpoint(ctx, privateEndpointName, privateEndpointParameters)
	if err != nil {
		t.Fatalf("failed to create private endpoint: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
		SKU: &armnetwork.PublicIPAddressSKU{Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr()},
	}
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters) //
	if err != nil {
		t.Fatalf("failed to create public ip address: %+v", err)
	}

	certPfx, err := os.ReadFile("../testdata/application_gateway_test.pfx")
	if err != nil {
		t.Fatal(err)
	}
	certB64 := base64.StdEncoding.EncodeToString(certPfx)

	applicationGatewayUrl := "/subscriptions/" + config.SubscriptionID() + "/resourceGroups/" + config.GroupName() + "/providers/Microsoft.Network/applicationGateways/" + applicationGatewayName

	applicationGatewayParameters := armnetwork.ApplicationGateway{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(config.Location()),
		},
		Properties: &armnetwork.ApplicationGatewayPropertiesFormat{
			SKU: &armnetwork.ApplicationGatewaySKU{
				Capacity: to.Int32Ptr(3),
				Name:     armnetwork.ApplicationGatewaySKUNameStandardV2.ToPtr(),
				Tier:     armnetwork.ApplicationGatewayTierStandardV2.ToPtr(),
			},
			GatewayIPConfigurations: &[]*armnetwork.ApplicationGatewayIPConfiguration{{
				Name: &gatewayIpConfiguration,
				Properties: &armnetwork.ApplicationGatewayIPConfigurationPropertiesFormat{
					Subnet: &armnetwork.SubResource{
						ID: &subnetAppgwId,
					},
				},
			}},
			SSLCertificates: &[]*armnetwork.ApplicationGatewaySSLCertificate{
				{
					Name: &sslCertificateName1,
					Properties: &armnetwork.ApplicationGatewaySSLCertificatePropertiesFormat{
						Data:     &certB64,
						Password: to.StringPtr("123456"),
					},
				},
			},

			FrontendIPConfigurations: &[]*armnetwork.ApplicationGatewayFrontendIPConfiguration{{
				Name: &frontendIpConfigurationName,
				Properties: &armnetwork.ApplicationGatewayFrontendIPConfigurationPropertiesFormat{
					PublicIPAddress: &armnetwork.SubResource{
						ID: &publicIpAddressId,
					},
				},
			}},
			FrontendPorts: &[]*armnetwork.ApplicationGatewayFrontendPort{{
				Name: &frontendPortName,
				Properties: &armnetwork.ApplicationGatewayFrontendPortPropertiesFormat{
					Port: to.Int32Ptr(443),
				},
			}, {
				Name: &frontendPortName2,
				Properties: &armnetwork.ApplicationGatewayFrontendPortPropertiesFormat{
					Port: to.Int32Ptr(80),
				},
			}},

			BackendAddressPools: &[]*armnetwork.ApplicationGatewayBackendAddressPool{
				{
					Name: &backendAddressPoolName,
					Properties: &armnetwork.ApplicationGatewayBackendAddressPoolPropertiesFormat{
						BackendAddresses: &[]*armnetwork.ApplicationGatewayBackendAddress{
							{
								IPAddress: to.StringPtr("10.0.1.1"),
							},
							{
								IPAddress: to.StringPtr("10.0.1.2"),
							},
						},
					},
				},
			},
			BackendHTTPSettingsCollection: &[]*armnetwork.ApplicationGatewayBackendHTTPSettings{{
				Name: &backendHttpSettingsCollectionName,
				Properties: &armnetwork.ApplicationGatewayBackendHTTPSettingsPropertiesFormat{
					CookieBasedAffinity: armnetwork.ApplicationGatewayCookieBasedAffinityDisabled.ToPtr(),
					Port:                to.Int32Ptr(80),
					Protocol:            armnetwork.ApplicationGatewayProtocolHTTP.ToPtr(),
					RequestTimeout:      to.Int32Ptr(30),
				},
			}},
			SSLProfiles: &[]*armnetwork.ApplicationGatewaySSLProfile{{
				Name: &sslProfileName,
				Properties: &armnetwork.ApplicationGatewaySSLProfilePropertiesFormat{
					ClientAuthConfiguration: &armnetwork.ApplicationGatewayClientAuthConfiguration{
						VerifyClientCertIssuerDN: to.BoolPtr(true),
					},
					SSLPolicy: &armnetwork.ApplicationGatewaySSLPolicy{
						PolicyName: armnetwork.ApplicationGatewaySSLPolicyNameAppGwSSLPolicy20170401.ToPtr(),
						PolicyType: armnetwork.ApplicationGatewaySSLPolicyTypePredefined.ToPtr(),
					},
				},
			}},
			HTTPListeners: &[]*armnetwork.ApplicationGatewayHTTPListener{{
				Name: &httpListenerName1,
				Properties: &armnetwork.ApplicationGatewayHTTPListenerPropertiesFormat{
					FrontendIPConfiguration: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
					},
					FrontendPort: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/frontendPorts/" + frontendPortName),
					},
					Protocol:                    armnetwork.ApplicationGatewayProtocolHTTPS.ToPtr(),
					RequireServerNameIndication: to.BoolPtr(false),
					SSLCertificate: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/sslCertificates/" + sslCertificateName1),
					},
					SSLProfile: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/sslProfiles/" + sslProfileName),
					},
				},
			}, {
				Name: &httpListenerName2,
				Properties: &armnetwork.ApplicationGatewayHTTPListenerPropertiesFormat{
					FrontendIPConfiguration: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/frontendIPConfigurations/" + frontendIpConfigurationName),
					},
					FrontendPort: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/frontendPorts/" + frontendPortName2),
					},
					Protocol: armnetwork.ApplicationGatewayProtocolHTTP.ToPtr(),
				},
			}},
			URLPathMaps: &[]*armnetwork.ApplicationGatewayURLPathMap{{
				Name: &urlPathMapName,
				Properties: &armnetwork.ApplicationGatewayURLPathMapPropertiesFormat{
					DefaultBackendAddressPool: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/backendAddressPools/" + backendAddressPoolName),
					},
					DefaultBackendHTTPSettings: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/backendHttpSettingsCollection/" + backendHttpSettingsCollectionName),
					},
					DefaultRewriteRuleSet: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/rewriteRuleSets/" + rewriteRuleSetName),
					},
					PathRules: &[]*armnetwork.ApplicationGatewayPathRule{{
						Name: to.StringPtr("apiPaths"),
						Properties: &armnetwork.ApplicationGatewayPathRulePropertiesFormat{
							BackendAddressPool: &armnetwork.SubResource{
								ID: to.StringPtr(applicationGatewayUrl + "/backendAddressPools/" + backendAddressPoolName),
							},
							BackendHTTPSettings: &armnetwork.SubResource{
								ID: to.StringPtr(applicationGatewayUrl + "/backendHttpSettingsCollection/" + backendHttpSettingsCollectionName),
							},
							Paths: &[]*string{to.StringPtr("/api"), to.StringPtr("/v1/api")},
							RewriteRuleSet: &armnetwork.SubResource{
								ID: to.StringPtr(applicationGatewayUrl + "/rewriteRuleSets/" + rewriteRuleSetName),
							},
						},
					}},
				},
			}},
			RequestRoutingRules: &[]*armnetwork.ApplicationGatewayRequestRoutingRule{{
				Name: to.StringPtr("appgwrule"),
				Properties: &armnetwork.ApplicationGatewayRequestRoutingRulePropertiesFormat{
					BackendAddressPool: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/backendAddressPools/" + backendAddressPoolName),
					},
					BackendHTTPSettings: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/backendHttpSettingsCollection/" + backendHttpSettingsCollectionName),
					},
					HTTPListener: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/httpListeners/" + httpListenerName1),
					},
					Priority: to.Int32Ptr(10),
					RewriteRuleSet: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/rewriteRuleSets/" + rewriteRuleSetName),
					},
					RuleType: armnetwork.ApplicationGatewayRequestRoutingRuleTypeBasic.ToPtr(),
				},
			}, {
				Name: to.StringPtr("appgwPathBasedRule"),
				Properties: &armnetwork.ApplicationGatewayRequestRoutingRulePropertiesFormat{
					HTTPListener: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/httpListeners/" + httpListenerName2),
					},
					Priority: to.Int32Ptr(20),
					RuleType: armnetwork.ApplicationGatewayRequestRoutingRuleTypePathBasedRouting.ToPtr(),
					URLPathMap: &armnetwork.SubResource{
						ID: to.StringPtr(applicationGatewayUrl + "/urlPathMaps/" + urlPathMapName),
					},
				},
			}},
			RewriteRuleSets: &[]*armnetwork.ApplicationGatewayRewriteRuleSet{{
				Name: &rewriteRuleSetName,
				Properties: &armnetwork.ApplicationGatewayRewriteRuleSetPropertiesFormat{
					RewriteRules: &[]*armnetwork.ApplicationGatewayRewriteRule{{
						ActionSet: &armnetwork.ApplicationGatewayRewriteRuleActionSet{
							RequestHeaderConfigurations: &[]*armnetwork.ApplicationGatewayHeaderConfiguration{{
								HeaderName:  to.StringPtr("X-Forwarded-For"),
								HeaderValue: to.StringPtr("{var_add_x_forwarded_for_proxy}"),
							}},
							ResponseHeaderConfigurations: &[]*armnetwork.ApplicationGatewayHeaderConfiguration{{
								HeaderName:  to.StringPtr("Strict-Transport-Security"),
								HeaderValue: to.StringPtr("max-age=31536000"),
							}},
							URLConfiguration: &armnetwork.ApplicationGatewayURLConfiguration{
								ModifiedPath: to.StringPtr("/abc"),
							},
						},
						Conditions: &[]*armnetwork.ApplicationGatewayRewriteRuleCondition{{
							IgnoreCase: to.BoolPtr(true),
							Negate:     to.BoolPtr(false),
							Pattern:    to.StringPtr("^Bearer"),
							Variable:   to.StringPtr("http_req_Authorization"),
						}},
						Name:         to.StringPtr("Set X-Forwarded-For"),
						RuleSequence: to.Int32Ptr(102),
					}},
				},
			}},
		},
	}
	err = CreateApplicationGateway(ctx, applicationGatewayName, applicationGatewayParameters)
	if err != nil {
		t.Fatalf("failed to create application gateway: % +v", err)
	}

	applicationGatewayPrivateEndpointConnectionParameters := armnetwork.ApplicationGatewayPrivateEndpointConnection{
		Name: &connectionName,
		Properties: &armnetwork.ApplicationGatewayPrivateEndpointConnectionProperties{
			PrivateEndpoint: &armnetwork.PrivateEndpoint{
				Resource: armnetwork.Resource{
					ID: &privateEndpointId,
				},
			},
			PrivateLinkServiceConnectionState: &armnetwork.PrivateLinkServiceConnectionState{
				Description: to.StringPtr("approved it for some reason."),
				Status:      to.StringPtr("Approved"),
			},
		},
	}
	err = UpdateApplicationGatewayPrivateEndpointConnection(ctx, applicationGatewayName, connectionName, applicationGatewayPrivateEndpointConnectionParameters)
	if err != nil {
		t.Fatalf("failed to update application gateway private endpoint connection: % +v", err)
	}
	t.Logf("updated application gateway private endpoint connection")

	err = GetApplicationGatewayPrivateEndpointConnection(ctx, applicationGatewayName, connectionName)
	if err != nil {
		t.Fatalf("failed to get application gateway private endpoint connection: %+v", err)
	}
	t.Logf("got application gateway private endpoint connection")

	err = ListApplicationGatewayPrivateEndpointConnection(ctx, connectionName)
	if err != nil {
		t.Fatalf("failed to list application gateway private endpoint connection: %+v", err)
	}
	t.Logf("listed application gateway private endpoint connection")

	err = DeleteApplicationGatewayPrivateEndpointConnection(ctx, applicationGatewayName, connectionName)
	if err != nil {
		t.Fatalf("failed to delete application gateway private endpoint connection: %+v", err)
	}
	t.Logf("deleted application gateway private endpoint connection")
}
