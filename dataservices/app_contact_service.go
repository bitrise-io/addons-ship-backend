package dataservices

import "github.com/bitrise-io/addons-ship-backend/models"

// AppContactService ...
type AppContactService interface {
	Create(appContact *models.AppContact) (*models.AppContact, []error, error)
	Find(appContact *models.AppContact) (*models.AppContact, error)
	FindAll(app *models.App) ([]models.AppContact, error)
	Update(appContact *models.AppContact, whitelist []string) error
	Delete(appContact *models.AppContact) error
}
