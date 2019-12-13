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

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils/provider_wrapper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type kafkaConfig struct {
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
type KafkaProvider struct {
	terraform_utils.Provider
	config kafkaConfig
}

func (p *KafkaProvider) Init(args []string) error {
	p.config.bootstrapServers = args[0] //strings.Split(args[0], ",")
	p.config.caCert = args[1]
	p.config.clientCert = args[2]
	p.config.clientCertKey = args[3]
	p.config.saslUsername = args[4]
	p.config.saslPassword = args[5]
	p.config.saslMechanism = args[6]
	p.config.skipTLSVerify, _ = strconv.ParseBool(args[7])
	p.config.tlsEnabled, _ = strconv.ParseBool(args[8])
	p.config.timeout, _ = strconv.Atoi(args[9])
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

	bootstrap, err := gocty.ToCtyValue(strings.Split(p.config.bootstrapServers, ","), cty.List(cty.String))
	if err != nil {
		panic(err)
	}
	config := cty.ObjectVal(map[string]cty.Value{
		"bootstrap_servers": bootstrap,
		"ca_cert":           cty.StringVal(p.config.caCert),
		"client_cert":       cty.StringVal(p.config.clientCert),
		"client_key":        cty.StringVal(p.config.clientCertKey),
		"sasl_username":     cty.StringVal(p.config.saslUsername),
		"sasl_password":     cty.StringVal(p.config.saslPassword),
		"sasl_mechanism":    cty.StringVal(p.config.saslMechanism),
		"skip_tls_verify":   cty.BoolVal(p.config.skipTLSVerify),
		"tls_enabled":       cty.BoolVal(p.config.tlsEnabled),
		"timeout":           cty.NumberIntVal(int64(p.config.timeout)),
	})
	return config
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
		"bootstrap_servers": p.config.bootstrapServers,
		"ca_cert":           p.config.caCert,
		"client_cert":       p.config.clientCert,
		"client_key":        p.config.clientCertKey,
		"sasl_username":     p.config.saslUsername,
		"sasl_password":     p.config.saslPassword,
		"sasl_mechanism":    p.config.saslMechanism,
		"skip_tls_verify":   p.config.skipTLSVerify,
		"tls_enabled":       p.config.tlsEnabled,
		"timeout":           p.config.timeout,
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
