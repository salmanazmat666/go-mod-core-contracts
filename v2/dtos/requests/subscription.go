//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package requests

import (
	"encoding/json"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/models"
)

var supportedChannelTypes = []string{v2.EMAIL, v2.REST}

// AddSubscriptionRequest defines the Request Content for POST Subscription DTO.
// This object and its properties correspond to the AddSubscriptionRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-notifications/2.x#/AddSubscriptionRequest
type AddSubscriptionRequest struct {
	common.BaseRequest `json:",inline"`
	Subscription       dtos.Subscription `json:"subscription"`
}

// Validate satisfies the Validator interface
func (request AddSubscriptionRequest) Validate() error {
	err := v2.Validate(request)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	for _, c := range request.Subscription.Channels {
		err = c.Validate()
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		} else if !contains(supportedChannelTypes, c.Type) {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "MQTT is not valid type for Channel", nil)
		}
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the AddSubscriptionRequest type
func (request *AddSubscriptionRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		Subscription dtos.Subscription
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*request = AddSubscriptionRequest(alias)

	// validate AddSubscriptionRequest DTO
	if err := request.Validate(); err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// AddSubscriptionReqToSubscriptionModels transforms the AddSubscriptionRequest DTO array to the AddSubscriptionRequest model array
func AddSubscriptionReqToSubscriptionModels(reqs []AddSubscriptionRequest) (s []models.Subscription) {
	for _, req := range reqs {
		d := dtos.ToSubscriptionModel(req.Subscription)
		s = append(s, d)
	}
	return s
}

// UpdateSubscriptionRequest defines the Request Content for PUT event as pushed DTO.
// This object and its properties correspond to the UpdateSubscriptionRequest object in the APIv2 specification:
// https://app.swaggerhub.com/apis-docs/EdgeXFoundry1/support-notifications/2.x#/UpdateSubscriptionRequest
type UpdateSubscriptionRequest struct {
	common.BaseRequest `json:",inline"`
	Subscription       dtos.UpdateSubscription `json:"subscription"`
}

// Validate satisfies the Validator interface
func (request UpdateSubscriptionRequest) Validate() error {
	err := v2.Validate(request)
	if err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	for _, c := range request.Subscription.Channels {
		err = c.Validate()
		if err != nil {
			return errors.NewCommonEdgeXWrapper(err)
		} else if !contains(supportedChannelTypes, c.Type) {
			return errors.NewCommonEdgeX(errors.KindContractInvalid, "MQTT is not valid type for Channel", nil)
		}
	}
	if request.Subscription.Categories != nil && request.Subscription.Labels != nil &&
		len(request.Subscription.Categories) == 0 && len(request.Subscription.Labels) == 0 {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "categories and labels can not be both empty", nil)
	}
	return nil
}

// UnmarshalJSON implements the Unmarshaler interface for the UpdateSubscriptionRequest type
func (request *UpdateSubscriptionRequest) UnmarshalJSON(b []byte) error {
	var alias struct {
		common.BaseRequest
		Subscription dtos.UpdateSubscription
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return errors.NewCommonEdgeX(errors.KindContractInvalid, "Failed to unmarshal request body as JSON.", err)
	}

	*request = UpdateSubscriptionRequest(alias)

	// validate UpdateSubscriptionRequest DTO
	if err := request.Validate(); err != nil {
		return errors.NewCommonEdgeXWrapper(err)
	}
	return nil
}

// ReplaceSubscriptionModelFieldsWithDTO replace existing Subscription's fields with DTO patch
func ReplaceSubscriptionModelFieldsWithDTO(s *models.Subscription, patch dtos.UpdateSubscription) {
	if patch.Channels != nil {
		s.Channels = dtos.ToAddressModels(patch.Channels)
	}
	if patch.Categories != nil {
		s.Categories = patch.Categories
	}
	if patch.Labels != nil {
		s.Labels = patch.Labels
	}
	if patch.Description != nil {
		s.Description = *patch.Description
	}
	if patch.Receiver != nil {
		s.Receiver = *patch.Receiver
	}
	if patch.ResendLimit != nil {
		s.ResendLimit = *patch.ResendLimit
	}
	if patch.ResendInterval != nil {
		s.ResendInterval = *patch.ResendInterval
	}
	if patch.AdminState != nil {
		s.AdminState = models.AdminState(*patch.AdminState)
	}
}

func NewAddSubscriptionRequest(dto dtos.Subscription) AddSubscriptionRequest {
	return AddSubscriptionRequest{
		BaseRequest:  common.NewBaseRequest(),
		Subscription: dto,
	}
}

func NewUpdateSubscriptionRequest(dto dtos.UpdateSubscription) UpdateSubscriptionRequest {
	return UpdateSubscriptionRequest{
		BaseRequest:  common.NewBaseRequest(),
		Subscription: dto,
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
