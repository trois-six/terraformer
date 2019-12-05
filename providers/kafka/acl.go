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

package kafka

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type ACLGenerator struct {
	KafkaService
}

type ACL struct {
	Principal      string
	Host           string
	Operation      string
	PermissionType string
}

type ACLs []ACL

var ACLAllowEmptyValues = []string{}

func (g ACLGenerator) createResources(acls ACLs) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, acl := range acls {
		resources = append(resources, terraform_utils.NewSimpleResource(
			acl.Principal, // TODO
			fmt.Sprintf("acl_%s", normalizeResourceName(acl.Principal)), // TOOD
			"kafka_acl",
			"kafka",
			ACLAllowEmptyValues,
		))
	}
	return resources
}

func (g *ACLGenerator) InitResources() error {
	var acls ACLs
	/* TODO
	 */
	g.Resources = g.createResources(acls)
	return nil
}
