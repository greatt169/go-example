package entity

import (
	"testing"
)

func TestNews_TableName(t *testing.T) {
	tests := []struct {
		name   string
		fields News
		want   string
	}{
		{name: "table name is correct", want: "news"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := News{}
			if got := n.TableName(); got != tt.want {
				t.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
