package service

import (
	"reflect"
	"testing"
)

func TestValidatorRulesMap_GetGetNewsMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"limit":  {"numeric", "min:0", "max:100"},
				"offset": {"numeric"},
				"sort":   {"in:active_from,date_create"},
				"order":  {"in:asc,desc"},
				"query":  {"min_string_len:3", "max_string_len:100"},
				"mode":   {"in:active,inactive"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetGetNewsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGetNewsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetCreateNewsMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"title":       {"required", "max_string_len:200"},
				"text":        {"required"},
				"textJson":    {"required"},
				"activeFrom":  {"required", "numeric"},
				"isImportant": {"bool"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetCreateNewsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCreateNewsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetCreateNewsTagMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"name": {"required", "max_string_len:50", "regex:[a-zа-я0-9]+$"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetCreateNewsTagMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCreateNewsTagMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetUpdateNewsTagMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"name": {"required", "max_string_len:50", "regex:[a-zа-я0-9]+$"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetUpdateNewsTagMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdateNewsTagMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetNewsFilesMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "Test answer",
			want: map[string]string{
				"file_ext":    "doc,xls,ppt,jpg,bmp,pdf,rtf,txt,zip",
				"file_size":   "5242880",
				"files_limit": "7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetNewsFilesMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNewsFilesMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetUpdateNewsMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"title":       {"required", "max_string_len:200"},
				"text":        {"required"},
				"textJson":    {"required"},
				"activeFrom":  {"required", "numeric"},
				"isImportant": {"bool"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetUpdateNewsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdateNewsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetGetPromoMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"limit":  {"numeric", "min:0", "max:100"},
				"sort":   {"in:active_from,date_create"},
				"order":  {"in:asc,desc"},
				"offset": {"numeric"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetGetPromoMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGetPromoMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetCreatePromoMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test answer",
			want: map[string][]string{
				"title":      {"required", "max_string_len:140"},
				"text":       {"required"},
				"textJson":   {"required"},
				"activeFrom": {"required", "numeric"},
				"period":     {"max_string_len:80"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetCreatePromoMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCreatePromoMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetUpdatePromoMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Correct Answer",
			want: map[string][]string{
				"title":      {"required", "max_string_len:140"},
				"text":       {"required"},
				"textJson":   {"required"},
				"activeFrom": {"required", "numeric"},
				"period":     {"max_string_len:80"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetUpdatePromoMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdatePromoMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetGetOneMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Correct Answer",
			want: map[string][]string{
				"id": {"uuid_v4"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetGetOneMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGetOneMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetDeleteMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Correct Answer",
			want: map[string][]string{
				"id": {"uuid_v4"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetDeleteMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDeleteMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetSubscribeMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Correct Answer",
			want: map[string][]string{
				"email": {"email"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetSubscribeMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSubscribeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatorRulesMap_GetGetOneBySlugMap(t *testing.T) {
	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Correct Answer",
			want: map[string][]string{
				"slug": {"required"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ValidatorRulesMap{}
			if got := v.GetGetOneBySlugMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGetOneBySlugMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
