/*******************************************************************************
 * Copyright 2019 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package operating

import (
	"github.com/salmanazmat666/go-mod-core-contracts/models"
	"testing"
)

func TestUpdateValidation(t *testing.T) {
	tests := []struct {
		name        string
		up          UpdateRequest
		expectError bool
	}{
		{"valid - enabled", UpdateRequest{OperatingState: models.OperatingState("ENABLED")}, false},
		{"valid - disabled", UpdateRequest{OperatingState: models.OperatingState("DISABLED")}, false},
		{"invalid - blank", UpdateRequest{OperatingState: models.OperatingState("")}, true},
		{"invalid - garbage", UpdateRequest{OperatingState: models.OperatingState("QWERTY")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.up.Validate()
			if err != nil {
				if !tt.expectError {
					t.Errorf("unexpected error: %v", err)
				}
				_, ok := err.(models.ErrContractInvalid)
				if !ok {
					t.Errorf("incorrect error type returned")
				}
			}
			if tt.expectError && err == nil {
				t.Errorf("did not receive expected error: %s", tt.name)
			}
		})
	}
}
