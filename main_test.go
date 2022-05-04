package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

//Sets up and destroys container used for all tests
func setup(t *testing.T) (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	require.NoError(t, err, "could not connect to Docker")

	resource, err := pool.Run("diag-upload-service-go", "latest", []string{})
	require.NoError(t, err, "could not start container")

	t.Cleanup(func() {
		require.NoError(t, pool.Purge(resource), "failed to remove container")
	})
	return pool, resource
}

//Functional test 1: HomePage responds with Diag Service
func TestHomePage(t *testing.T) {

	pool, resource := setup(t)

	var resp *http.Response
	url := fmt.Sprint("http://localhost:", resource.GetPort("8000/tcp"), "/")

	err := pool.Retry(func() error {
		resp, _ = http.Get(url)
		return nil
	})

	require.NoError(t, err, "HTTP Error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "HTTP status code is not 200")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Could not read HTTP body")

	require.Equal(t, string(body), "Diag Service", "Site home page does not return 'Diag Service'")

}

//Functional test 2: Upload function works, accepts only *.tgz files
func TestUpload(t *testing.T) {

	pool, resource := setup(t)

	url := fmt.Sprint("http://localhost:", resource.GetPort("8000/tcp"), "/upload")

	//Generating payload for *.tgz file test
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	part, _ := writer.CreateFormFile("diag", "upload.tgz")
	part.Write([]byte(`sample`))
	err1 := writer.Close()
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp *http.Response
	err = pool.Retry(func() error {
		resp, err = client.Do(req)
		if err != nil {
			t.Log("waiting for container startup..")
			return err
		}
		return nil
	})

	require.NoError(t, err, "HTTP Error")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Upload of *.tgz file was not successful")

	//Generating payload for non *.tgz file
	payload2 := &bytes.Buffer{}
	writer2 := multipart.NewWriter(payload2)
	part2, _ := writer2.CreateFormFile("diag", "upload.txt")
	part2.Write([]byte(`sample`))
	err2 := writer2.Close()
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	client2 := &http.Client{}

	req2, err := http.NewRequest("POST", url, payload2)
	if err != nil {
		fmt.Println(err)
		return
	}

	req2.Header.Set("Content-Type", writer2.FormDataContentType())

	var resp2 *http.Response
	err = pool.Retry(func() error {
		resp2, err = client2.Do(req2)
		if err != nil {
			t.Log("waiting for container startup..")
			return err
		}
		return nil
	})
	require.NoError(t, err, "HTTP Error")
	defer resp2.Body.Close()

	//Files without the .tgz extension must be rejected
	require.Equal(t, http.StatusUnsupportedMediaType, resp2.StatusCode, "Prevents uploading non *.tgz files")

}
