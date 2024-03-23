package config_test

import (
	"fmt"
	. "go-template/internal/config"
	"os"
	"testing"

	. "github.com/agiledragon/gomonkey/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "key_arg",
			},
			want:    "",
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "key_arg",
			},
			want:    "value",
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.success {
				os.Setenv(tt.args.key, tt.want)
			}

			if got := GetString(tt.args.key); got != tt.want {
				t.Errorf("GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "int_arg",
			},
			want:    0,
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "int_arg",
			},
			want:    100,
			success: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.success {
				os.Setenv(tt.args.key, fmt.Sprintf("%d", tt.want))
			}
			if got := GetInt(tt.args.key); got != tt.want {
				t.Errorf("GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBool(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		success bool
	}{
		{
			name: "Failed to fetch value from env var",
			args: args{
				key: "bool_arg",
			},
			want:    false,
			success: false,
		},
		{
			name: "Successfully getting env var",
			args: args{
				key: "bool_arg",
			},
			want:    true,
			success: true,
		},
	}
	for _, tt := range tests {
		if tt.success {
			os.Setenv(tt.args.key, fmt.Sprintf("%v", tt.want))
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBool(tt.args.key); got != tt.want {
				t.Errorf("GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileName(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Successfully get .env.local",
			want: ".env.local",
		},
		{
			name: "nil env",
			want: ".env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "nil env" {
				os.Setenv("ENVIRONMENT_NAME", "")
			} else {
				os.Setenv("ENVIRONMENT_NAME", "local")
			}

			if got := FileName(); got != tt.want {
				t.Errorf("FileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

type args struct {
	env      string
	err      string
	tapped   bool
	dbSecret string
}

func TestLoadEnv(t *testing.T) {
	tests := getTestCases()

	for _, tt := range tests {
		mockEnvLoad(tt)
		setEnvironmentVariables(tt.args)

		t.Run(tt.name, func(t *testing.T) {
			testLoadEnv(t, tt)
		})
	}
}

func getTestCases() []struct {
	name    string
	wantErr bool
	args    args
} {
	return []struct {
		name    string
		wantErr bool
		args    args
	}{
		{
			name:    "Successfully load local env if ENVIRONMENT_NAME doesn't have a value",
			wantErr: false,
			args:    args{env: "", tapped: false},
		},
		// Add more test cases here...
	}
}

func mockEnvLoad(tt struct {
	name    string
	wantErr bool
	args    args
}) {
	ApplyFunc(godotenv.Load, func(filenames ...string) (err error) {
		tt.args.tapped = true
		if tt.args.err == "" {
			if tt.name == "Env varInjection Error" && len(filenames) > 0 && filenames[0] == ".env.local" {
				return fmt.Errorf(tt.args.err)
			}
			return nil
		}
		return fmt.Errorf(tt.args.err)
	})
}

func setEnvironmentVariables(args args) {
	os.Setenv("ENVIRONMENT_NAME", args.env)
	if args.dbSecret != "" {
		os.Setenv("DB_SECRET", args.dbSecret)
	}
}

func testLoadEnv(t *testing.T, tt struct {
	name    string
	wantErr bool
	args    args
}) {
	username := "go_template_role"
	host := "localhost"
	dbname := "go_template"
	password := "go_template_role456"
	port := "5432"
	if tt.name == "dbCredsInjected True" {
		ApplyFunc(GetBool, func(key string) bool {
			return true
		})
	}

	tapped := tt.args.tapped

	if err := LoadEnv(); (err != nil) != tt.wantErr {
		t.Errorf("LoadEnv() error = %v, wantErr %v", err, tt.wantErr)
	}
	assert.Equal(t, tapped, !tt.args.tapped)
	if tt.args.dbSecret != "" {
		assert.Equal(t, os.Getenv("PSQL_USER"), username)
		assert.Equal(t, os.Getenv("PSQL_HOST"), host)
		assert.Equal(t, os.Getenv("PSQL_DBNAME"), dbname)
		assert.Equal(t, os.Getenv("PSQL_PASS"), password)
		assert.Equal(t, os.Getenv("PSQL_PORT"), port)
	}
}
