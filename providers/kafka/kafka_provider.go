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
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/zclconf/go-cty/cty"
)

type KafkaProvider struct { //nolint
	terraformutils.Provider
	bootstrapServers    []string
	caCert              string
	clientCert          string
	clientKey           string
	clientKeyPassphrase string
	saslUsername        string
	saslPassword        string
	saslMechanism       string
	skipTLSVerify       bool
	tlsEnabled          bool
}

func (p *KafkaProvider) Init(args []string) error {
	p.bootstrapServers = strings.Split(args[0], ",")
	p.caCert = args[1]
	p.clientCert = args[2]
	p.clientKey = args[3]
	p.clientKeyPassphrase = args[4]
	p.saslUsername = args[5]
	p.saslPassword = args[6]
	p.saslMechanism = args[7]
	p.skipTLSVerify, _ = strconv.ParseBool(args[8])
	p.tlsEnabled, _ = strconv.ParseBool(args[9])

	return nil
}

func (p *KafkaProvider) GetName() string {
	return "kafka"
}

func (p *KafkaProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{}
}

func (p *KafkaProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"bootstrap_servers":     cty.StringVal(strings.Join(p.bootstrapServers, ",")),
		"ca_cert":               cty.StringVal(p.caCert),
		"client_cert":           cty.StringVal(p.clientCert),
		"client_key":            cty.StringVal(p.clientKey),
		"client_key_passphrase": cty.StringVal(p.clientKeyPassphrase),
		"sasl_username":         cty.StringVal(p.saslUsername),
		"sasl_password":         cty.StringVal(p.saslPassword),
		"sasl_mechanism":        cty.StringVal(p.saslMechanism),
		"skip_tls_verify":       cty.BoolVal(p.skipTLSVerify),
		"tls_enabled":           cty.BoolVal(p.tlsEnabled),
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
		"bootstrap_servers":     p.bootstrapServers,
		"ca_cert":               p.caCert,
		"client_cert":           p.clientCert,
		"client_key":            p.clientKey,
		"client_key_passphrase": p.clientKeyPassphrase,
		"sasl_username":         p.saslUsername,
		"sasl_password":         p.saslPassword,
		"sasl_mechanism":        p.saslMechanism,
		"skip_tls_verify":       p.skipTLSVerify,
		"tls_enabled":           p.tlsEnabled,
	})

	return nil
}

func (p *KafkaProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {

	return map[string]terraformutils.ServiceGenerator{
		"acls":   &ACLGenerator{},
		"quotas": &QuotaGenerator{},
		"topics": &TopicGenerator{},
	}
}

func (KafkaProvider) GetResourceConnections() map[string]map[string][]string {

	return map[string]map[string][]string{}
}
