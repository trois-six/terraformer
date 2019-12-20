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

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/ciscoecosystem/aci-go-client/client"
)

type CiscoACIService struct {
	terraform_utils.Service
}

func (s *CiscoACIService) getClient(isInsecure bool) *client.Client {
	if os.Getenv("ACI_PASSWORD") != "" {
		return client.GetClient(os.Getenv("ACI_URL"), os.Getenv("ACI_USERNAME"), client.Password(os.Getenv("ACI_PASSWORD")), client.Insecure(isInsecure), client.ProxyUrl(os.Getenv("ACI_PROXY_URL")))
	} else {
		return client.GetClient(os.Getenv("ACI_URL"), os.Getenv("ACI_USERNAME"), client.PrivateKey(os.Getenv("ACI_PRIVATE_KEY")), client.AdminCert(os.Getenv("ACI_CERT_NAME")), client.Insecure(isInsecure), client.ProxyUrl(os.Getenv("ACI_PROXY_URL")))
	}
}
