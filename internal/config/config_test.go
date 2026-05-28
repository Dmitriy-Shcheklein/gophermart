package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig_FlagsOverrideEnv(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() {
		flag.CommandLine = oldFlagSet
		_ = os.Unsetenv("RUN_ADDRESS")
	}()

	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	_ = os.Setenv("RUN_ADDRESS", "localhost:9090")

	srvAddr := NewSrvAddr()
	srvAddr.ApplyEnv()

	result := srvAddr.String()
	assert.Equal(t, "localhost:9090", result)
}

func TestNewConfig_ValidEnvVars(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() {
		flag.CommandLine = oldFlagSet
		_ = os.Unsetenv("DATABASE_URI")
		_ = os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
	}()

	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	_ = os.Setenv("DATABASE_URI", "postgres://localhost:5432/testdb")
	_ = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")

	dsn := NewDSN()
	dsn.ApplyEnv()
	accrual := NewAccrual()
	accrual.ApplyEnv()

	assert.Equal(t, "postgres://localhost:5432/testdb", dsn.value)
	assert.Equal(t, "http://localhost:8081", accrual.value)
}

func TestSrvAddr_SetValidAddresses(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
		wantErr bool
	}{
		{"Valid address", "localhost:8080", "localhost:8080", false},
		{"IP address", "192.168.1.1:9090", "192.168.1.1:9090", false},
		{"Different port", "localhost:3000", "localhost:3000", false},
	}

	oldFlagSet := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagSet }()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	srvAddr := NewSrvAddr()

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := srvAddr.Set(tt.address)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					assert.Equal(t, tt.want, srvAddr.String())
				}
			},
		)
	}
}

func TestSrvAddr_SetInvalidFormats(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{"Missing port", "localhost", true},
		{"No colon separator", "localhost8080", true},
		{"Multiple colons", "host::8080", true},
		{"Invalid port", "localhost:abc", true},
		{"Empty address", "", true},
		{"Just colon", ":", true},
	}

	oldFlagSet := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagSet }()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
	srvAddr := NewSrvAddr()

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				err := srvAddr.Set(tt.address)

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			},
		)
	}
}

func TestSrvAddr_DefaultValues(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagSet }()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	srvAddr := NewSrvAddr()

	assert.Equal(t, "localhost:8080", srvAddr.String())
}

func TestApplyEnv_EmptyEnvVars(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagSet }()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	srvAddr := NewSrvAddr()
	srvAddr.ApplyEnv()

	assert.Equal(t, "localhost:8080", srvAddr.String())
}

func TestDSN_ApplyEnvValidDSNs(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			"PostgreSQL standard", "postgres://user:pass@localhost:5432/dbname",
			"postgres://user:pass@localhost:5432/dbname",
		},
		{"With host", "postgres://localhost:5432/testdb", "postgres://localhost:5432/testdb"},
		{
			"With all parameters", "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			"postgres://user:pass@localhost:5432/dbname?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				oldFlagSet := flag.CommandLine
				defer func() { flag.CommandLine = oldFlagSet }()
				flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

				_ = os.Setenv("DATABASE_URI", tt.dsn)
				defer func() {
					_ = os.Unsetenv("DATABASE_URI")
				}()

				dsn := NewDSN()
				dsn.ApplyEnv()

				assert.Equal(t, tt.want, dsn.value)
			},
		)
	}
}

func TestAccrual_ApplyEnvValidAddresses(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    string
	}{
		{"Valid HTTP", "http://localhost:8080", "http://localhost:8080"},
		{"Valid HTTPS", "https://localhost:8443", "https://localhost:8443"},
		{"Full URL with path", "http://example.com/api/orders", "http://example.com/api/orders"},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				oldFlagSet := flag.CommandLine
				defer func() { flag.CommandLine = oldFlagSet }()
				flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

				_ = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", tt.address)
				defer func() {
					_ = os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
				}()

				accrual := NewAccrual()
				accrual.ApplyEnv()

				assert.Equal(t, tt.want, accrual.value)
			},
		)
	}
}

func TestConfig_GetterMethods(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() {
		flag.CommandLine = oldFlagSet
		_ = os.Unsetenv("DATABASE_URI")
		_ = os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
	}()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	_ = os.Setenv("DATABASE_URI", "postgres://localhost:5432/testdb")
	_ = os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "http://localhost:8081")

	config, err := New()

	require.NoError(t, err)
	require.NotNil(t, config)

	srvAddr := config.GetSrvAddr()
	require.NotEmpty(t, srvAddr)

	dbDsn := config.DbDsn()
	require.Equal(t, "postgres://localhost:5432/testdb", dbDsn)

	accrualAddr := config.GetAccrualSrvAddr()
	require.Equal(t, "http://localhost:8081", accrualAddr)
}

func TestConfig_MissingRequiredConfig(t *testing.T) {
	oldFlagSet := flag.CommandLine
	defer func() { flag.CommandLine = oldFlagSet }()
	flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)

	config, err := New()

	require.Error(t, err)
	assert.Nil(t, config)
}
