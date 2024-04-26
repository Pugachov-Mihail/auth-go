package auth_validate

import (
	"testing"
)

func TestValidateLenValuesString(t *testing.T) {
	type args struct {
		value  string
		maxLen int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Base",
			args: args{
				value:  "dasdasda",
				maxLen: 10,
			},
			want: true,
		},
		{
			name: "Base2",
			args: args{
				value:  "dasdasda",
				maxLen: 5,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLenValuesString(tt.args.value, tt.args.maxLen); got != tt.want {
				t.Errorf("ValidateLenValuesString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "case 1",
			args:    args{email: "asdasd@mail.ru"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "case 2",
			args:    args{email: "asdasd@mailru"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "case 3",
			args:    args{email: "asdasdmail.ru"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "case 4",
			args:    args{email: "asdasd"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "case 5",
			args:    args{email: "Test@test.hui"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "case 6",
			args:    args{email: "¢[¥}¢£✓@£¢°€¥^.ru"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "case 7",
			args:    args{email: "test@test.ru@test.ru"},
			want:    false,
			wantErr: false,
		},
		{
			name:    "case 8",
			args:    args{email: "te_st@test.ru@test.ru"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
