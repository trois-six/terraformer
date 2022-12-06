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
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type TopicGenerator struct {
	KafkaService
}

type Topic struct {
	Name              string `json:"name"`
	Partition         string `json:"partition"`
	ReplicationFactor string `json:"replication_factor"`
}

type Topics []Topic

var TopicAllowEmptyValues = []string{}
var TopicAdditionalFields = map[string]interface{}{}

func (t TopicGenerator) createResources(topics Topics) []terraformutils.Resource {
	var resources []terraformutils.Resource
	return resources
}

func (t *TopicGenerator) InitResources() error {
	var topics Topics
	t.Resources = t.createResources(topics)
	return nil
}
