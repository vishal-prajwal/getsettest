package froller

import (
	"reflect"
	"testing"
)

func TestNewFeatureRoller(t *testing.T) {
	type args struct {
		features []Feature[string]
	}
	tests := []struct {
		name    string
		args    args
		want    FeatureRoller[string]
		wantErr bool
	}{
		{
			name: "valid features",
			args: args{
				features: []Feature[string]{
					NewFeature("feature1", 50),
					NewFeature("feature2", 30),
					NewFeature("feature3", 20),
				},
			},
			want: &featureRollerImpl[string]{
				features: map[string]int{
					"feature1": 50,
					"feature2": 30,
					"feature3": 20,
				},
			},
			wantErr: false,
		},
		{
			name: "no features",
			args: args{
				features: []Feature[string]{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid percentage",
			args: args{
				features: []Feature[string]{
					NewFeature("feature1", 50),
					NewFeature("feature2", -10), // Invalid percentage
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFeatureRoller(tt.args.features)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFeatureRoller() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFeatureRoller() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEnabled(t *testing.T) {
	type args struct {
		features []Feature[string]
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Testing isEnabled with valid features and userId",
			args: args{
				features: []Feature[string]{
					NewFeature("feature1", 50),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFeatureRoller(tt.args.features)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFeatureRoller() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.IsEnabled("feature1", 123) {
				t.Errorf("IsEnabled() = false, want true")
			}
			if got.IsEnabled("feature2", 123) {
				t.Errorf("IsEnabled() = true, want false for non-existent feature")
			}
			if got.IsEnabled("feature1", 456) {
				t.Errorf("IsEnabled() = true, want false for different userId")
			}
			if got.IsEnabled("feature1", 1077) {
				t.Errorf("IsEnabled() = true, want false for different userId")
			}
			if !got.IsEnabled("feature1", 1011) {
				t.Errorf("IsEnabled() = true, want false for different userId")
			}
			// testing distribution according to percentage without userid.
			enabledCount := 0
			for i := 0; i < 1000; i++ {
				if got.IsEnabled("feature1") {
					enabledCount++
				}
			}
			expectedPercentage := 50
			enabledPercentage := (enabledCount * 100) / 1000
			difference := (expectedPercentage - enabledPercentage)
			t.Logf("Total %d, EnabledCount %d, Enabled Percentage: %d%%, Expected Percentage: %d%%", 1000, enabledCount, enabledPercentage, expectedPercentage)
			if !(difference > -5 && difference < 5) { // Allowing a 5% margin of error
				t.Errorf("IsEnabled() = %d%%, want around %d%%", enabledPercentage, expectedPercentage)
			}
		})
	}
}
