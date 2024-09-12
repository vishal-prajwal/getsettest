package validate

import "testing"

func TestValidatePanNumber(t *testing.T) {
	type args struct {
		panNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid Pan Number",
			args:    args{panNumber: "ascpy1892J"},
			wantErr: false,
		},
		{
			name:    "Invalid Pan Number",
			args:    args{panNumber: "acpy1892J"},
			wantErr: true,
		},
		{
			name:    "Not a individual pan number",
			args:    args{panNumber: "ABCDE1892J"},
			wantErr: true,
		},
		{
			name:    "Not a valid pan number",
			args:    args{panNumber: "ABCDE1892"},
			wantErr: true,
		},
		{
			name:    "Not a valid pan number",
			args:    args{panNumber: " ABCDE1892"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePanNumber(tt.args.panNumber); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePanNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
