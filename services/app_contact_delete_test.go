package services_test

import (
	"net/http"
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactDeleteHandler(t *testing.T) {
	httpMethod := "DELETE"
	url := "/apps/{app-slug}/contacts/{contact-id}"
	handler := services.AppContactDeleteHandler

	behavesAsServiceCravingHandler(t, httpMethod, url, handler, []string{"AppContactService"}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	behavesAsContextCravingHandler(t, httpMethod, url, handler, []ctxpkg.RequestContextKey{services.ContextKeyAuthorizedAppContactID}, ControllerTestCase{
		contextElements: map[ctxpkg.RequestContextKey]interface{}{
			services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
		},
		env: &env.AppEnv{
			AppContactService: &testAppContactService{},
		},
	})

	t.Run("ok - minimal", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					deleteFn: func(appContact *models.AppContact) error {
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactDeleteResponse{
				Data: &models.AppContact{},
			},
		})
	})

	t.Run("ok - more complex", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{Email: "someones@email.addr"}, nil
					},
					deleteFn: func(appContact *models.AppContact) error {
						return nil
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: services.AppContactDeleteResponse{
				Data: &models.AppContact{Email: "someones@email.addr"},
			},
		})
	})

	t.Run("when error happens at finding app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return nil, gorm.ErrRecordNotFound
					},
					deleteFn: func(appContact *models.AppContact) error {
						return nil
					},
				},
			},
			expectedInternalErr: "SQL Error: record not found",
		})
	})

	t.Run("when error happens at deleting app contact", func(t *testing.T) {
		performControllerTest(t, httpMethod, url, handler, ControllerTestCase{
			contextElements: map[ctxpkg.RequestContextKey]interface{}{
				services.ContextKeyAuthorizedAppContactID: uuid.NewV4(),
			},
			env: &env.AppEnv{
				AppContactService: &testAppContactService{
					findFn: func(appContact *models.AppContact) (*models.AppContact, error) {
						return &models.AppContact{}, nil
					},
					deleteFn: func(appContact *models.AppContact) error {
						return errors.New("SOME-SQL-ERROR")
					},
				},
			},
			expectedInternalErr: "SQL Error: SOME-SQL-ERROR",
		})
	})
}
