//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package responses

import (
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"
)

// DeviceProfileResponse defines the Response Content for GET DeviceProfile DTOs.
// This object and its properties correspond to the DeviceProfileResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/DeviceProfileResponse
type DeviceProfileResponse struct {
	common.BaseResponse `json:",inline"`
	Profile             dtos.DeviceProfile `json:"profile"`
}

func NewDeviceProfileResponse(requestId string, message string, statusCode int, deviceProfile dtos.DeviceProfile) DeviceProfileResponse {
	return DeviceProfileResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Profile:      deviceProfile,
	}
}

// MultiDeviceProfilesResponse defines the Response Content for GET multiple DeviceProfile DTOs.
// This object and its properties correspond to the MultiDeviceProfilesResponse object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/MultiDeviceProfilesResponse
type MultiDeviceProfilesResponse struct {
	common.BaseResponse `json:",inline"`
	Profiles            []dtos.DeviceProfile `json:"profiles"`
}

func NewMultiDeviceProfilesResponse(requestId string, message string, statusCode int, deviceProfiles []dtos.DeviceProfile) MultiDeviceProfilesResponse {
	return MultiDeviceProfilesResponse{
		BaseResponse: common.NewBaseResponse(requestId, message, statusCode),
		Profiles:     deviceProfiles,
	}
}
