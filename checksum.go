package zep_nakama

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/heroiclabs/nakama-common/runtime"
)

// RequestPayload represents the structure of the request payload
type RequestPayload struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

// ResponsePayload represents the structure of the response payload
type ResponsePayload struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
	Content string `json:"content,omitempty"`
}

// CalculateHash computes the SHA-256 hash of the given content
func CalculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// ReadFileContent reads the content of a file at the given path
func ReadFileContent(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// RpcCheckSum handles the RPC request
func RpcCheckSum(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var reqPayload RequestPayload
	err := json.Unmarshal([]byte(payload), &reqPayload)
	if err != nil {
		return "", runtime.NewError("Invalid request payload", http.StatusBadRequest)
	}

	if reqPayload.Type == "" {
		reqPayload.Type = "core"
	}
	if reqPayload.Version == "" {
		reqPayload.Version = "1.0.0"
	}

	basePath := os.Getenv("FILE_BASE_PATH")
	if basePath == "" {
		return "", runtime.NewError("FILE_BASE_PATH environment variable is not set", http.StatusInternalServerError)
	}

	filePath := filepath.Join(basePath, reqPayload.Type, reqPayload.Version+".json")
	fileContent, err := ReadFileContent(filePath)
	if err != nil {
		return "", runtime.NewError("File not found", http.StatusNotFound)
	}

	responsePayload := ResponsePayload{
		Type:    reqPayload.Type,
		Version: reqPayload.Version,
		Hash:    CalculateHash(fileContent),
	}

	if reqPayload.Hash != "" && reqPayload.Hash != responsePayload.Hash {
		responsePayload.Content = ""
	} else {
		responsePayload.Content = fileContent
	}

	_, err = db.Exec(`INSERT INTO file_data (type, version, hash, content) VALUES ($1, $2, $3, $4)`,
		responsePayload.Type, responsePayload.Version, responsePayload.Hash, responsePayload.Content)
	if err != nil {
		return "", runtime.NewError("Failed to save to database", http.StatusInternalServerError)
	}

	responsePayloadBytes, err := json.Marshal(responsePayload)
	if err != nil {
		return "", err
	}

	return string(responsePayloadBytes), nil
}
