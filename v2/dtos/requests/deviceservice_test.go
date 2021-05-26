//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testAddDeviceService = AddDeviceServiceRequest{
	BaseRequest: common.BaseRequest{
		RequestId: ExampleUUID,
	},
	Service: dtos.DeviceService{
		Name:        TestDeviceServiceName,
		BaseAddress: TestBaseAddress,
		Labels:      []string{"MODBUS", "TEMP"},
		AdminState:  models.Locked,
	},
}

var testUpdateDeviceService = UpdateDeviceServiceRequest{
	BaseRequest: common.BaseRequest{
		RequestId: ExampleUUID,
	},
	Service: mockDeviceServiceDTO(),
}

func mockDeviceServiceDTO() dtos.UpdateDeviceService {
	testUUID := ExampleUUID
	testName := TestDeviceServiceName
	testBaseAddress := TestBaseAddress
	testAdminState := models.Locked
	ds := dtos.UpdateDeviceService{}
	ds.Id = &testUUID
	ds.Name = &testName
	ds.BaseAddress = &testBaseAddress
	ds.AdminState = &testAdminState
	ds.Labels = testLabels
	return ds
}

func TestAddDeviceServiceRequest_Validate(t *testing.T) {
	emptyString := " "
	valid := testAddDeviceService
	noReqId := testAddDeviceService
	noReqId.RequestId = ""
	invalidReqId := testAddDeviceService
	invalidReqId.RequestId = "jfdw324"
	noName := testAddDeviceService
	noName.Service.Name = emptyString
	noAdminState := testAddDeviceService
	noAdminState.Service.AdminState = emptyString
	invalidAdminState := testAddDeviceService
	invalidAdminState.Service.AdminState = "invalid"
	noBaseAddress := testAddDeviceService
	noBaseAddress.Service.BaseAddress = emptyString
	invalidBaseAddress := testAddDeviceService
	invalidBaseAddress.Service.BaseAddress = "invalid"
	tests := []struct {
		name          string
		DeviceService AddDeviceServiceRequest
		expectError   bool
	}{
		{"valid AddDeviceServiceRequest", valid, false},
		{"valid AddDeviceServiceRequest, no Request Id", noReqId, false},
		{"invalid AddDeviceServiceRequest, Request Id is not an uuid", invalidReqId, true},
		{"invalid AddDeviceServiceRequest, no Name", noName, true},
		{"invalid AddDeviceServiceRequest, no AdminState", noAdminState, true},
		{"invalid AddDeviceServiceRequest, invalid AdminState", invalidAdminState, true},
		{"invalid AddDeviceServiceRequest, no BaseAddress", noBaseAddress, true},
		{"invalid AddDeviceServiceRequest, no BaseAddress", invalidBaseAddress, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.DeviceService.Validate()
			assert.Equal(t, tt.expectError, err != nil, "Unexpected addDeviceServiceRequest validation result.", err)
		})
	}

	serviceNameWithUnreservedChars := testAddDeviceService
	serviceNameWithUnreservedChars.Service.Name = nameWithUnreservedChars

	err := serviceNameWithUnreservedChars.Validate()
	assert.NoError(t, err, fmt.Sprintf("AddDeviceServiceRequest with service name containing unreserved chars %s should pass validation", nameWithUnreservedChars))

	// Following tests verify if service name containing reserved characters should be detected with an error
	for _, n := range namesWithReservedChar {
		serviceNameWithReservedChar := testAddDeviceService
		serviceNameWithReservedChar.Service.Name = n

		err := serviceNameWithReservedChar.Validate()
		assert.Error(t, err, fmt.Sprintf("AddDeviceServiceRequest with service name containing reserved char %s should return error during validation", n))
	}
}

func TestAddDeviceService_UnmarshalJSON(t *testing.T) {
	valid := testAddDeviceService
	resultTestBytes, _ := json.Marshal(testAddDeviceService)
	type args struct {
		data []byte
	}
	tests := []struct {
		name             string
		addDeviceService AddDeviceServiceRequest
		args             args
		wantErr          bool
	}{
		{"unmarshal AddDeviceServiceRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid AddDeviceServiceRequest, empty data", AddDeviceServiceRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid AddDeviceServiceRequest, string data", AddDeviceServiceRequest{}, args{[]byte("Invalid AddDeviceServiceRequest")}, true},
	}
	fmt.Println(string(resultTestBytes))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.addDeviceService
			err := tt.addDeviceService.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, tt.addDeviceService, "Unmarshal did not result in expected AddDeviceServiceRequest.")
			}
		})
	}
}

func TestAddDeviceServiceReqToDeviceServiceModels(t *testing.T) {
	requests := []AddDeviceServiceRequest{testAddDeviceService}
	expectedDeviceServiceModel := []models.DeviceService{{
		Name:        TestDeviceServiceName,
		BaseAddress: TestBaseAddress,
		Labels:      []string{"MODBUS", "TEMP"},
		AdminState:  models.Locked,
	}}
	resultModels := AddDeviceServiceReqToDeviceServiceModels(requests)
	assert.Equal(t, expectedDeviceServiceModel, resultModels, "AddDeviceServiceReqToDeviceServiceModels did not result in expected DeviceService model.")
}

