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
	"log"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils/provider_wrapper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type kafkaConfig struct {
	BootstrapServers string
	CACert           string
	ClientCert       string
	ClientCertKey    string
	SASLUsername     string
	SASLPassword     string
	SASLMechanism    string
	SkipTLSVerify    bool
	TLSEnabled       bool
	Timeout          int
}
type KafkaProvider struct {
	terraform_utils.Provider
	config kafkaConfig
}

func (p *KafkaProvider) Init(args []string) error {
	p.config.BootstrapServers = args[0] //strings.Split(args[0], ",")
	p.config.CACert = args[1]
	p.config.ClientCert = args[2]
	p.config.ClientCertKey = args[3]
	p.config.SASLUsername = args[4]
	p.config.SASLPassword = args[5]
	p.config.SASLMechanism = args[6]
	p.config.SkipTLSVerify, _ = strconv.ParseBool(args[7])
	p.config.TLSEnabled, _ = strconv.ParseBool(args[8])
	p.config.Timeout, _ = strconv.Atoi(args[9])
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

	bootstrap, err := gocty.ToCtyValue(strings.Split(p.config.BootstrapServers, ","), cty.List(cty.String))
	if err != nil {
		log.Fatal("Cannot convert BootstrapServers to slice: ", err.Error())
	}
	config := cty.ObjectVal(map[string]cty.Value{
		"bootstrap_servers": bootstrap,
		"ca_cert":           cty.StringVal(p.config.CACert),
		"client_cert":       cty.StringVal(p.config.ClientCert),
		"client_key":        cty.StringVal(p.config.ClientCertKey),
		"sasl_username":     cty.StringVal(p.config.SASLUsername),
		"sasl_password":     cty.StringVal(p.config.SASLPassword),
		"sasl_mechanism":    cty.StringVal(p.config.SASLMechanism),
		"skip_tls_verify":   cty.BoolVal(p.config.SkipTLSVerify),
		"tls_enabled":       cty.BoolVal(p.config.TLSEnabled),
		"timeout":           cty.NumberIntVal(int64(p.config.Timeout)),
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
		"bootstrap_servers": p.config.BootstrapServers,
		"ca_cert":           p.config.CACert,
		"client_cert":       p.config.ClientCert,
		"client_key":        p.config.ClientCertKey,
		"sasl_username":     p.config.SASLUsername,
		"sasl_password":     p.config.SASLPassword,
		"sasl_mechanism":    p.config.SASLMechanism,
		"skip_tls_verify":   p.config.SkipTLSVerify,
		"tls_enabled":       p.config.TLSEnabled,
		"timeout":           p.config.Timeout,
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
