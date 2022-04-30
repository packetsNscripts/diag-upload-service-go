package main

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

//Functional test 1: HomePage responds with Diag Service
func TestRespondsWithDiagService(t *testing.T) {

	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	resource, err := pool.Run("diag-upload-service-go", "latest", []string{})
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})

	var resp *http.Response
	err = pool.Retry(func() error {
		resp, err = http.Get(fmt.Sprint("http://localhost:", resource.GetPort("8000/tcp"), "/"))
		if err != nil {
			t.Log("waiting for container startup..")
			return err
		}
		return nil
	})

	require.NoError(t, err, "HTTP Error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code is not 200")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Could not read HTTP body")

	require.Equal(t, string(body), "Diag Service", "Site home page does not return 'Diag Service'")

}