func TestUpdateDeviceService_UnmarshalJSON(t *testing.T) {
	valid := testUpdateDeviceService
	resultTestBytes, _ := json.Marshal(testUpdateDeviceService)
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		req     UpdateDeviceServiceRequest
		args    args
		wantErr bool
	}{
		{"unmarshal UpdateDeviceServiceRequest with success", valid, args{resultTestBytes}, false},
		{"unmarshal invalid UpdateDeviceServiceRequest, empty data", UpdateDeviceServiceRequest{}, args{[]byte{}}, true},
		{"unmarshal invalid UpdateDeviceServiceRequest, string data", UpdateDeviceServiceRequest{}, args{[]byte("Invalid UpdateDeviceServiceRequest")}, true},
	}
	fmt.Println(string(resultTestBytes))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expected = tt.req
			err := tt.req.UnmarshalJSON(tt.args.data)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expected, tt.req, "Unmarshal did not result in expected UpdateDeviceServiceRequest.", err)
			}
		})
	}
}

func TestUpdateDeviceServiceRequest_Validate(t *testing.T) {
	emptyString := " "
	invalidUUID := "invalidUUID"
	invalidUrl := "http127.0.0.1"

	valid := testUpdateDeviceService
	noReqId := valid
	noReqId.RequestId = ""
	invalidReqId := valid
	invalidReqId.RequestId = invalidUUID

	// id
	validOnlyId := valid
	validOnlyId.Service.Name = nil
	invalidId := valid
	invalidId.Service.Id = &invalidUUID
	// name
	validOnlyName := valid
	validOnlyName.Service.Id = nil
	invalidEmptyName := valid
	invalidEmptyName.Service.Name = &emptyString
	// no id and name
	noIdAndName := valid
	noIdAndName.Service.Id = nil
	noIdAndName.Service.Name = nil
	// baseAddress
	validNilBaseAddress := valid
	validNilBaseAddress.Service.BaseAddress = nil
	invalidBaseAddress := valid
	invalidBaseAddress.Service.BaseAddress = &invalidUrl

	invalid := "invalid"
	invalidAdminState := valid
	invalidAdminState.Service.AdminState = &invalid
	tests := []struct {
		name        string
		req         UpdateDeviceServiceRequest
		expectError bool
	}{
		{"valid", valid, false},
		{"valid, no Request Id", noReqId, false},
		{"invalid, Request Id is not an uuid", invalidReqId, true},

		{"valid, only id", validOnlyId, false},
		{"invalid, invalid Id", invalidId, true},
		{"valid, only name", validOnlyName, false},
		{"invalid, empty name", invalidEmptyName, true},

		{"invalid, no Id and name", noIdAndName, true},

		{"valid, nil baseAddress", validNilBaseAddress, false},
		{"invalid, invalid baseAddress", invalidBaseAddress, true},

		{"invalid, invalid AdminState", invalidAdminState, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.expectError, err != nil, "Unexpected updateDeviceServiceRequest validation result.", err)
		})
	}

	serviceNameWithUnreservedChars := testUpdateDeviceService
	serviceNameWithUnreservedChars.Service.Name = &nameWithUnreservedChars

	err := serviceNameWithUnreservedChars.Validate()
	assert.NoError(t, err, fmt.Sprintf("UpdateDeviceServiceRequest with service name containing unreserved chars %s should pass validation", nameWithUnreservedChars))

	// Following tests verify if service name containing reserved characters should be detected with an error
	for _, n := range namesWithReservedChar {
		serviceNameWithReservedChar := testUpdateDeviceService
		serviceNameWithReservedChar.Service.Name = &n

		err := serviceNameWithReservedChar.Validate()
		assert.Error(t, err, fmt.Sprintf("UpdateDeviceServiceRequest with service name containing reserved char %s should return error during validation", n))
	}
}

func TestUpdateDeviceServiceRequest_UnmarshalJSON_NilField(t *testing.T) {
	reqJson := `{
        "requestId":"7a1707f0-166f-4c4b-bc9d-1d54c74e0137",
		"service":{
			"name":"test device service"
		}
	}`
	var req UpdateDeviceServiceRequest

	err := req.UnmarshalJSON([]byte(reqJson))

	require.NoError(t, nil, err)
	// Nil field checking is used to update with patch
	assert.Nil(t, req.Service.BaseAddress)
	assert.Nil(t, req.Service.AdminState)
	assert.Nil(t, req.Service.Labels)
}

func TestUpdateDeviceServiceRequest_UnmarshalJSON_EmptySlice(t *testing.T) {
	reqJson := `{
        "requestId":"7a1707f0-166f-4c4b-bc9d-1d54c74e0137",
		"service":{
			"name":"TestDevice",
			"labels":[]
		}
	}`
	var req UpdateDeviceServiceRequest

	err := req.UnmarshalJSON([]byte(reqJson))

	require.NoError(t, err)
	// Empty slice is used to remove the data
	assert.NotNil(t, req.Service.Labels)
}

func TestReplaceDeviceServiceModelFieldsWithDTO(t *testing.T) {
	ds := models.DeviceService{
		Id:   "7a1707f0-166f-4c4b-bc9d-1d54c74e0137",
		Name: "test device service",
	}
	patch := mockDeviceServiceDTO()

	ReplaceDeviceServiceModelFieldsWithDTO(&ds, patch)

	assert.Equal(t, TestBaseAddress, ds.BaseAddress)
	assert.Equal(t, models.Locked, string(ds.AdminState))
	assert.Equal(t, testLabels, ds.Labels)
}
