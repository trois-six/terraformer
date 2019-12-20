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
	_ "os"

	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type TenantGenerator struct {
	CiscoACIService
}

var TenantAllowEmptyValues = []string{}

func (g TenantGenerator) createResources(tenants *[]*models.Tenant) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, tenant := range *tenants {
		resources = append(resources, terraform_utils.NewSimpleResource(
			tenant.DistinguishedName,
			"tenant_"+normalizeResourceName(tenant.DistinguishedName[len("uni/tn-"):]),
			"aci_tenant",
			"aci",
			TenantAllowEmptyValues,
		))
	}
	return resources
}

func (g *TenantGenerator) InitResources() error {
	/*
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	os.Stdout = null
	*/

	aciClient := g.getClient(g.Args["insecure"].(bool))
	path := "api/node/class/fvTenant.json?query-target=self&rsp-subtree=no&rsp-prop-include=config-only"
	fvTenantCont, err := aciClient.GetViaURL(path)
	if err != nil {
		return err
	}

	tenants := models.TenantListFromContainer(fvTenantCont)
	//os.Stdout = sout

	g.Resources = g.createResources(&tenants)

	return nil
}
