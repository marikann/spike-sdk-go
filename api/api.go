//    \\ SPIKE: Secure your secrets with SPIFFE.
//  \\\\\ Copyright 2024-present SPIKE contributors.
// \\\\\\\ SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"

	"github.com/spiffe/go-spiffe/v2/workloadapi"

	"github.com/spiffe/spike-sdk-go/api/entity/data"
	"github.com/spiffe/spike-sdk-go/api/internal/impl/api/acl"
	"github.com/spiffe/spike-sdk-go/api/internal/impl/api/operator"
	"github.com/spiffe/spike-sdk-go/api/internal/impl/api/secret"
	"github.com/spiffe/spike-sdk-go/spiffe"
)

// Api is the SPIKE API.
type Api struct {
	source *workloadapi.X509Source
}

// New creates and returns a new instance of Api configured with a SPIFFE source.
func New() *Api {
	defaultEndpointSocket := spiffe.EndpointSocket()

	source, _, err := spiffe.Source(context.Background(), defaultEndpointSocket)
	if err != nil {
		return nil
	}

	return &Api{source: source}
}

// NewWithSource initializes a new Api instance with the given X509Source.
func NewWithSource(source *workloadapi.X509Source) *Api {
	return &Api{source: source}
}

// Close releases any resources held by the Api instance.
// It ensures proper cleanup of the underlying source.
func (a *Api) Close() {
	spiffe.CloseSource(a.source)
}

// CreatePolicy creates a new policy in the system. It establishes a mutual
// TLS connection using the X.509 source and sends a policy creation request
// to the server.
//
// The function takes the following parameters:
//   - name string: The name of the policy to be created
//   - spiffeIdPattern string: The SPIFFE ID pattern that this policy will apply
//     to
//   - pathPattern string: The path pattern that this policy will match against
//   - permissions []data.PolicyPermission: A slice of PolicyPermission defining
//     the access rights for this policy
//
// The function returns an error if any of the following operations fail:
//   - Marshaling the policy creation request
//   - Creating the mTLS client
//   - Making the HTTP POST request
//   - Unmarshaling the response
//   - Server-side policy creation (indicated in the response)
//
// Example usage:
//
//	permissions := []data.PolicyPermission{
//	    {
//	        Action: "read",
//	        Resource: "documents/*",
//	    },
//	}
//
//	err = api.CreatePolicy(
//	    "doc-reader",
//	    "spiffe://example.org/service/*",
//	    "/api/documents/*",
//	    permissions,
//	)
//	if err != nil {
//	    log.Printf("Failed to create policy: %v", err)
//	    return
//	}
func (a *Api) CreatePolicy(
	name string, spiffeIdPattern string, pathPattern string,
	permissions []data.PolicyPermission,
) error {
	return acl.CreatePolicy(a.source,
		name, spiffeIdPattern, pathPattern, permissions)
}

// DeletePolicy removes an existing policy from the system using its name.
//
// The function takes the following parameters:
//   - name string: The name of the policy to be deleted
//
// The function returns an error if any of the following operations fail:
//   - Marshaling the policy deletion request
//   - Creating the mTLS client
//   - Making the HTTP POST request
//   - Unmarshaling the response
//   - Server-side policy deletion (indicated in the response)
//
// Example usage:
//
//	err = api.DeletePolicy("doc-reader")
//	if err != nil {
//	    log.Printf("Failed to delete policy: %v", err)
//	    return
//	}
func (a *Api) DeletePolicy(name string) error {
	return acl.DeletePolicy(a.source, name)
}

// GetPolicy retrieves a policy from the system using its name.
//
// The function takes the following parameters:
//   - name string: The name of the policy to retrieve
//
// The function returns:
//   - (*data.Policy, nil) if the policy is found
//   - (nil, nil) if the policy is not found
//   - (nil, error) if an error occurs during the operation
//
// Errors can occur during:
//   - Marshaling the policy retrieval request
//   - Creating the mTLS client
//   - Making the HTTP POST request (except for not found cases)
//   - Unmarshaling the response
//   - Server-side policy retrieval (indicated in the response)
//
// Example usage:
//
//	policy, err := api.GetPolicy("doc-reader")
//	if err != nil {
//	    log.Printf("Error retrieving policy: %v", err)
//	    return
//	}
//	if policy == nil {
//	    log.Printf("Policy not found")
//	    return
//	}
//
//	log.Printf("Found policy: %+v", policy)
func (a *Api) GetPolicy(name string) (*data.Policy, error) {
	return acl.GetPolicy(a.source, name)
}

// ListPolicies retrieves all policies from the system.
//
// The function returns:
//   - (*[]data.Policy, nil) containing all policies if successful
//   - (nil, nil) if no policies are found
//   - (nil, error) if an error occurs during the operation
//
// Note: The returned slice pointer should be dereferenced before use:
//
//	policies := *result
//
// Errors can occur during:
//   - Marshaling the policy list request
//   - Creating the mTLS client
//   - Making the HTTP POST request (except for not found cases)
//   - Unmarshaling the response
//   - Server-side policy listing (indicated in the response)
//
// Example usage:
//
//	result, err := api.ListPolicies()
//	if err != nil {
//	    log.Printf("Error listing policies: %v", err)
//	    return
//	}
//	if result == nil {
//	    log.Printf("No policies found")
//	    return
//	}
//
//	policies := *result
//	for _, policy := range policies {
//	    log.Printf("Found policy: %+v", policy)
//	}
func (a *Api) ListPolicies() (*[]data.Policy, error) {
	return acl.ListPolicies(a.source)
}

