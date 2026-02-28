package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPermissionsBuilder(t *testing.T) {
	pb := NewPermissionsBuilder()

	assert.NotNil(t, pb.Permissions)
	assert.Empty(t, pb.Permissions)
}

func TestPermissionsBuilder_Allow(t *testing.T) {
	tests := []struct {
		name      string
		resource  Resource
		operation Operation
	}{
		{"allow APP CREATE", APP, CREATE},
		{"allow KEY READ", KEY, READ},
		{"allow DELIVERY UPDATE", DELIVERY, UPDATE},
		{"allow OneTimeToken DESTROY", OneTimeToken, DESTROY},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := NewPermissionsBuilder()

			result := pb.Allow(tt.resource, tt.operation)

			assert.Equal(t, &pb, result)
			assert.Contains(t, pb.Permissions[tt.resource], tt.operation)
		})
	}
}

func TestPermissionsBuilder_AllowChaining(t *testing.T) {
	pb := NewPermissionsBuilder()

	pb.Allow(APP, CREATE).
		Allow(APP, READ).
		Allow(KEY, CREATE)

	assert.Contains(t, pb.Permissions[APP], CREATE)
	assert.Contains(t, pb.Permissions[APP], READ)
	assert.Contains(t, pb.Permissions[KEY], CREATE)
}

func TestPermissionsBuilder_AllowNoDuplicates(t *testing.T) {
	pb := NewPermissionsBuilder()

	pb.Allow(APP, CREATE)
	pb.Allow(APP, CREATE)
	pb.Allow(APP, CREATE)

	count := 0
	for _, op := range pb.Permissions[APP] {
		if op == CREATE {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestPermissionsBuilder_Deny(t *testing.T) {
	pb := NewPermissionsBuilder()

	pb.Allow(APP, CREATE)
	pb.Allow(APP, READ)

	result := pb.Deny(APP, CREATE)

	assert.Equal(t, &pb, result)
	assert.NotContains(t, pb.Permissions[APP], CREATE)
	assert.Contains(t, pb.Permissions[APP], READ)
}

func TestPermissionsBuilder_DenyNonExistent(t *testing.T) {
	pb := NewPermissionsBuilder()

	result := pb.Deny(APP, CREATE)

	assert.Equal(t, &pb, result)
	assert.Empty(t, pb.Permissions[APP])
}

func TestPermissionsBuilder_Can(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() PermissionsBuilder
		resource   Resource
		operation  Operation
		wantResult bool
	}{
		{
			name: "can with direct permission",
			setup: func() PermissionsBuilder {
				pb := NewPermissionsBuilder()
				pb.Allow(APP, CREATE)
				return pb
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: true,
		},
		{
			name: "cannot without permission",
			setup: func() PermissionsBuilder {
				pb := NewPermissionsBuilder()
				pb.Allow(APP, READ)
				return pb
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: false,
		},
		{
			name: "can with RESOURCES_ALL wildcard",
			setup: func() PermissionsBuilder {
				pb := NewPermissionsBuilder()
				pb.Allow(RESOURCES_ALL, OPERATIONS_ALL)
				return pb
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: true,
		},
		{
			name: "can with RESOURCES_ALL and specific operation",
			setup: func() PermissionsBuilder {
				pb := NewPermissionsBuilder()
				pb.Allow(RESOURCES_ALL, CREATE)
				return pb
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: true,
		},
		{
			name: "cannot with RESOURCES_ALL but different operation",
			setup: func() PermissionsBuilder {
				pb := NewPermissionsBuilder()
				pb.Allow(RESOURCES_ALL, READ)
				return pb
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: false,
		},
		{
			name: "cannot with no permissions at all",
			setup: func() PermissionsBuilder {
				return NewPermissionsBuilder()
			},
			resource:   APP,
			operation:  CREATE,
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := tt.setup()
			result := pb.Can(tt.resource, tt.operation)
			assert.Equal(t, tt.wantResult, result)
		})
	}
}

func TestPermissionsBuilder_String(t *testing.T) {
	pb := NewPermissionsBuilder()
	pb.Allow(APP, CREATE)
	pb.Allow(APP, READ)

	str, err := pb.String()

	require.NoError(t, err)
	assert.NotEmpty(t, str)

	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(str), &parsed)
	require.NoError(t, err)
}

func TestParsePermissions(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantErr   bool
		checkFunc func(t *testing.T, pb *PermissionsBuilder)
	}{
		{
			name:    "valid permissions JSON",
			input:   `{"permissions":{"APP":["CREATE","READ"]}}`,
			wantErr: false,
			checkFunc: func(t *testing.T, pb *PermissionsBuilder) {
				assert.True(t, pb.Can(APP, CREATE))
				assert.True(t, pb.Can(APP, READ))
				assert.False(t, pb.Can(APP, DESTROY))
			},
		},
		{
			name:    "empty permissions",
			input:   `{"permissions":{}}`,
			wantErr: false,
			checkFunc: func(t *testing.T, pb *PermissionsBuilder) {
				assert.False(t, pb.Can(APP, CREATE))
			},
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			wantErr: true,
			checkFunc: func(t *testing.T, pb *PermissionsBuilder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb, err := ParsePermissions([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, pb)
			tt.checkFunc(t, pb)
		})
	}
}

func TestResourceConstants(t *testing.T) {
	assert.Equal(t, Resource("APP"), APP)
	assert.Equal(t, Resource("KEY"), KEY)
	assert.Equal(t, Resource("ONE_TIME_TOKEN"), OneTimeToken)
	assert.Equal(t, Resource("DELIVERY"), DELIVERY)
	assert.Equal(t, Resource("RESOURCES_ALL"), RESOURCES_ALL)
}

func TestOperationConstants(t *testing.T) {
	assert.Equal(t, Operation("CREATE"), CREATE)
	assert.Equal(t, Operation("READ"), READ)
	assert.Equal(t, Operation("DESTROY"), DESTROY)
	assert.Equal(t, Operation("UPDATE"), UPDATE)
	assert.Equal(t, Operation("OPERATIONS_ALL"), OPERATIONS_ALL)
}
