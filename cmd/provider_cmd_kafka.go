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

package cmd

import (
	"errors"
	"os"
	"strconv"

	kafka_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/kafka"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/spf13/cobra"
)

const (
	defaultKafkaBootStrapServer = "localhost"
	defaultKafkaTimeout         = 120
)

func newCmdKafkaImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "Import current state to Terraform configuration from Kafka",
		Long:  "Import current state to Terraform configuration from Kafka",
		RunE: func(cmd *cobra.Command, args []string) error {
			bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
			if len(bootstrapServers) == 0 {
				endpoint = defaultKafkaBootStrapServer
			}
			caCert := os.Getenv("KAFKA_CA_CERT")
			clientCert := os.Getenv("KAFKA_CLIENT_CERT")
			clientKey := os.Getenv("KAFKA_CLIENT_KEY")
			saslUsername := os.Getenv("KAFKA_SASL_USERNAME")
			saslPassword := os.Getenv("KAFKA_SASL_PASSWORD")
			saslMechanism := os.Getenv("KAFKA_SASL_MECHANISM")
			switch saslMechanism {
				case "scram-sha512", "scram-sha256", "plain":
				default:
					return errors.New("Invalid sasl mechanism \"%s\": can only be \"scram-sha256\", \"scram-sha512\" or \"plain\"", saslMechanism)
			}
			skipTLSVerify := strconv.ParseBool(os.Getenv("KAFKA_SKIP_VERIFY"))
			tlsEnabled := strconv.ParseBool(os.Getenv("KAFKA_ENABLE_TLS"))
			timeout := strconv.Atoi(os.Getenv("KAFKA_TIMEOUT"))
			if timeout == 0 {
				timemout = defaultKafkaTimeout
			}
			provider := newKafkaProvider()
			err := Import(provider, options, []string{bootstrapServers, caCert, clientCert, clientKey, saslUsername, saslPassword, saslMechanism, skipTLSVerify, tlsEnabled, timeout})
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.AddCommand(listCmd(newKafkaProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "topics", "kafka_type=id1:id2:id4")
	return cmd
}

func newKafkaProvider() terraform_utils.ProviderGenerator {
	return &kafka_terraforming.KafkaProvider{}
}
