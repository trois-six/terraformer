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

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

type TopicGenerator struct {
	KafkaService
}

type Topic struct {
	Name              string
	Partitions        int32
	ReplicationFactor int16
	Config            map[string]*string
}

type Topics []Topic

var TopicAllowEmptyValues = []string{}

func (g TopicGenerator) createResources(topics Topics) []terraform_utils.Resource {
	var resources []terraform_utils.Resource
	for _, topic := range topics {
		resources = append(resources, terraform_utils.NewSimpleResource(
			topic.Name,
			fmt.Sprintf("topic_%s", normalizeResourceName(topic.Name)),
			"kafka_topic",
			"kafka",
			TopicAllowEmptyValues,
		))
	}
	return resources
}

func (g *TopicGenerator) InitResources() error {
	var topics Topics
	/* TODO
	 */
	g.Resources = g.createResources(topics)
	return nil
}
