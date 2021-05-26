//
// Copyright (C) 2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package dtos

import (
	"github.com/salmanazmat666/go-mod-core-contracts/v2/dtos/common"
	"github.com/salmanazmat666/go-mod-core-contracts/v2/models"
)

// ProvisionWatcher and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/ProvisionWatcher
type ProvisionWatcher struct {
	common.Versionable  `json:",inline"`
	Id                  string              `json:"id,omitempty" validate:"omitempty,uuid"`
	Name                string              `json:"name" validate:"required,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	Labels              []string            `json:"labels,omitempty"`
	Identifiers         map[string]string   `json:"identifiers" validate:"gt=0,dive,keys,required,endkeys,required"`
	BlockingIdentifiers map[string][]string `json:"blockingIdentifiers,omitempty"`
	ProfileName         string              `json:"profile" validate:"required,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	ServiceName         string              `json:"service" validate:"required,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	AdminState          string              `json:"adminState" validate:"oneof='LOCKED' 'UNLOCKED'"`
	AutoEvents          []AutoEvent         `json:"autoEvents,omitempty" validate:"dive"`
}

// UpdateProvisionWatcher and its properties are defined in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/core-metadata/2.x#/UpdateProvisionWatcher
type UpdateProvisionWatcher struct {
	common.Versionable  `json:",inline"`
	Id                  *string             `json:"id" validate:"required_without=Name,edgex-dto-uuid"`
	Name                *string             `json:"name" validate:"required_without=Id,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	Labels              []string            `json:"labels"`
	Identifiers         map[string]string   `json:"identifiers" validate:"omitempty,gt=0,dive,keys,required,endkeys,required"`
	BlockingIdentifiers map[string][]string `json:"blockingIdentifiers"`
	ProfileName         *string             `json:"profileName" validate:"omitempty,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	ServiceName         *string             `json:"serviceName" validate:"omitempty,edgex-dto-none-empty-string,edgex-dto-rfc3986-unreserved-chars"`
	AdminState          *string             `json:"adminState" validate:"omitempty,oneof='LOCKED' 'UNLOCKED'"`
	AutoEvents          []AutoEvent         `json:"autoEvents" validate:"dive"`
}

// ToProvisionWatcherModel transforms the ProvisionWatcher DTO to the ProvisionWatcher model
func ToProvisionWatcherModel(dto ProvisionWatcher) models.ProvisionWatcher {
	return models.ProvisionWatcher{
		Timestamps:          models.Timestamps{},
		Id:                  dto.Id,
		Name:                dto.Name,
		Labels:              dto.Labels,
		Identifiers:         dto.Identifiers,
		BlockingIdentifiers: dto.BlockingIdentifiers,
		ProfileName:         dto.ProfileName,
		ServiceName:         dto.ServiceName,
		AdminState:          models.AdminState(dto.AdminState),
		AutoEvents:          ToAutoEventModels(dto.AutoEvents),
	}
}

// FromProvisionWatcherModelToDTO transforms the ProvisionWatcher Model to the ProvisionWatcher DTO
func FromProvisionWatcherModelToDTO(pw models.ProvisionWatcher) ProvisionWatcher {
	return ProvisionWatcher{
		Id:                  pw.Id,
		Name:                pw.Name,
		Labels:              pw.Labels,
		Identifiers:         pw.Identifiers,
		BlockingIdentifiers: pw.BlockingIdentifiers,
		ProfileName:         pw.ProfileName,
		ServiceName:         pw.ServiceName,
		AdminState:          string(pw.AdminState),
		AutoEvents:          FromAutoEventModelsToDTOs(pw.AutoEvents),
	}
}
