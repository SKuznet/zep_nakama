package zep_nakama

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// TestCalculateHash tests the CalculateHash function
func TestCalculateHash(t *testing.T) {
	content := "test content"
	expectedHash := sha256.Sum256([]byte(content))
	expectedHashStr := hex.EncodeToString(expectedHash[:])
	hash := CalculateHash(content)
	assert.Equal(t, expectedHashStr, hash)
}

// TestReadFileContent tests the ReadFileContent function
func TestReadFileContent(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := ioutil.TempFile("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Write content to the temporary file
	content := "test content"
	_, err = tmpFile.WriteString(content)
	assert.NoError(t, err)

	// Close the temporary file
	err = tmpFile.Close()
	assert.NoError(t, err)

	// Read the content of the file using ReadFileContent
	readContent, err := ReadFileContent(tmpFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, content, readContent)

	// Test case when file is not found
	_, err = ReadFileContent("nonexistentfile")
	assert.Error(t, err)
}

// TestRpcCheckSum tests the RpcCheckSum function using SQLite
func TestRpcCheckSum(t *testing.T) {
	// Create temporary directory and file for testing
	tmpDir, err := ioutil.TempDir("", "testdir")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	coreFilePath := tmpDir + "/core/1.0.0.json"
	err = os.MkdirAll(tmpDir+"/core", os.ModePerm)
	assert.NoError(t, err)
	err = ioutil.WriteFile(coreFilePath, []byte("file content"), 0644)
	assert.NoError(t, err)

	// Set the environment variable to the temporary directory
	os.Setenv("FILE_BASE_PATH", tmpDir)
	defer os.Unsetenv("FILE_BASE_PATH")

	// Create an in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	// Create the table
	_, err = db.Exec(`CREATE TABLE file_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		version TEXT,
		hash TEXT,
		content TEXT
	)`)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		payload    string
		expected   string
		shouldFail bool
	}{
		{
			name:       "Valid payload with matching hash",
			payload:    `{"type":"core","version":"1.0.0","hash":"e0ac3601005dfa1864f5392aabaf7d898b1b5bab854f1acb4491bcd806b76b0c"}`,
			expected:   `{"type":"core","version":"1.0.0","hash":"e0ac3601005dfa1864f5392aabaf7d898b1b5bab854f1acb4491bcd806b76b0c","content":"file content"}`,
			shouldFail: false,
		},
		{
			name:       "Valid payload with non-matching hash",
			payload:    `{"type":"core","version":"1.0.0","hash":"invalidhash"}`,
			expected:   `{"type":"core","version":"1.0.0","hash":"e0ac3601005dfa1864f5392aabaf7d898b1b5bab854f1acb4491bcd806b76b0c","content":""}`,
			shouldFail: false,
		},
		{
			name:       "Missing file",
			payload:    `{"type":"core","version":"2.0.0","hash":"somehash"}`,
			expected:   `{"error":"File not found"}`,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := RpcCheckSum(context.Background(), nil, db, nil, tt.payload)
			if tt.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				var respPayload ResponsePayload
				err = json.Unmarshal([]byte(resp), &respPayload)
				assert.NoError(t, err)

				var expectedPayload ResponsePayload
				err = json.Unmarshal([]byte(tt.expected), &expectedPayload)
				assert.NoError(t, err)

				assert.Equal(t, expectedPayload, respPayload)
			}
		})
	}
}
