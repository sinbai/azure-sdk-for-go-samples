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

func TestApplicationGateway(t *testing.T) {
	groupName := config.GenerateGroupName("network")
	config.SetGroupName(groupName)

	applicationGatewayName := config.AppendRandomSuffix("applicationgateway")
	publicIpAddressName := config.AppendRandomSuffix("pipaddress")
	virtualNetworkName := config.AppendRandomSuffix("virtualnetwork")
	subnetName := config.AppendRandomSuffix("subnet")
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

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
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

	subnetParameters := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.0.0/16"),
		},
	}
	subNetId, err := CreateSubnet(ctx, virtualNetworkName, subnetName, subnetParameters)
	if err != nil {
		t.Fatalf("failed to create sub net: % +v", err)
	}

	publicIPAddressParameters := armnetwork.PublicIPAddress{
		Resource: armnetwork.Resource{Location: to.StringPtr(config.Location())},
		Properties: &armnetwork.PublicIPAddressPropertiesFormat{
			PublicIPAllocationMethod: armnetwork.IPAllocationMethodStatic.ToPtr(),
		},
		SKU: &armnetwork.PublicIPAddressSKU{Name: armnetwork.PublicIPAddressSKUNameStandard.ToPtr()},
	}
	publicIpAddressId, err := CreatePublicIPAddress(ctx, publicIpAddressName, publicIPAddressParameters)
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
						ID: &subNetId,
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
	t.Logf("created application gateway")

	err = GetApplicationGatewaySSLPredefinedPolicy(ctx, string(armnetwork.ApplicationGatewaySSLPolicyNameAppGwSSLPolicy20170401))
	if err != nil {
		t.Fatalf("failed to get application gateway ssl predefined policy: %+v", err)
	}
	t.Logf("got application gateway ssl predefined policy")

	err = ListApplicationGatewayAvailableSSLPredefinedPolicie(ctx)
	if err != nil {
		t.Fatalf("failed to list all SSL predefined policies for configuring Ssl policy: %+v", err)
	}
	t.Logf("listed all SSL predefined policies for configuring Ssl policy")

	err = ListApplicationGatewayAvailableSSLOptions(ctx)
	if err != nil {
		t.Fatalf("failed to list available Ssl options for configuring Ssl policy: %+v", err)
	}
	t.Logf("listed available Ssl options for configuring Ssl policy")

	err = GetApplicationGateway(ctx, applicationGatewayName)
	if err != nil {
		t.Fatalf("failed to get application gateway: %+v", err)
	}
	t.Logf("got application gateway")

	err = ListApplicationGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list application gateway: %+v", err)
	}
	t.Logf("listed application gateway")

	err = ListApplicationGatewayAvailableServerVariables(ctx)
	if err != nil {
		t.Fatalf("failed to list all available server variables: %+v", err)
	}
	t.Logf("listed all available server variables")

	err = ListApplicationGatewayAvailableResponseHeaders(ctx)
	if err != nil {
		t.Fatalf("failed to list all available response headers: %+v", err)
	}
	t.Logf("listed all available response headers")

	err = ListApplicationGatewayAvailableWafRuleSets(ctx)
	if err != nil {
		t.Fatalf("failed to list all available web application firewall rule sets: %+v", err)
	}
	t.Logf("listed all available web application firewall rule sets")

	err = ListAllApplicationGateway(ctx)
	if err != nil {
		t.Fatalf("failed to list all application gateway: %+v", err)
	}
	t.Logf("listed all application gateway")

	probeRequestParameters := armnetwork.ApplicationGatewayOnDemandProbe{
		BackendAddressPool: &armnetwork.SubResource{
			ID: to.StringPtr(applicationGatewayUrl + "/backendaddressPools/" + backendAddressPoolName),
		},
		BackendHTTPSettings: &armnetwork.SubResource{
			ID: to.StringPtr(applicationGatewayUrl + "/backendHttpSettingsCollection/" + backendHttpSettingsCollectionName),
		},
		Path:                                to.StringPtr("/"),
		PickHostNameFromBackendHTTPSettings: to.BoolPtr(true),
		Protocol:                            armnetwork.ApplicationGatewayProtocolHTTP.ToPtr(),
		Timeout:                             to.Int32Ptr(30),
	}
	err = GetApplicationGatewayBackendHealthOnDemand(ctx, applicationGatewayName, probeRequestParameters)
	if err != nil {
		t.Fatalf("failed to get the backend health for given combination of backend pool and http setting of the specified application gateway: %+v", err)
	}
	t.Logf("got the backend health for given combination of backend pool and http setting of the specified application gateway")

	err = GetApplicationGatewayBackendHealth(ctx, applicationGatewayName)
	if err != nil {
		t.Fatalf("failed to get the backend health of the specified application gateway: %+v", err)
	}
	t.Logf("got the backend health of the specified application gateway")

	err = StartApplicationGateway(ctx, applicationGatewayName)
	if err != nil {
		t.Fatalf("failed to start the specified application gateway: %+v", err)
	}
	t.Logf("started the specified application gateway")

	err = StopApplicationGateway(ctx, applicationGatewayName)
	if err != nil {
		t.Fatalf("failed to stop the specified application gateway in a resource group: %+v", err)
	}
	t.Logf("stopped the specified application gateway in a resource group")

	tagsObjectParameters := armnetwork.TagsObject{
		Tags: &map[string]*string{"tag1": to.StringPtr("value1"), "tag2": to.StringPtr("value2")},
	}
	err = UpdateApplicationGatewayTags(ctx, applicationGatewayName, tagsObjectParameters)
	if err != nil {
		t.Fatalf("failed to update tags for application gateway: %+v", err)
	}
	t.Logf("updated application gateway tags")

	err = DeleteApplicationGateway(ctx, applicationGatewayName)
	if err != nil {
		t.Fatalf("failed to delete application gateway: %+v", err)
	}
	t.Logf("deleted application gateway")
}
