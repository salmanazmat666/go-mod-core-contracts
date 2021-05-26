//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/salmanazmat666/go-mod-core-contracts/v2"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func eventRequestData() AddEventRequest {
	return AddEventRequest{
		BaseRequest: common.BaseRequest{
			RequestId: ExampleUUID,
		},
		Event: dtos.Event{
			Id:          ExampleUUID,
			DeviceName:  TestDeviceName,
			ProfileName: TestDeviceProfileName,
			Origin:      TestOriginTime,
			Readings: []dtos.BaseReading{{
				DeviceName:   TestDeviceName,
				ResourceName: TestDeviceResourceName,
				ProfileName:  TestDeviceProfileName,
				Origin:       TestOriginTime,
				ValueType:    v2.ValueTypeUint8,
				SimpleReading: dtos.SimpleReading{
					Value: TestReadingValue,
				},
			}},
			Tags: map[string]string{
				"GatewayId": "Houston-0001",
			},
		},
	}
}

func TestAddEventRequest_Validate(t *testing.T) {
	valid := eventRequestData()
	noReqId := eventRequestData()
	noReqId.RequestId = ""
	invalidReqId := eventRequestData()
	invalidReqId.RequestId = "xxy"
	noEventId := eventRequestData()
	noEventId.Event.Id = ""
	invalidEventId := eventRequestData()
	invalidEventId.Event.Id = "gj93j2-v92hvi3h"
	noDeviceName := eventRequestData()
	noDeviceName.Event.DeviceName = ""
	noProfileName := eventRequestData()
	noProfileName.Event.ProfileName = ""
	noOrigin := eventRequestData()
	noOrigin.Event.Origin = 0

	noReading := eventRequestData()
	noReading.Event.Readings = nil

	invalidReadingNoDevice := eventRequestData()
	invalidReadingNoDevice.Event.Readings[0].DeviceName = ""
	invalidReadingNoResourceName := eventRequestData()
	invalidReadingNoResourceName.Event.Readings[0].ResourceName = ""
	invalidReadingNoProfileName := eventRequestData()
	invalidReadingNoProfileName.Event.Readings[0].ProfileName = ""
	invalidReadingNoOrigin := eventRequestData()
	invalidReadingNoOrigin.Event.Readings[0].Origin = 0

	invalidReadingNoValueType := eventRequestData()
	invalidReadingNoValueType.Event.Readings[0].ValueType = ""
	invalidReadingInvalidValueType := eventRequestData()
	invalidReadingInvalidValueType.Event.Readings[0].ValueType = "BadType"

	invalidSimpleReadingNoValue := eventRequestData()
	invalidSimpleReadingNoValue.Event.Readings[0].SimpleReading.Value = ""

	invalidBinaryReadingNoValue := eventRequestData()
	invalidBinaryReadingNoValue.Event.Readings[0].ValueType = v2.ValueTypeBinary
	invalidBinaryReadingNoValue.Event.Readings[0].BinaryReading.MediaType = TestBinaryReadingMediaType
	invalidBinaryReadingNoValue.Event.Readings[0].BinaryReading.BinaryValue = []byte{}

	invalidBinaryReadingNoMedia := eventRequestData()
	invalidBinaryReadingNoMedia.Event.Readings[0].ValueType = v2.ValueTypeBinary
	invalidBinaryReadingNoMedia.Event.Readings[0].BinaryReading.MediaType = ""
	invalidBinaryReadingNoMedia.Event.Readings[0].BinaryReading.BinaryValue = []byte(TestReadingBinaryValue)

	tests := []struct {
		name        string
		event       AddEventRequest
		expectError bool
	}{
		{"valid AddEventRequest", valid, false},
		{"valid AddEventRequest, no Request Id", noReqId, false},
		{"invalid AddEventRequest, Request Id is not an uuid", invalidReqId, true},
		{"invalid AddEventRequest, no Event Id", noEventId, true},
		{"invalid AddEventRequest, Event Id is not an uuid", invalidEventId, true},
		{"invalid AddEventRequest, no DeviceName", noDeviceName, true},
		{"invalid AddEventRequest, no ProfileName", noProfileName, true},
		{"invalid AddEventRequest, no Origin", noOrigin, true},
		{"invalid AddEventRequest, no Reading", noReading, true},
		{"invalid AddEventRequest, no Reading DeviceName", invalidReadingNoDevice, true},
		{"invalid AddEventRequest, no Resource Name", invalidReadingNoResourceName, true},
		{"invalid AddEventRequest, no Profile Name", invalidReadingNoProfileName, true},
		{"invalid AddEventRequest, no Reading Origin", invalidReadingNoOrigin, true},
		{"invalid AddEventRequest, no Reading ValueType", invalidReadingNoValueType, true},
		{"invalid AddEventRequest, invalid Reading ValueType", invalidReadingInvalidValueType, true},
		{"invalid AddEventRequest, no SimpleReading Value", invalidSimpleReadingNoValue, true},
		{"invalid AddEventRequest, no BinaryReading BinaryValue", invalidBinaryReadingNoValue, true},
		{"invalid AddEventRequest, no BinaryReading MediaType", invalidBinaryReadingNoMedia, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	type testForNameField struct {
		name        string
		event       AddEventRequest
		expectError bool
	}

	deviceNameWithUnreservedChar := eventRequestData()
	deviceNameWithUnreservedChar.Event.DeviceName = nameWithUnreservedChars
	profileNameWithUnreservedChar := eventRequestData()
	profileNameWithUnreservedChar.Event.ProfileName = nameWithUnreservedChars
	readingDeviceNameWithUnreservedChar := eventRequestData()
	readingDeviceNameWithUnreservedChar.Event.Readings[0].DeviceName = nameWithUnreservedChars
	readingResourceNameWithUnreservedChar := eventRequestData()
	readingResourceNameWithUnreservedChar.Event.Readings[0].ResourceName = nameWithUnreservedChars
	readingProfileNameWithUnreservedChar := eventRequestData()
	readingProfileNameWithUnreservedChar.Event.Readings[0].ProfileName = nameWithUnreservedChars

	// Following tests verify if name fields containing unreserved characters should pass edgex-dto-rfc3986-unreserved-chars check
	testsForNameFields := []testForNameField{
		{"Valid AddEventRequest with device name containing unreserved chars", deviceNameWithUnreservedChar, false},
		{"Valid AddEventRequest with profile name containing unreserved chars", profileNameWithUnreservedChar, false},
		{"Valid AddEventRequest with reading device name containing unreserved chars", readingDeviceNameWithUnreservedChar, false},
		{"Valid AddEventRequest with reading resource name containing unreserved chars", readingResourceNameWithUnreservedChar, false},
		{"Valid AddEventRequest with reading profile name containing unreserved chars", readingProfileNameWithUnreservedChar, false},
	}

	// Following tests verify if name fields containing reserved characters should be detected with an error
	for _, n := range namesWithReservedChar {
		deviceNameWithReservedChar := eventRequestData()
		deviceNameWithReservedChar.Event.DeviceName = n
		profileNameWithReservedChar := eventRequestData()
		profileNameWithReservedChar.Event.ProfileName = n
		readingDeviceNameWithReservedChar := eventRequestData()
		readingDeviceNameWithReservedChar.Event.Readings[0].DeviceName = n
		readingResourceNameWithReservedChar := eventRequestData()
		readingResourceNameWithReservedChar.Event.Readings[0].ResourceName = n
		readingProfileNameWithReservedChar := eventRequestData()
		readingProfileNameWithReservedChar.Event.Readings[0].ProfileName = n

		testsForNameFields = append(testsForNameFields,
			testForNameField{"Invalid AddEventRequest with device name containing reserved char", deviceNameWithReservedChar, true},
			testForNameField{"Invalid AddEventRequest with profile name containing reserved char", profileNameWithReservedChar, true},
			testForNameField{"Invalid AddEventRequest with reading device name containing reserved char", readingDeviceNameWithReservedChar, true},
			testForNameField{"Invalid AddEventRequest with reading resource name containing reserved char", readingResourceNameWithReservedChar, true},
			testForNameField{"Invalid AddEventRequest with reading profile name containing reserved char", readingProfileNameWithReservedChar, true},
		)
	}

	for _, tt := range testsForNameFields {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.event.Validate()
			if tt.expectError {
				assert.Error(t, err, fmt.Sprintf("expect error but not : %s", tt.name))
			} else {
				assert.NoError(t, err, fmt.Sprintf("unexpected error occurs : %s", tt.name))
			}
		})
	}
}

