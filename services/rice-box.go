package services

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "workflows.yml",
		FileModTime: time.Unix(1567614147, 0),

		Content: string("format_version: '7'\ndefault_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git\napp:\n  envs:\n    - SHIP_ADDON_CONFIG_ANDROID: $CONFIG_JSON_URL\nworkflows:\n  resign_archive_app_store:\n    steps:\n      - activate-ssh-key@4.0.3:\n          run_if: '{{getenv \"SSH_RSA_PRIVATE_KEY\" | ne \"\"}}'\n      - git::git@github.com:bitrise-io/addons-ship-metadata-downloader-ios.git@update:\n          inputs:\n            - bitrise_ship_data_source: '$CONFIG_JSON_URL'\n      - certificate-and-profile-installer@1.10.1: {}\n      - script@1.1.5:\n          inputs:\n            - content: |-\n                #!/usr/bin/env bash\n                # fail if any commands fails\n                set -ex\n\n                mkdir zip_tmp\n                unzip -o $BITRISE_SHIP_ARTIFACT -d ./zip_tmp\n                mv zip_tmp/*.xcarchive ./ship.xcarchive\n      - export-xcarchive@1.0.1:\n          inputs:\n            - export_method: app-store\n            - archive_path: './ship.xcarchive'\n            - upload_bitcode: '$BITRISE_SHIP_INCLUDE_BITCODE'\n            - custom_export_options_plist_content: '$BITRISE_SHIP_CUSTOM_EXPORT_OPTION_PLIST'\n            - team_id: '$BITRISE_SHIP_FORCE_TEAM'\n      - git::git@github.com:bitrise-io/addons-ship-bg-worker-task-ios.git@master:\n          inputs:\n            - apple_user: '$BITRISE_SHIP_APPLE_USER'\n            - apple_app_specific_password: '$BITRISE_SHIP_APP_SPECIFIC_PASSWORD'\n            - sku: '$BITRISE_SHIP_SKU'\n            - metadata_file_structure_path: '$BITRISE_SHIP_DATA_PATH'\n  resign_android:\n    title: Re-sign Android artifact and deploy to store\n    steps:\n      - activate-ssh-key@4.0.3:\n          run_if: '{{getenv \"SSH_RSA_PRIVATE_KEY\" | ne \"\"}}'\n      - git::https://github.com/bitrise-io/addons-ship-bg-worker-task-android-prepare.git@master:\n      - sign-apk:\n          run_if: true\n          inputs:\n            - android_app: '$APP_LIST'\n            - keystore_url: '$KEYSTORE_URL'\n            - keystore_password: '$KEYSTORE_PASSWORD'\n            - keystore_alias: '$KEYSTORE_ALIAS'\n            - private_key_password: '$KEYSTORE_PRIVATE_KEY_PASSWORD'\n      - google-play-deploy:\n          inputs:\n            - service_account_json_key_path: '$SERVICE_ACCOUNT_JSON_URL'\n            - package_name: '$PACKAGE_NAME'\n            - expansionfile_path: '$EXPANSION_FILE_PATH'\n            - track: '$TRACK'\n            - whatsnews_dir: '$WHATS_NEW_DIR_PATH'\n            - mapping_file: '$MAPPING_PATH'\n      - git::https://github.com/bitrise-io/addons-ship-bg-worker-task-android-sync.git@master:\n          inputs:\n            - service_account_json_key_path: '$SERVICE_ACCOUNT_JSON_URL'\n            - package_name: '$PACKAGE_NAME'\n            - metadata_dir_path: '$METADATA_DIR_PATH'\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1567595557, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "workflows.yml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../utility`, &embedded.EmbeddedBox{
		Name: `../utility`,
		Time: time.Unix(1567595557, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"workflows.yml": file2,
		},
	})
}
