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
	"os"
	"strings"

	kafka_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/kafka"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/spf13/cobra"
)

const (
	defaultKafkaBootstrapServers = "localhost:9092"
	defaultKafkaSASLMechanism    = "plain"
	defaultKafkaSkipTLSVerify    = "false"
	defaultKafkaTLSEnabled       = "true"
)

func newCmdKafkaImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kafka",
		Short: "Import current state to Terraform configuration from Kafka",
		Long:  "Import current state to Terraform configuration from Kafka",
		RunE: func(cmd *cobra.Command, args []string) error {
			bootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
			if len(bootstrapServers) == 0 {
				bootstrapServers = defaultKafkaBootstrapServers
			}

			caCert := strings.TrimSpace(os.Getenv("KAFKA_CA_CERT"))
			clientCert := strings.TrimSpace(os.Getenv("KAFKA_CLIENT_CERT"))
			clientKey := strings.TrimSpace(os.Getenv("KAFKA_CLIENT_KEY"))
			clientKeyPassphrase := strings.TrimSpace(os.Getenv("KAFKA_CLIENT_KEY_PASSPHRASE"))
			saslUsername := strings.TrimSpace(os.Getenv("KAFKA_SASL_USERNAME"))
			saslPassword := strings.TrimSpace(os.Getenv("KAFKA_SASL_PASSWORD"))

			saslMechanism := strings.TrimSpace(os.Getenv("KAFKA_SASL_MECHANISM"))
			if len(saslMechanism) == 0 {
				saslMechanism = defaultKafkaSASLMechanism
			}

			skipTLSVerify := strings.TrimSpace(os.Getenv("KAFKA_SKIP_VERIFY"))
			if len(skipTLSVerify) == 0 {
				skipTLSVerify = defaultKafkaSkipTLSVerify
			}

			tlsEnabled := strings.TrimSpace(os.Getenv("KAFKA_SKIP_VERIFY"))
			if len(tlsEnabled) == 0 {
				tlsEnabled = defaultKafkaTLSEnabled
			}

			provider := newKafkaProvider()
			err := Import(provider, options, []string{
				bootstrapServers,
				caCert,
				clientCert,
				clientKey,
				clientKeyPassphrase,
				saslUsername,
				saslPassword,
				saslMechanism,
				skipTLSVerify,
				tlsEnabled,
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(listCmd(newKafkaProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "acls,quotas,topics", "")

	return cmd
}

func newKafkaProvider() terraformutils.ProviderGenerator {
	return &kafka_terraforming.KafkaProvider{}
}
