package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
)

func Test_getUsers(t *testing.T) {
	tests := []struct {
		name        string
		config      *stns.Config
		params      map[string]string
		wantErr     bool
		wantStatus  int
		wantID      int
		wantRecords int
	}{
		{
			name:   "id ok",
			config: testConfig(),
			params: map[string]string{
				"id": "1",
			},
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantID:      1,
			wantRecords: 1,
		},
		{
			name:   "id notfound",
			config: testConfig(),
			params: map[string]string{
				"id": "999999",
			},
			wantErr:    false,
			wantStatus: http.StatusNotFound,
		},
		{
			name:   "bad request",
			config: testConfig(),
			params: map[string]string{
				"ng": "999999",
			},
			wantErr:    false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:   "name ok",
			config: testConfig(),
			params: map[string]string{
				"name": "user1",
			},
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantID:      1,
			wantRecords: 1,
		},
		{
			name:   "name notfound",
			config: testConfig(),
			params: map[string]string{
				"name": "hoge",
			},
			wantErr:    false,
			wantStatus: http.StatusNotFound,
		},
		{
			name:        "all ok",
			config:      testConfig(),
			wantErr:     false,
			wantStatus:  http.StatusOK,
			wantID:      1,
			wantRecords: 2,
		},
	}
	for _, tt := range tests {
		ctx, rec := newContext("/users", tt.params, tt.config)

		t.Run(tt.name, func(t *testing.T) {
			if err := getUsers(ctx); (err != nil) != tt.wantErr {
				t.Errorf("getUsers() error = %v, wantErr %v", err, tt.wantErr)
			}

			if rec.Code != tt.wantStatus {
				t.Errorf("getUsers status code does not match, expected %d, got %d", tt.wantStatus, rec.Code)
			}

			if tt.wantID != 0 {
				users := []model.User{}
				if err := json.Unmarshal(rec.Body.Bytes(), &users); err != nil {
					t.Errorf(err.Error())
				}

				if len(users) != tt.wantRecords {
					t.Error("getUsers Record has not been acquired")
				}

				if len(users) == 1 && users[0].ID != tt.wantID {
					t.Errorf("getUsers ID does not match, expected %d, got %d", tt.wantID, users[0].ID)

				}
			}
		})
	}
}
