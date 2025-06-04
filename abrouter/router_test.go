package abrouter

import (
	"reflect"
	"testing"
)

func TestNewABRouter(t *testing.T) {
	type args struct {
		choices []Choice[string]
	}
	tests := []struct {
		name    string
		args    args
		want    ABRouter[string]
		wantErr bool
	}{
		{
			name: "valid choices",
			args: args{
				choices: []Choice[string]{
					NewChoice("A", 50),
					NewChoice("B", 50),
				},
			},
			want: &abRouterImpl[string]{
				choices: []Choice[string]{
					NewChoice("A", 50),
					NewChoice("B", 50),
				},
				totalWeight: 100,
			},
			wantErr: false,
		},
		{
			name: "no choices",
			args: args{
				choices: []Choice[string]{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid weight",
			args: args{
				choices: []Choice[string]{
					NewChoice("A", 50),

					NewChoice("B", -10), // Invalid weight
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid choices with zero weight",
			args: args{
				choices: []Choice[string]{
					NewChoice("A", 0),
					NewChoice("B", 100),
				},
			},
			want: &abRouterImpl[string]{
				choices: []Choice[string]{
					NewChoice("A", 0),
					NewChoice("B", 100),
				},
				totalWeight: 100,
			},
			wantErr: false,
		},
		{
			name: "total weights is not zero",
			args: args{
				choices: []Choice[string]{
					NewChoice("A", 0),
					NewChoice("B", 0),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid choices with different weights",
			args: args{
				choices: []Choice[string]{
					NewChoice("A", 30),
					NewChoice("B", 70),
				},
			},
			want: &abRouterImpl[string]{
				choices: []Choice[string]{
					NewChoice("A", 30),
					NewChoice("B", 70),
				},
				totalWeight: 100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewABRouter(tt.args.choices)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewABRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewABRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChoose(t *testing.T) {
	type args struct {
		choices []Choice[string]
	}
	tests := []struct {
		name string
	}{
		{
			name: "valid choices",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewABRouter([]Choice[string]{
				NewChoice("A", 50),
				NewChoice("B", 50),
			})
			if err != nil {
				t.Fatalf("NewABRouter() error = %v", err)
			}
			if choice := got.Choose(1); choice != "A" {
				t.Errorf("Choose() = %v, want 'A'", choice)
			}

			if choice := got.Choose(49); choice != "A" {
				t.Errorf("Choose() = %v, want 'A'", choice)
			}
			if choice := got.Choose(50); choice != "B" {
				t.Errorf("Choose() = %v, want 'B'", choice)
			}
			if choice := got.Choose(99); choice != "B" {
				t.Errorf("Choose() = %v, want 'B'", choice)
			}
			if choice := got.Choose(125); choice != "A" {
				t.Errorf("Choose() = %v, want 'A'", choice)
			}
			if choice := got.Choose(158); choice != "B" {
				t.Errorf("Choose() = %v, want 'B'", choice)
			}

			// Test with no userId and check the distribution
			choicesCount := map[string]int{
				"A": 0,
				"B": 0,
			}
			for i := 0; i < 1000; i++ {
				choice := got.Choose()
				choicesCount[choice]++
			}
			choiceAPercentage := (choicesCount["A"] * 100) / 1000
			choiceBPercentage := (choicesCount["B"] * 100) / 1000
			difference := choiceAPercentage - choiceBPercentage
			t.Logf("Count of A=%d and B=%d, Distribution of choices: A = %d%%, B = %d%%, difference = %d%%", choicesCount["A"], choicesCount["B"], choiceAPercentage, choiceBPercentage, difference)
			if difference < -5 || difference > 5 {
				t.Errorf("Distribution of choices is not balanced: A = %d%%, B = %d%%, difference = %d%%", choiceAPercentage, choiceBPercentage, difference)
			}
		})
	}
}
