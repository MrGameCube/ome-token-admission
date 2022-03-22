package token_admission

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MrGameCube/ome-token-admission/token-admission/internal/stream"
	"github.com/MrGameCube/ome-token-admission/token-admission/internal/token"
	"github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	ErrInvalidSignature = errors.New("hmac signature invalid")
	ErrTokenMissing     = errors.New("token missing")
	ErrCantCreateToken  = errors.New("token creation failed")
)

// TokenAdmission contains the logic to admit OME stream requests based on tokens
type TokenAdmission struct {
	// db contains the SQLite DB which is used to store the token and stream information
	db         *sql.DB
	tokenRepo  *token.SQLiteRepository
	streamRepo *stream.SQLiteRepository
}

// New creates a new TokenAdmission object and initializes the needed SQLite Repositories
func New(db *sql.DB) (*TokenAdmission, error) {
	tokenRepo, err := token.NewSQLiteRepository(db)
	if err != nil {
		return nil, err
	}
	streamRepo, err := stream.NewSQLiteRepository(db)
	if err != nil {
		return nil, err
	}
	return &TokenAdmission{
		db:         db,
		tokenRepo:  tokenRepo,
		streamRepo: streamRepo,
	}, nil
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

	// TODO: use expiry information?
	return &OMEAdmissionResponse{
		Allowed: tA.canAccess(&admissionReq),
	}, nil
}

func (tA *TokenAdmission) canAccess(request *OMEAdmissionBody) bool {
	reqAppName, reqStreamName, reqToken := parseStreamFromURL(request.Request.URL)
	streamInfo, err := tA.streamRepo.FindByName(reqStreamName, reqAppName)

	if err != nil {
		log.Println("canAccess:", err)
		return false
	}

	if streamInfo.Public {
		return true
	}

	tokenData, err := tA.tokenRepo.FindByToken(reqToken)
	if err != nil {
		log.Println("canAccess:", err)
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

func (tA *TokenAdmission) CreateStream(options *ta_models.StreamRequest) (*ta_models.StreamResponse, error) {
	entity, err := tA.streamRepo.Create(options.StreamOptions)
	if err != nil {
		return nil, err
	}
	var streamResp = ta_models.StreamResponse{}
	streamResp.Entity = entity
	if options.CreateTokens {
		watchToken, err := tA.CreateToken(&ta_models.TokenOptions{
			Direction:   ta_models.DirectionOutgoing,
			Stream:      entity.StreamName,
			Application: entity.ApplicationName,
			ExpiresAt:   options.ExpireAt,
		})
		if err != nil {
			return &streamResp, ErrCantCreateToken
		}
		streamToken, err := tA.CreateToken(&ta_models.TokenOptions{
			Direction:   ta_models.DirectionIncoming,
			Stream:      entity.StreamName,
			Application: entity.ApplicationName,
			ExpiresAt:   options.ExpireAt,
		})
		if err != nil {
			return &streamResp, ErrCantCreateToken
		}
		streamResp.StreamToken = streamToken.Token
		streamResp.WatchToken = watchToken.Token
	}
	return &streamResp, nil
}

func (tA *TokenAdmission) CreateToken(options *ta_models.TokenOptions) (*ta_models.TokenEntity, error) {
	genToken, err := generateToken()
	if err != nil {
		return nil, err
	}
	tokenEntity, err := tA.tokenRepo.Create(ta_models.TokenEntity{
		Token:       genToken,
		Direction:   options.Direction,
		Stream:      options.Stream,
		Application: options.Application,
		ExpiresAt:   options.ExpiresAt,
	})
	return tokenEntity, err
}
