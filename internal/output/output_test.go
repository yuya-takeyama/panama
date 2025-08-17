package output

import (
	"testing"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Format
		wantErr bool
	}{
		{
			name:    "path format",
			input:   "path",
			want:    FormatPath,
			wantErr: false,
		},
		{
			name:    "cd format",
			input:   "cd",
			want:    FormatCD,
			wantErr: false,
		},
		{
			name:    "json format",
			input:   "json",
			want:    FormatJSON,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}