// DeleteSecretVersions deletes specified versions of a secret at the given
// path
//
// It constructs a delete request and sends it to the secrets API endpoint.
//
// Parameters:
//   - path string: Path to the secret to delete
//   - versions []int: Array of version numbers to delete
//
// Returns:
//   - error: nil on success, unauthorized error if not logged in, or wrapped
//     error on request/parsing failure
//
// Example:
//
//	err := api.DeleteSecretVersions("secret/path", []int{1, 2})
func (a *Api) DeleteSecretVersions(path string, versions []int) error {
	return secret.Delete(a.source, path, versions)
}

// DeleteSecret deletes the entire secret at the given path
//
// Parameters:
//   - path string: Path to the secret to delete
//
// Returns:
//   - error: nil on success, unauthorized error if not logged in, or wrapped
//     error on request/parsing failure
//
// Example:
//
//	err := api.DeleteSecret("secret/path")
func (a *Api) DeleteSecret(path string) error {
	return secret.Delete(a.source, path, []int{})
}

// GetSecretVersion retrieves a specific version of a secret at the given
// path.
//
// Parameters:
//   - path string: Path to the secret to retrieve
//   - version int: Version number of the secret to retrieve
//
// Returns:
//   - *data.Secret: Secret data if found, nil if secret not found
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	secret, err := api.GetSecretVersion("secret/path", 1)
func (a *Api) GetSecretVersion(
	path string, version int,
) (*data.Secret, error) {
	return secret.Get(a.source, path, version)
}

// GetSecret retrieves the latest version of the secret at the given path.
//
// Parameters:
//   - path string: Path to the secret to retrieve
//
// Returns:
//   - *data.Secret: Secret data if found, nil if secret not found
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	secret, err := api.GetSecret("secret/path")
func (a *Api) GetSecret(path string) (*data.Secret, error) {
	return secret.Get(a.source, path, 0)
}

// ListSecretKeys retrieves all secret keys.
//
// Returns:
//   - *[]string: Pointer to array of secret keys if found, nil if none found
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	keys, err := api.ListSecretKeys()
func (a *Api) ListSecretKeys() (*[]string, error) {
	return secret.ListKeys(a.source)
}

// GetSecretMetadata retrieves metadata for a specific version of a secret at
// the given path.
//
// Parameters:
//   - path string: Path to the secret to retrieve metadata for
//   - version int: Version number of the secret to retrieve metadata for
//
// Returns:
//   - *data.SecretMetadata: Secret metadata if found, nil if secret not found
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	metadata, err := api.GetSecretMetadata("secret/path", 1)
func (a *Api) GetSecretMetadata(
	path string, version int,
) (*data.SecretMetadata, error) {
	return secret.GetMetadata(a.source, path, version)
}

// PutSecret creates or updates a secret at the specified path with the given
// values.
//
// Parameters:
//   - path string: Path where the secret should be stored
//   - data map[string]string: Map of key-value pairs representing the secret
//     data
//
// Returns:
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	err := api.PutSecret("secret/path", map[string]string{"key": "value"})
func (a *Api) PutSecret(path string, data map[string]string) error {
	return secret.Put(a.source, path, data)
}

// UndeleteSecret restores previously deleted versions of a secret at the
// specified path.
//
// Parameters:
//   - path string: Path to the secret to restore
//   - versions []int: Array of version numbers to restore. Empty array
//     attempts no restoration
//
// Returns:
//   - error: nil on success, unauthorized error if not logged in, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	err := api.UndeleteSecret("secret/path", []int{1, 2})
func (a *Api) UndeleteSecret(path string, versions []int) error {
	return secret.Undelete(a.source, path, versions)
}

// Recover returns recovery partitions for SPIKE Nexus to be used in a
// break-the-glass recovery operation if SPIKE Nexus auto-recovery mechanism
// isn't successful.
//
// The returned shards are sensitive and should be securely stored out-of-band
// in encrypted form.
//
// Returns:
//   - *[][32]byte: Pointer to array of recovery shards as 32-byte arrays
//   - error: nil on success, unauthorized error if not authorized, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	shards, err := api.Recover()
func (a *Api) Recover() (map[int]*[32]byte, error) {
	return operator.Recover(a.source)
}

// Restore SPIKE Nexus backing using recovery shards when SPIKE Keepers cannot
// provide adequate shards and SPIKE Nexus cannot recall its root key either.
//
// This is a break-the-glass superuser-only operation that a well-architected
// SPIKE deployment should not need.
//
// Parameters:
//   - shard *[32]byte: Pointer to a 32-byte array containing the shard to seed
//
// Returns:
//   - *data.RestorationStatus: Status of the restoration process if successful
//   - error: nil on success, unauthorized error if not authorized, or
//     wrapped error on request/parsing failure
//
// Example:
//
//	status, err := api.Restore(shardPtr)
func (a *Api) Restore(index int, shard *[32]byte) (*data.RestorationStatus, error) {
	return operator.Restore(a.source, index, shard)
}
