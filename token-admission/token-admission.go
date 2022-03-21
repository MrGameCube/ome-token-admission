package token_admission

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MrGameCube/ome-token-admission/token-admission/internal/stream"
	"github.com/MrGameCube/ome-token-admission/token-admission/internal/token"
	"io/ioutil"
	"net/http"
)

var (
	ErrInvalidSignature = errors.New("hmac signature invalid")
	ErrTokenMissing     = errors.New("token missing")
)

// TokenAdmission contains the logic to admit OME stream requests based on tokens
type TokenAdmission struct {
	// db contains the SQLite DB which is used to store the token and stream information
	db         *sql.DB
	tokenRepo  *token.SQLiteRepository
	streamRepo *stream.SQLiteRepository
}

// New creates a new TokenAdmission object and initializes the needed SQLite Repositories
func New(db *sql.DB) *TokenAdmission {
	return &TokenAdmission{
		db:         db,
		tokenRepo:  token.NewSQLiteRepository(db),
		streamRepo: stream.NewSQLiteRepository(db),
	}
}

func (tA *TokenAdmission) HandleAdmissionRequest(request *http.Request) (*OMEAdmissionResponse, error) {

	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}

	if !ValidateHMACRequest(request, bodyBytes) {
		return nil, ErrInvalidSignature
	}

	admissionReq := OMEAdmissionBody{}
	err = json.Unmarshal(bodyBytes, &admissionReq)
	if err != nil {
		return nil, err
	}

	return &OMEAdmissionResponse{
		Allowed: tA.canAccess(&admissionReq),
	}, nil
}

func (tA *TokenAdmission) canAccess(request *OMEAdmissionBody) bool {
	reqAppName, reqStreamName, reqToken := parseStreamFromURL(request.Request.URL)
	streamInfo, err := tA.streamRepo.FindByName(reqStreamName, reqAppName)

	if err != nil {
		return false
	}

	if streamInfo.Public {
		return true
	}

	tokenData, err := tA.tokenRepo.FindByToken(reqToken)
	if err != nil {
		return false
	}

	if request.Request.Direction != tokenData.Direction {
		return false
	}

	if tokenData.Stream != reqStreamName || tokenData.Application != reqAppName {
		return false
	}

	return true
}