func TestAddEvent_UnmarshalJSON(t *testing.T) {
	expected := eventRequestData()
	validData, err := json.Marshal(eventRequestData())
	require.NoError(t, err)
	validValueTypeLowerCase := eventRequestData()
	validValueTypeLowerCase.Event.Readings[0].ValueType = "uint8"
	validValueTypeLowerCaseData, err := json.Marshal(validValueTypeLowerCase)
	require.NoError(t, err)
	validValueTypeUpperCase := eventRequestData()
	validValueTypeUpperCase.Event.Readings[0].ValueType = "UINT8"
	validValueTypeUpperCaseData, err := json.Marshal(validValueTypeUpperCase)
	require.NoError(t, err)

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"unmarshal AddEventRequest with success", validData, false},
		{"unmarshal AddEventRequest with success, valid value type uint8", validValueTypeLowerCaseData, false},
		{"unmarshal AddEventRequest with success, valid value type UINT8", validValueTypeUpperCaseData, false},
		{"unmarshal invalid AddEventRequest, empty data", []byte{}, true},
		{"unmarshal invalid AddEventRequest, string data", []byte("Invalid AddEventRequest"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var addEvent AddEventRequest
			err := addEvent.UnmarshalJSON(tt.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, addEvent, "Unmarshal did not result in expected AddEventRequest.")
			}
		})
	}
}

func Test_AddEventReqToEventModels(t *testing.T) {
	valid := eventRequestData()
	s := models.SimpleReading{
		BaseReading: models.BaseReading{
			DeviceName:   TestDeviceName,
			ResourceName: TestDeviceResourceName,
			ProfileName:  TestDeviceProfileName,
			Origin:       TestOriginTime,
			ValueType:    v2.ValueTypeUint8,
		},
		Value: TestReadingValue,
	}
	expectedEventModel := models.Event{
		Id:          ExampleUUID,
		DeviceName:  TestDeviceName,
		ProfileName: TestDeviceProfileName,
		Origin:      TestOriginTime,
		Readings:    []models.Reading{s},
		Tags: map[string]string{
			"GatewayId": "Houston-0001",
		},
	}

	tests := []struct {
		name        string
		addEventReq AddEventRequest
	}{
		{"valid AddEventRequest", valid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventModel := AddEventReqToEventModel(tt.addEventReq)
			assert.Equal(t, expectedEventModel, eventModel, "AddEventReqToEventModel did not result in expected Event model.")
		})
	}
}
