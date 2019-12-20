// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aci

import (
	"os"
	"strconv"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils/provider_wrapper"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type CiscoACIProvider struct {
	terraform_utils.Provider
	isInsecure bool
}

func (p *CiscoACIProvider) Init(args []string) error {
	if os.Getenv("ACI_USERNAME") == "" {
		return errors.New("aci: Username must be provided for the ACI provider")
	}

	if os.Getenv("ACI_PASSWORD") == "" {
		if os.Getenv("ACI_PRIVATE_KEY") == "" && os.Getenv("ACI_CERT_NAME") == "" {
			return errors.New("aci: Either of private_key/cert_name or password is required")
		} else if os.Getenv("ACI_PRIVATE_KEY") == "" || os.Getenv("ACI_CERT_NAME") == "" {
			return errors.New("aci: private_key and cert_name both must be provided")
		}
	}

	if os.Getenv("ACI_URL") == "" {
		return errors.New("aci: The URL must be provided for the ACI provider")
	}

	p.isInsecure, _ = strconv.ParseBool(os.Getenv("ACI_INSECURE"))

	return nil
}

func (p *CiscoACIProvider) GetName() string {
	return "aci"
}

func (p *CiscoACIProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			p.GetName(): map[string]interface{}{
				"version": provider_wrapper.GetProviderVersion(p.GetName()),
			},
		},
	}
}

func (p *CiscoACIProvider) GetBasicConfig() cty.Value {
	return p.GetConfig()
}

func (p *CiscoACIProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"insecure": p.isInsecure,
	})
	return nil
}

func (p *CiscoACIProvider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"tenants": &TenantGenerator{},
	}
}

func (CiscoACIProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}
