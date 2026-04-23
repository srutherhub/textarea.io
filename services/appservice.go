package services

import (
	"app/utils"
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strconv"
	"strings"
)

type AppService struct {
	query *DBService
}

func NewAppService(db *DBService) *AppService {
	return &AppService{query: db}
}

func (as *AppService) CreateSpace(name string) (id string, key string, spaceName string, err error) {
	id = as.generateSlug()
	key = as.generateKey()
	encryptedKey, err := utils.Encrypt(key)

	if err != nil {
		return "", "", "", err
	}

	if name == "" {
		name = id
	}

	err = as.query.Spaces.CreateSpace(id, name, encryptedKey)

	if err != nil {
		return "", "", "", err
	}

	return id, key, name, nil
}

func (as *AppService) GetByID(id string) (Space, error) {
	space, err := as.query.Spaces.GetByID(id)

	if err != nil {
		return Space{}, err
	}

	space.Key, _ = utils.Decrypt(space.Key)

	if err != nil {
		return Space{}, err
	}

	go as.query.Spaces.UpdateLastAccessed(id)

	return space, err
}

func (as *AppService) AddContent(id, content string) error {
	err := as.query.Spaces.UpdateContent(id, content)

	if err != nil {
		return err
	}

	return nil

}

func (as *AppService) generateSlug() string {
	b := make([]byte, 6)
	rand.Read(b)
	return strings.ToLower(base64.URLEncoding.EncodeToString(b)[:8])
}

func (as *AppService) generateKey() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(90000))
	return strconv.FormatInt(n.Int64()+10000, 10)
}
