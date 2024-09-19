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
	"github.com/IBM/sarama"
)

type QuotaGenerator struct {
	KafkaService
}

type Quota struct {
	EntityName string `json:"entity_name"`
	EntityType string `json:"entity_type"`
}

type Quotas []Quota

var QuotaAllowEmptyValues = []string{}
var QuotaAdditionalFields = map[string]interface{}{}

func (q QuotaGenerator) createResources(quotas Quotas) []terraformutils.Resource {
	var resources []terraformutils.Resource
	return resources
}

func (q *QuotaGenerator) InitResources() error {
	broker, err := q.client.Controller()
	if err != nil {
		return err
	}

	quotas, err = broker.DescribeClientQuotas(sarama.NewMetadataRequest())
	if err != nil {
		return err
	}

	q.Resources = q.createResources(quotas)
	return nil
}
