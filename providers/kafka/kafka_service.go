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
	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
//	"github.com/Mongey/terraform-provider-kafka/kafka"
//	"github.com/Shopify/sarama"
//	"strings"
)

type Config struct {
	BootstrapServers *[]string
	Timeout          int
	CACertFile       string
	ClientCertFile   string
	ClientCertKey    string
	TLSEnabled       bool
	SkipTLSVerify    bool
	SASLUsername     string
	SASLPassword     string
	SASLMechanism    string
}

type KafkaService struct {
	terraform_utils.Service
}

// func (s *KafkaService) kafkaConfig() (*sarama.Config, error) {
// 	bootstrapServers := strings.Split(s.Args["bootstrap_servers"].(string), ",")
// 	config := Config{
// 		BootstrapServers: &bootstrapServers,
// 		CACertFile:       s.Args["ca_cert"].(string),
// 		ClientCertFile:   s.Args["client_cert"].(string),
// 		ClientCertKey:    s.Args["client_key"].(string),
// 		SkipTLSVerify:    s.Args["skip_tls_verify"].(bool),
// 		SASLUsername:     s.Args["sasl_username"].(string),
// 		SASLPassword:     s.Args["sasl_password"].(string),
// 		SASLMechanism:    s.Args["sasl_mechanism"].(string),
// 		TLSEnabled:       s.Args["tls_enabled"].(bool),
// 		Timeout:          s.Args["tls_enabled"].(int),
// 	}
// 	//return kafka.NewKafkaConfig(config)
// }
