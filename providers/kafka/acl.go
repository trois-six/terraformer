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
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type ACLGenerator struct {
	KafkaService
}

type ACL struct {
	ResourceName      string `json:"resource_name"`
	ResourceType      string `json:"resource_type"`
	ACLPrincipal      string `json:"acl_principal"`
	ACLHost           string `json:"acl_host"`
	ACLOperation      string `json:"acl_operation"`
	ACLPermissionType string `json:"acl_permission_type"`
}

type ACLs []ACL

var ACLAllowEmptyValues = []string{}
var ACLAdditionalFields = map[string]interface{}{}

func (a ACLGenerator) createResources(acls ACLs) []terraformutils.Resource {
	var resources []terraformutils.Resource
	return resources
}

func (a *ACLGenerator) InitResources() error {
	var acls ACLs
	a.Resources = a.createResources(acls)
	return nil
}
