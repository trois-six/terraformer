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
	"fmt"
	"log"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
	"github.com/Shopify/sarama"
)

type TopicGenerator struct {
	KafkaService
}

var TopicAllowEmptyValues = []string{}

func (g TopicGenerator) createResources(topics []string) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, topic := range topics {
		resources = append(resources, terraform_utils.NewSimpleResource(
			topic, //ID
			fmt.Sprintf("topic_%s", normalizeResourceName(topic)), //resourceName
			"kafka_topic", //resourceType
			"kafka",       //provider
			TopicAllowEmptyValues,
		))
	}
	return resources
}

func (g *TopicGenerator) InitResources() error {
	var topics []string

	bootstrapServers := strings.Split(g.Args["bootstrap_servers"].(string), ",")
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin(bootstrapServers, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	defer func() { _ = admin.Close() }()
	topicslist, _ := admin.ListTopics()
	for topic := range topicslist {
		topics = append(topics, topic)
	}
	g.Resources = g.createResources(topics)
	return nil
}
