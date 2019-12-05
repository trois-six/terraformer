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
	"errors"
	"strconv"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils/provider_wrapper"
	"github.com/zclconf/go-cty/cty"
)

type KafkaProvider struct {
	terraform_utils.Provider
	bootstrapServers string
	caCert           string
	clientCert       string
	clientCertKey    string
	saslUsername     string
	saslPassword     string
	saslMechanism    string
	skipTLSVerify    bool
	tlsEnabled       bool
	timeout          int
}

func (p *KafkaProvider) Init(args []string) error {
	p.bootstrapServers = args[0]
	p.caCert = args[1]
	p.clientCert = args[2]
	p.clientCertKey = args[3]
	p.saslUsername = args[4]
	p.saslPassword = args[5]
	p.saslMechanism = args[6]
	p.skipTLSVerify, _ = strconv.ParseBool(args[7])
	p.tlsEnabled, _ = strconv.ParseBool(args[8])
	p.timeout, _ = strconv.Atoi(args[9])
	return nil
}

func (p *KafkaProvider) GetName() string {
	return "kafka"
}

func (p *KafkaProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			p.GetName(): map[string]interface{}{
				"version": provider_wrapper.GetProviderVersion(p.GetName()),
			},
		},
	}
}

func (p *KafkaProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"bootstrap_servers": cty.StringVal(p.bootstrapServers),
		"ca_cert":           cty.StringVal(p.caCert),
		"client_cert":       cty.StringVal(p.clientCert),
		"client_key":        cty.StringVal(p.clientCertKey),
		"sasl_username":     cty.StringVal(p.saslUsername),
		"sasl_password":     cty.StringVal(p.saslPassword),
		"sasl_mechanism":    cty.StringVal(p.saslMechanism),
		"skip_tls_verify":   cty.BoolVal(p.skipTLSVerify),
		"tls_enabled":       cty.BoolVal(p.tlsEnabled),
		"timeout":           cty.NumberIntVal(int64(p.timeout)),
	})
}

func (p *KafkaProvider) GetBasicConfig() cty.Value {
	return p.GetConfig()
}

func (p *KafkaProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"bootstrap_servers": p.bootstrapServers,
		"ca_cert":           p.caCert,
		"client_cert":       p.clientCert,
		"client_key":        p.clientCertKey,
		"sasl_username":     p.saslUsername,
		"sasl_password":     p.saslPassword,
		"sasl_mechanism":    p.saslMechanism,
		"skip_tls_verify":   p.skipTLSVerify,
		"tls_enabled":       p.tlsEnabled,
		"timeout":           p.timeout,
	})
	return nil
}

func (p *KafkaProvider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"topics": &TopicGenerator{},
		"acls":   &ACLGenerator{},
	}
}

func (KafkaProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}
