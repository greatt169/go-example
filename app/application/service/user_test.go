package service

import (
	helpers "github.com/AeroAgency/golang-helpers-lib"
	"testing"
)

func TestNewUser(t *testing.T) {
	NewUser()
}

func TestUser_GetNameByToken(t *testing.T) {
	type fields struct {
		errorFormatter *helpers.ErrorFormatter
		jwt            *helpers.Jwt
		meta           *helpers.Meta
	}
	type args struct {
		tokenString string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Correct token [Case Аркадьев Аркадий Аркадьевич -> Аркадьев А. А.]",
			fields: fields{
				errorFormatter: &helpers.ErrorFormatter{},
				jwt:            &helpers.Jwt{},
				meta:           &helpers.Meta{},
			},
			args: args{
				tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQUk1fSUQiOlsiMTExIiwiMjIyIiwiMzMzIl0sImFjY2Vzc191dWlkIjoiOTRiNjA3NmItN2U1Ny00MjkzLTg5NzMtMWE3YzI2NDRhODBiIiwiYnJhbmNoIjoiYnJhbmNoLXNsdWciLCJleHAiOjE1OTkxNDM5NTAsImZhbWlseV9uYW1lIjoi0JDRgNC60LDQtNGM0LXQsiIsImZpbyI6ItCQ0YDQutCw0LTRjNC10LIg0JDRgNC60LDQtNC40Lkg0JDRgNC60LDQtNGM0LXQstC40YciLCJnaXZlbl9uYW1lIjoi0JDRgNC60LDQtNC40LkiLCJtaWRkbGVfbmFtZSI6ItCQ0YDQutCw0LTRjNC10LLQuNGHIiwicG9zaXRpb24iOiLQkdC40LfQvdC10YEg0LDQtNC80LjQvdC40YHRgtGA0LDRgtC-0YAg0JzQpCIsInJlc291cmNlX2FjY2VzcyI6eyJkZWFsZXItcG9ydGFsIjp7InJvbGVzIjpbIm1lZ2Fmb25fYnVpc25lc3NfYWRtaW4iXX19LCJzdWIiOiJkYzc5NTQxYi1hODU0LTRjOWQtYTQyZi1lYzA5ZTBlMzY4ODcifQ.kRTF23_M60hicAo3pa-0sSZHgUPEB6yOwFx-1ap8Ezo",
			},
			want: "Аркадьев А.А",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{
				errorFormatter: tt.fields.errorFormatter,
				jwt:            tt.fields.jwt,
				meta:           tt.fields.meta,
			}
			if got := u.GetNameByToken(tt.args.tokenString); got != tt.want {
				t.Errorf("GetNameByToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
