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
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/IBM/sarama"
)

type KafkaService struct { //nolint
	terraformutils.Service

	client sarama.Client
}

func (p *KafkaService) Initialize() error {
	bootstrapServers := strings.Split(p.Args["bootstrap_servers"].(string), "")
	config := sarama.NewConfig()

	var err error

	if p.Args["sasl_username"] != nil {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = p.Args["sasl_username"].(string)
		config.Net.SASL.Password = p.Args["sasl_password"].(string)
		config.Net.SASL.Mechanism = sarama.SASLMechanism(p.Args["sasl_mechanism"].(string))
	}

	if p.Args["tls_enabled"] != nil {
		tlsConfig, err := newTLSConfig(
			p.Args["client_cert"].(string),
			p.Args["client_key"].(string),
			p.Args["ca_cert"].(string),
			p.Args["client_key_passphrase"].(string),
		)

		if err != nil {
			return err
		}

		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Config.InsecureSkipVerify = p.Args["skip_tls_verify"].(bool)
	}

	p.client, err = sarama.NewClient(bootstrapServers, config)

	return err
}

func parsePemOrLoadFromFile(input string) (*pem.Block, []byte, error) {
	// attempt to parse
	var inputBytes = []byte(input)
	inputBlock, _ := pem.Decode(inputBytes)

	if inputBlock == nil {
		//attempt to load from file
		log.Printf("[INFO] Attempting to load from file '%s'", input)
		var err error
		inputBytes, err = os.ReadFile(input)
		if err != nil {
			return nil, nil, err
		}
		inputBlock, _ = pem.Decode(inputBytes)
		if inputBlock == nil {
			return nil, nil, fmt.Errorf("[ERROR] Error unable to decode pem")
		}
	}
	return inputBlock, inputBytes, nil
}

func newTLSConfig(clientCert, clientKey, caCert, clientKeyPassphrase string) (*tls.Config, error) {
	tlsConfig := tls.Config{}

	if clientCert != "" && clientKey != "" {
		_, certBytes, err := parsePemOrLoadFromFile(clientCert)
		if err != nil {
			log.Printf("[ERROR] Unable to read certificate %s", err)
			return &tlsConfig, err
		}

		keyBlock, keyBytes, err := parsePemOrLoadFromFile(clientKey)
		if err != nil {
			log.Printf("[ERROR] Unable to read private key %s", err)
			return &tlsConfig, err
		}

		if x509.IsEncryptedPEMBlock(keyBlock) { //nolint:staticcheck
			log.Printf("[INFO] Using encrypted private key")
			var err error

			keyBytes, err = x509.DecryptPEMBlock(keyBlock, []byte(clientKeyPassphrase)) //nolint:staticcheck
			if err != nil {
				log.Printf("[ERROR] Error decrypting private key with passphrase %s", err)
				return &tlsConfig, err
			}
			keyBytes = pem.EncodeToMemory(&pem.Block{
				Type:  keyBlock.Type,
				Bytes: keyBytes,
			})
		}

		cert, err := tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			log.Printf("[ERROR] Error creating X509KeyPair %s", err)
			return &tlsConfig, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if caCert == "" {
		log.Println("[WARN] no CA file set skipping")
		return &tlsConfig, nil
	}

	caCertPool, _ := x509.SystemCertPool()
	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}

	_, caBytes, err := parsePemOrLoadFromFile(caCert)
	if err != nil {
		log.Printf("[ERROR] Unable to read CA %s", err)
		return &tlsConfig, err
	}
	ok := caCertPool.AppendCertsFromPEM(caBytes)
	log.Printf("[TRACE] set cert pool %v", ok)
	if !ok {
		return &tlsConfig, fmt.Errorf("Couldn't add the caPem")
	}

	tlsConfig.RootCAs = caCertPool
	return &tlsConfig, nil
}
