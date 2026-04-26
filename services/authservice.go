package services

import (
	"app/utils"
	"encoding/json"
	"errors"
)

type SessionToken []SpaceInfo

type SpaceInfo struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}

type AuthService struct {
	query *DBService
}

func NewAuthService(db *DBService) *AuthService {
	return &AuthService{query: db}
}

func (au *AuthService) CreateSessionToken(info SpaceInfo) (string, error) {
	token := SessionToken{info}
	return au.encryptTokens(token)
}

func (au *AuthService) AddToSessionToken(existing string, info SpaceInfo) (string, error) {
	tokens, err := au.DecryptTokens(existing)
	if err != nil {
		return au.CreateSessionToken(info)
	}

	for i, token := range tokens {
		if token.ID == info.ID {
			tokens[i] = info
		} else {
			tokens = append(tokens, info)
		}

	}
	return au.encryptTokens(tokens)
}

func (au *AuthService) encryptTokens(tokens SessionToken) (string, error) {
	data, err := json.Marshal(tokens)
	if err != nil {
		return "", err
	}
	return utils.Encrypt(string(data))
}

func (au *AuthService) Authenticate(id, key string) (SpaceInfo, error) {
	if id == "" || key == "" {
		return SpaceInfo{}, errors.New("missing required inputs: id, key")
	}

	areaInfo, err := au.query.Spaces.Authenticate(id, key)

	if err != nil {
		return SpaceInfo{}, err
	}

	return areaInfo, nil
}

func (au *AuthService) DecryptTokens(cookie string) (SessionToken, error) {
	decrypted, err := utils.Decrypt(cookie)
	if err != nil {
		return nil, err
	}
	var tokens SessionToken
	err = json.Unmarshal([]byte(decrypted), &tokens)
	return tokens, err
}

func (au *AuthService) HasAccess(session, id string) bool {
	if session == "" {
		return false
	}

	decrypted, err := au.DecryptTokens(session)

	if err != nil {
		return false
	}

	for _, info := range decrypted {
		if info.ID == id {
			return true
		}
	}

	return false
}
