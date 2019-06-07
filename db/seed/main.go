package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/bitrise-io/addons-ship-backend/dataservices"
	"github.com/bitrise-io/addons-ship-backend/models"
	"github.com/go-yaml/yaml"
	"github.com/satori/go.uuid"
)

func main() {
	err := dataservices.InitializeConnection(dataservices.ConnectionParams{}, true)
	if err != nil {
		fmt.Printf("Failed to initialize: %#v", err)
	}
	defer dataservices.Close()

	db := dataservices.GetDB()

	var testData testData

	yamlFile, err := ioutil.ReadFile("test_data.yml")
	if err != nil {
		log.Printf("Failed to read test_data.yml: %s", err)
	}
	err = yaml.Unmarshal(yamlFile, &testData)
	if err != nil {
		log.Fatalf("Failed to parse test_data.yml: %v", err)
	}

	// create test apps
	for _, appData := range testData.Apps {
		app := models.App{
			Record:          models.Record{ID: appData.ID},
			AppSlug:         appData.Slug,
			APIToken:        appData.APIToken,
			BitriseAPIToken: appData.BitriseAPIToken,
			Plan:            appData.Plan,
		}
		if err := db.Create(&app).Error; err != nil {
			fmt.Printf("Failed to seed db with app: %#v, app: %#v", err, app)
		}
	}

	// create test app versions
	for _, appVersionData := range testData.AppVersions {
		appStoreInfoBytes, err := json.Marshal(appVersionData.AppStoreInfo)
		if err != nil {
			fmt.Printf("Failed to marshal app store info: %#v, app store info: %#v", err, appVersionData.AppStoreInfo)
		}
		appVersion := models.AppVersion{
			Record:           models.Record{ID: appVersionData.ID},
			AppID:            appVersionData.AppID,
			Platform:         appVersionData.Platform,
			Version:          appVersionData.Version,
			BuildNumber:      appVersionData.BuildNumber,
			LastUpdate:       appVersionData.LastUpdate,
			Scheme:           appVersionData.Scheme,
			Configuration:    appVersionData.Configuration,
			AppStoreInfoData: appStoreInfoBytes,
		}
		if err := db.Create(&appVersion).Error; err != nil {
			fmt.Printf("Failed to seed db with app version: %#v, app version: %#v", err, appVersion)
		}
	}
}

type app struct {
	ID              uuid.UUID `yaml:"id"`
	Slug            string    `yaml:"slug"`
	Plan            string    `yaml:"plan"`
	BitriseAPIToken string    `yaml:"bitrise_api_token"`
	APIToken        string    `yaml:"api_token"`
}

type appStoreInfo struct {
	ShortDescription string `yaml:"short_description"`
	FullDescription  string `yaml:"full_description"`
	WhatsNew         string `yaml:"whats_new"`
	PromotionalText  string `yaml:"promotional_text"`
	Keywords         string `yaml:"keywords"`
	ReviewNotes      string `yaml:"review_notes"`
	SupportURL       string `yaml:"support_url"`
	MarketingURL     string `yaml:"marketing_url"`
}

type appVersion struct {
	ID            uuid.UUID    `yaml:"id"`
	AppID         uuid.UUID    `yaml:"app_id"`
	Version       string       `yaml:"version"`
	Platform      string       `yaml:"platform"`
	BuildNumber   string       `yaml:"build_number"`
	BuildSlug     string       `yaml:"build_slug"`
	LastUpdate    time.Time    `yaml:"last_update"`
	Scheme        string       `yaml:"scheme"`
	Configuration string       `yaml:"configuration"`
	AppStoreInfo  appStoreInfo `yaml:"app_store_info"`
}

type testData struct {
	Apps        []app        `yaml:"apps"`
	AppVersions []appVersion `yaml:"app_versions"`
}
