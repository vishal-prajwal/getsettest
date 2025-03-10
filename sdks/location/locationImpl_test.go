package location

import (
	"reflect"
	"testing"
)

func TestLocationImpl_ExtractStateFromLatLong(t *testing.T) {
	locSDK := New(&LocationConfig{
		// change this address before running tests
		BaseURL:           "http://geolocation.jwr-uat-5.jungleerummyuat.com",
		DefaultAPITimeout: 10,
	})
	type args struct {
		locationrequest LocationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *LocationResponse
		wantErr bool
	}{
		{
			name: "Karnataka location",
			args: args{
				locationrequest: LocationRequest{
					Latitude:  12.93495402680534,
					Longitude: 77.63073307041866,
				},
			},
			want: &LocationResponse{
				StateName:        "Karnataka",
				StateShortCode:   "KA",
				CountryName:      "India",
				CountryShortCode: "IN",
				CityName:         "Bangalore Division",
				LocalityName:     "Bengaluru",
				PostalCode:       "560095",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := locSDK.ExtractStateFromLatLong(tt.args.locationrequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocationImpl.ExtractStateFromLatLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocationImpl.ExtractStateFromLatLong() = %v, want %v", got, tt.want)
			}
		})
	}
}
