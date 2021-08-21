package entity

import (
	"testing"
)

func TestFile_TableName(t *testing.T) {
	tests := []struct {
		name   string
		fields File
		want   string
	}{
		{name: "table name is correct", want: "files"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := File{}
			if got := f.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
