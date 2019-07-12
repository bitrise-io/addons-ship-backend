package services_test

import (
	"testing"

	"github.com/bitrise-io/addons-ship-backend/env"
	"github.com/bitrise-io/addons-ship-backend/services"
	ctxpkg "github.com/bitrise-io/api-utils/context"
	uuid "github.com/satori/go.uuid"
)

func Test_AppContactPutHandler(t *testing.T) {
	httpMethod := "PATCH"
	url := "/apps/{app-slug}/contacts/{contact-id}"
	handler := services.AppContactPutHandler

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
}
