//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"
	"os"

	"github.com/fxamacker/cbor/v2"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/models"
)

// AddEventRequest defines the Request Content for POST event DTO.
// This object and its properties correspond to the AddEventRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-data/2.x#/AddEventRequest
type AddEventRequest struct {
	common.BaseRequest `json:",inline"`
	Event              dtos.Event `json:"event" validate:"required"`
}

// NewAddEventRequest creates, initializes and returns an AddEventRequests
func NewAddEventRequest(event dtos.Event) AddEventRequest {
	return AddEventRequest{
		BaseRequest: common.NewBaseRequest(),
		Event:       event,
	}
}

// Validate satisfies the Validator interface
func (a AddEventRequest) Validate() error {
	if err := v2.Validate(a); err != nil {
		return err
	}

	// BaseReading has the skip("-") validation annotation for BinaryReading and SimpleReading
	// Otherwise error will occur as only one of them exists
	// Therefore, need to validate the nested BinaryReading and SimpleReading struct here
	for _, r := range a.Event.Readings {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type unmarshal func([]byte, interface{}) error

func (a *AddEventRequest) UnmarshalJSON(b []byte) error {
	return a.Unmarshal(b, json.Unmarshal)
}

func (a *AddEventRequest) UnmarshalCBOR(b []byte) error {
	return a.Unmarshal(b, cbor.Unmarshal)
}

func (a *AddEventRequest) Unmarshal(b []byte, f unmarshal) error {
	// To avoid recursively invoke unmarshaler interface, intentionally create a struct to represent AddEventRequest DTO
	var addEvent struct {
		common.BaseRequest
		Event dtos.Event
	}
	if err := f(b, &addEvent); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal the byte array.", err)
	}

	*a = AddEventRequest(addEvent)

	// validate AddEventRequest DTO
	if err := a.Validate(); err != nil {
		return err
	}

	// Normalize reading's value type
	for i, r := range a.Event.Readings {
		valueType, err := v2.NormalizeValueType(r.ValueType)
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		}
		a.Event.Readings[i].ValueType = valueType
	}
	return nil
}

func (a *AddEventRequest) Encode() ([]byte, string, error) {
	var encoding = clients.ContentTypeJSON

	for _, r := range a.Event.Readings {
		if r.ValueType == v2.ValueTypeBinary {
			encoding = clients.ContentTypeCBOR
			break
		}
	}
	if v := os.Getenv(v2.EnvEncodeAllEvents); v == v2.ValueTrue {
		encoding = clients.ContentTypeCBOR
	}

	var err error
	var encodedData []byte
	switch encoding {
	case clients.ContentTypeCBOR:
		encodedData, err = cbor.Marshal(a)
		if err != nil {
			return nil, "", errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to encode AddEventRequest to CBOR", err)
		}
	case clients.ContentTypeJSON:
		encodedData, err = json.Marshal(a)
		if err != nil {
			return nil, "", errors.NewCommonEdgeX(errors.KindContractInvalid, "failed to encode AddEventRequest to JSON", err)
		}
	}

	return encodedData, encoding, nil
}

// AddEventReqToEventModel transforms the AddEventRequest DTO to the Event model
func AddEventReqToEventModel(addEventReq AddEventRequest) (event models.Event) {
	readings := make([]models.Reading, len(addEventReq.Event.Readings))
	for i, r := range addEventReq.Event.Readings {
		readings[i] = dtos.ToReadingModel(r)
	}

	tags := make(map[string]string)
	for tag, value := range addEventReq.Event.Tags {
		tags[tag] = value
	}

	return models.Event{
		Id:          addEventReq.Event.Id,
		DeviceName:  addEventReq.Event.DeviceName,
		ProfileName: addEventReq.Event.ProfileName,
		SourceName:  addEventReq.Event.SourceName,
		Origin:      addEventReq.Event.Origin,
		Readings:    readings,
		Tags:        tags,
	}
}
