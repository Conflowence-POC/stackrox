package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/image/policies"
	"bitbucket.org/stack-rox/apollo/pkg/clientconn"
	"bitbucket.org/stack-rox/apollo/pkg/defaults"
	"bitbucket.org/stack-rox/apollo/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	policy = &v1.Policy{
		Name:        "test policy " + fmt.Sprintf("%d", time.Now().UnixNano()),
		Description: "description",
		Severity:    v1.Severity_HIGH_SEVERITY,
		Categories:  []string{"Image Assurance", "Privileges Capabilities"},
		Disabled:    false,
		ImagePolicy: &v1.ImagePolicy{
			ImageName: &v1.ImageNamePolicy{
				Tag: "latest",
			},
		},
		PrivilegePolicy: &v1.PrivilegePolicy{
			SetPrivileged: &v1.PrivilegePolicy_Privileged{
				Privileged: true,
			},
		},
	}

	logger = logging.LoggerForModule()
)

func TestDefaultPolicies(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := clientconn.UnauthenticatedGRPCConnection(apiEndpoint)
	require.NoError(t, err)

	defaults.PoliciesPath = policies.Directory()
	defaultPolicies, err := defaults.Policies()
	require.NoError(t, err)

	service := v1.NewPolicyServiceClient(conn)
	resp, err := service.GetPolicies(ctx, &v1.RawQuery{})
	require.NoError(t, err)

	policiesMap := make(map[string]*v1.Policy)
	for _, p := range resp.GetPolicies() {
		p.Id = ""
		policiesMap[p.GetName()] = p
	}

	assert.Equal(t, len(defaultPolicies), len(resp.GetPolicies()))

	for _, p := range defaultPolicies {
		assert.Equal(t, p, policiesMap[p.GetName()])
	}
}

func TestPoliciesCRUD(t *testing.T) {
	conn, err := clientconn.UnauthenticatedGRPCConnection(apiEndpoint)
	require.NoError(t, err)

	service := v1.NewPolicyServiceClient(conn)

	subtests := []struct {
		name string
		test func(t *testing.T, service v1.PolicyServiceClient)
	}{
		{
			name: "create",
			test: verifyCreatePolicy,
		},
		{
			name: "read",
			test: verifyReadPolicy,
		},
		{
			name: "update",
			test: verifyUpdatePolicy,
		},
		{
			name: "delete",
			test: verifyDeletePolicy,
		},
	}

	for _, sub := range subtests {
		t.Run(sub.name, func(t *testing.T) {
			sub.test(t, service)
		})
	}
}

func verifyCreatePolicy(t *testing.T, service v1.PolicyServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	postResp, err := service.PostPolicy(ctx, policy)
	require.NoError(t, err)

	policy.Id = postResp.GetId()
	assert.Equal(t, policy, postResp)
}

func verifyReadPolicy(t *testing.T, service v1.PolicyServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	getResp, err := service.GetPolicy(ctx, &v1.ResourceByID{Id: policy.GetId()})
	cancel()
	require.NoError(t, err)
	assert.Equal(t, policy, getResp)

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	getManyResp, err := service.GetPolicies(ctx, &v1.RawQuery{
		Query: getPolicyQuery(policy.GetName()),
	})
	cancel()
	require.NoError(t, err)
	assert.Equal(t, 1, len(getManyResp.GetPolicies()))
	if len(getManyResp.GetPolicies()) > 0 {
		assert.Equal(t, policy, getManyResp.GetPolicies()[0])
	}
}

func verifyUpdatePolicy(t *testing.T, service v1.PolicyServiceClient) {
	policy.Severity = v1.Severity_LOW_SEVERITY
	policy.Description = "updated description"
	policy.Disabled = true
	policy.ImagePolicy.SetScanAgeDays = &v1.ImagePolicy_ScanAgeDays{ScanAgeDays: 10}
	policy.PrivilegePolicy.AddCapabilities = []string{"CAP_SYS_MODULE"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := service.PutPolicy(ctx, policy)
	cancel()
	require.NoError(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	getResp, err := service.GetPolicy(ctx, &v1.ResourceByID{Id: policy.GetId()})
	cancel()
	require.NoError(t, err)
	assert.Equal(t, policy, getResp)
}

func verifyDeletePolicy(t *testing.T, service v1.PolicyServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	_, err := service.DeletePolicy(ctx, &v1.ResourceByID{Id: policy.GetId()})
	cancel()

	require.NoError(t, err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	_, err = service.GetPolicy(ctx, &v1.ResourceByID{Id: policy.GetId()})
	cancel()
	s, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, s.Code())
}
