package models

import (
	"os"

	"github.com/bitrise-io/go-crypto/crypto"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// App ...
type App struct {
	Record
	AppSlug           string `json:"app_slug"`
	Plan              string `json:"plan"`
	BitriseAPIToken   string `json:"-"`
	APIToken          string `json:"-"`
	EncryptedSecret   []byte `json:"-"`
	EncryptedSecretIV []byte `json:"-"`

	AppVersions []AppVersion `gorm:"foreignkey:AppID" json:"app_versions"`
	AppSettings AppSettings  `gorm:"foreignkey:AppsID" json:"app_settings"`
}

// BeforeCreate ...
func (a *App) BeforeCreate(scope *gorm.Scope) error {
	if uuid.Equal(a.ID, uuid.UUID{}) {
		a.ID = uuid.NewV4()
	}

	if len(a.EncryptedSecretIV) != 0 {
		return nil
	}

	var err error
	secret, err := crypto.SecureRandomHex(12)
	if err != nil {
		return errors.Wrap(err, "Failed to generate secret")
	}
	for {
		a.EncryptedSecretIV, err = crypto.GenerateIV()
		if err != nil {
			return errors.WithStack(err)
		}
		var appWebhookCount int64
		err := scope.DB().Model(&App{}).Where("encrypted_secret_iv = ?", a.EncryptedSecretIV).Count(&appWebhookCount).Error
		if err != nil {
			return errors.WithStack(err)
		}
		if appWebhookCount == 0 {
			break
		}
	}

	err = a.encryptSecret(secret)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (a *App) encryptSecret(secret string) error {
	encryptKey, ok := os.LookupEnv("APP_WEBHOOK_SECRET_ENCRYPT_KEY")
	if !ok {
		return errors.New("No encrypt key provided")
	}
	encryptedSecret, err := crypto.AES256GCMCipher(secret, a.EncryptedSecretIV, encryptKey)
	if err != nil {
		return errors.WithStack(err)
	}
	a.EncryptedSecret = encryptedSecret

	return nil
}

// Secret ...
func (a *App) Secret() (string, error) {
	encryptKey, ok := os.LookupEnv("APP_WEBHOOK_SECRET_ENCRYPT_KEY")
	if !ok {
		return "", errors.New("No encrypt key provided")
	}
	secret, err := crypto.AES256GCMDecipher(a.EncryptedSecret, a.EncryptedSecretIV, encryptKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return secret, nil
}
