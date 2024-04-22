package util

import "testing"

func TestReplaceEnvVars(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		exist bool
	}{
		{
			name: "Take default value",
			args: args{"${MY_OWN_HOME_ENV:/home}"},
			want: "/home",
		},
		{
			name:  "Take environment variable value",
			args:  args{"${PATH:/path}"},
			exist: true,
		},
		{
			name: "Take common variable value",
			args: args{"/home"},
			want: "/home",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceEnvVars(tt.args.value); (tt.exist == true && got == "") ||
				(tt.exist == false && got != tt.want) {
				t.Errorf("ReplaceEnvVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
