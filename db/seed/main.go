package main

import (
	"fmt"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
)

func main() {
	err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
	if err != nil {
		fmt.Printf("Failed to initialize: %#v", err)
	}
	defer dataservices.Close()

	db := dataservices.GetDB()

	// create test apps
	testApps := []models.App{
		models.App{
			AppSlug: "test-app-slug",
		},
		models.App{
			AppSlug: "9447ce3906c65e2d",
		},
	}
	for idx := range testApps {
		if err := db.Create(&testApps[idx]).Error; err != nil {
			fmt.Printf("Failed to seed db with app: %#v, app: %#v", err, testApps[idx])
		}
	}

	// create test app versions for test app
	for _, testAppVersion := range []models.AppVersion{
		models.AppVersion{
			AppID:     testApps[0].ID,
			Platform:  "ios",
			BuildSlug: "build-slug-ios",
		},
		models.AppVersion{
			AppID:     testApps[0].ID,
			Platform:  "android",
			BuildSlug: "build-slug-android",
		},
		models.AppVersion{
			AppID:     testApps[1].ID,
			Platform:  "ios",
			BuildSlug: "42e8d9b1def4074d",
		},
	} {
		if err := db.Create(&testAppVersion).Error; err != nil {
			fmt.Printf("Failed to seed db with app version: %#v, app version: %#v", err, testAppVersion)
		}
	}
}
