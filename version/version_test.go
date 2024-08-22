package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple",
			value: t.Name(),
			want:  t.Name(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVersion := version
			version = tt.value
			t.Cleanup(func() {
				version = oldVersion
			})

			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "simple",
			value: t.Name(),
			want:  t.Name(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldCommit := commit
			commit = tt.value
			t.Cleanup(func() {
				commit = oldCommit
			})

			if got := Commit(); got != tt.want {
				t.Errorf("Commit() = %v, want %v", got, tt.want)
			}
		})
	}
}
