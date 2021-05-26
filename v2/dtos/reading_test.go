//
// Copyright (C) 2020 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"testing"

	v2 "github.com/salmanazmat666/go-mod-core-contracts/v2"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/models"

	"github.com/stretchr/testify/assert"
)

var testSimpleReading = BaseReading{
	DeviceName:   TestDeviceName,
	ResourceName: TestReadingName,
	ProfileName:  TestDeviceProfileName,
	Origin:       TestTimestamp,
	ValueType:    TestValueType,
	SimpleReading: SimpleReading{
		Value: TestValue,
	},
}

func Test_ToReadingModel(t *testing.T) {
	valid := testSimpleReading
	expectedSimpleReading := models.SimpleReading{
		BaseReading: models.BaseReading{
			DeviceName:   TestDeviceName,
			ResourceName: TestReadingName,
			ProfileName:  TestDeviceProfileName,
			Origin:       TestTimestamp,
			ValueType:    TestValueType,
		},
		Value: TestValue,
	}
	tests := []struct {
		name    string
		reading BaseReading
	}{
		{"valid Reading", valid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readingModel := ToReadingModel(tt.reading)
			assert.Equal(t, expectedSimpleReading, readingModel, "ToReadingModel did not result in expected Reading model.")
		})
	}
}

func TestFromReadingModelToDTO(t *testing.T) {
	valid := models.SimpleReading{
		BaseReading: models.BaseReading{
			Id:           TestUUID,
			Created:      TestTimestamp,
			Origin:       TestTimestamp,
			DeviceName:   TestDeviceName,
			ResourceName: TestReadingName,
			ProfileName:  TestDeviceProfileName,
			ValueType:    TestValueType,
		},
		Value: TestValue,
	}
	expectedDTO := BaseReading{
		Versionable:  common.Versionable{ApiVersion: v2.ApiVersion},
		Id:           TestUUID,
		Created:      TestTimestamp,
		Origin:       TestTimestamp,
		DeviceName:   TestDeviceName,
		ResourceName: TestReadingName,
		ProfileName:  TestDeviceProfileName,
		ValueType:    TestValueType,
		SimpleReading: SimpleReading{
			Value: TestValue,
		},
	}

	tests := []struct {
		name    string
		reading models.Reading
	}{
		{"success to convert from reading model to DTO ", valid},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromReadingModelToDTO(tt.reading)
			assert.Equal(t, expectedDTO, result, "FromReadingModelToDTO did not result in expected Reading DTO.")
		})
	}
}
