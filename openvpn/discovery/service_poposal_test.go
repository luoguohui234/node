package discovery

import (
	"encoding/json"
	nats_discovery "github.com/mysterium/node/communication/nats/discovery"
	"github.com/mysterium/node/openvpn"
	dto_openvpn "github.com/mysterium/node/openvpn/discovery/dto"
	dto_discovery "github.com/mysterium/node/service_discovery/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	nats_discovery.Bootstrap()
	openvpn.Bootstrap()
}

func TestServiceProposalUnserialize(t *testing.T) {
	jsonData := []byte(`{
		"id": 1,
		"format": "service-proposal/v1",
		"service_type": "openvpn",
		"service_definition": {},
		"payment_method_type": "PER_TIME",
		"payment_method": {},
		"provider_id": "node",
		"provider_contacts": []
	}`)

	var actual dto_discovery.ServiceProposal
	err := json.Unmarshal(jsonData, &actual)
	assert.NoError(t, err)

	expected := dto_discovery.ServiceProposal{
		ID:                1,
		Format:            "service-proposal/v1",
		ServiceType:       "openvpn",
		ServiceDefinition: dto_openvpn.ServiceDefinition{},
		PaymentMethodType: "PER_TIME",
		PaymentMethod:     dto_openvpn.PaymentPerTime{},
		ProviderID:        "node",
		ProviderContacts:  []dto_discovery.Contact{},
	}
	assert.Equal(t, expected, actual)
}

func TestServiceProposalUnserializeUnknownService(t *testing.T) {
	jsonData := []byte(`{
		"service_type": "unknown",
		"service_definition": {},
		"payment_method_type": "PER_TIME",
		"payment_method": {},
		"provider_contacts": []
	}`)

	var actual dto_discovery.ServiceProposal
	err := json.Unmarshal(jsonData, &actual)

	assert.EqualError(t, err, "Service unserializer 'unknown' doesn't exist")
	assert.Equal(t, "unknown", actual.ServiceType)
	assert.Nil(t, actual.ServiceDefinition)
}

func TestServiceProposalUnserializePerTimePaymentMethod(t *testing.T) {
	jsonData := []byte(`{
		"service_type": "openvpn",
		"service_definition": {},
		"payment_method_type": "PER_TIME",
		"payment_method": {},
		"provider_contacts": []
	}`)

	var actual dto_discovery.ServiceProposal
	err := json.Unmarshal(jsonData, &actual)

	assert.Nil(t, err)
	assert.Exactly(t, dto_openvpn.PaymentPerTime{}, actual.PaymentMethod)
}

func TestServiceProposalUnserializeUnknownPaymentMethod(t *testing.T) {
	jsonData := []byte(`{
		"service_type": "openvpn",
		"service_definition": {},
		"payment_method_type": "unknown",
		"payment_method": {},
		"provider_contacts": []
	}`)

	var actual dto_discovery.ServiceProposal
	err := json.Unmarshal(jsonData, &actual)

	assert.EqualError(t, err, "Payment method unserializer 'unknown' doesn't exist")
	assert.Equal(t, "unknown", actual.PaymentMethodType)
	assert.Nil(t, actual.PaymentMethod)
}

func TestServiceProposalSerialize(t *testing.T) {
	expectedJSON := `{
		"id": 1,
		"format": "service-proposal/v1",
		"service_type": "openvpn",
		"service_definition": {
			"location": {},
			"location_originate": {}
		},
		"payment_method_type": "PER_TIME",
		"payment_method": {
			"price": {},
			"duration": 0
		},
		"provider_id": "node",
		"provider_contacts": []
	}`

	sp := dto_discovery.ServiceProposal{
		ID:                1,
		Format:            "service-proposal/v1",
		ServiceType:       "openvpn",
		ServiceDefinition: dto_openvpn.ServiceDefinition{},
		PaymentMethodType: "PER_TIME",
		PaymentMethod:     dto_openvpn.PaymentPerTime{},
		ProviderID:        "node",
		ProviderContacts:  []dto_discovery.Contact{},
	}

	actualJSON, err := json.Marshal(sp)
	assert.NoError(t, err)
	assert.JSONEq(t, expectedJSON, string(actualJSON))
}
