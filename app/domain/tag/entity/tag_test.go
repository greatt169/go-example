package entity

import (
	"testing"
)

func TestTag_TableName(t1 *testing.T) {
	tests := []struct {
		name   string
		fields Tag
		want   string
	}{
		{name: "table name is correct", want: "tag"},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Tag{
				Id:   tt.fields.Id,
				Name: tt.fields.Name,
			}
			if got := t.TableName(); got != tt.want {
				t1.Errorf("TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
