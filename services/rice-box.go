package services

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "workflows.yml",
		FileModTime: time.Unix(1563348034, 0),

		Content: string("format_version: '7'\ndefault_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git\nworkflows:\n  resign_ios:\n    steps:\n      - certificate-and-profile-installer@1.10.1: {}\n      - script@1.1.5:\n          inputs:\n            - content: |-\n                #!/usr/bin/env bash\n                # fail if any commands fails\n                set -e\n                # debug log\n                set -x\n\n                brew install jq\n          title: Install jq\n      # - script:\n      #     title: Auth to Bitrise API\n      #     inputs:\n      #       - content: |\n      #           #!/usr/bin/env bash\n      #           set -ex\n\n      #           curl -H \"Authorization: ${BITRISE_ACCESS_TOKEN}\" https://api.bitrise.io/v0.1/me\n      - script:\n          title: Get artifact from the build\n          inputs:\n            - content: |\n                #!/usr/bin/env bash\n                set -ex\n\n                download_url=$(curl -X GET \"https://api.bitrise.io/v0.1/apps/${BITRISE_APP_SLUG}/builds/${BITRISE_BUILD_SLUG}/artifacts/${BITRISE_ARTIFACT_SLUG}\" -H \"accept: application/json\" -H \"Authorization: ${BITRISE_ACCESS_TOKEN}\" | jq -r '.data.expiring_download_url')\n                envman add --key BITRISE_DOWNLOAD_URL --value $download_url\n      - resource-archive@2.0.1:\n          inputs:\n            - extract_to_path: './'\n            - archive_url: '$BITRISE_DOWNLOAD_URL'\n      - export-xcarchive@1.0.1:\n          inputs:\n            - archive_path: unarchived/Xcode-10_default.xcarchive\n            - upload_bitcode: 'no'\n            - compile_bitcode: 'no'\n            - export_method: app-store\n      - deploy-to-bitrise-io:\n          inputs:\n            - notify_user_groups: none\n  resign_android:\n    title: Re-sign Android artifact and deploy to store\n    steps:\n      - git-clone@4.0.14: {}\n      - path::./prepare:\n      - sign-apk:\n          run_if: true\n          inputs:\n            - android_app: '$APP_LIST'\n            - keystore_url: '$KEYSTORE_URL'\n            - keystore_password: '$KEYSTORE_PASSWORD'\n            - keystore_alias: '$KEYSTORE_ALIAS'\n            - private_key_password: '$KEYSTORE_PRIVATE_KEY_PASSWORD'\n      - google-play-deploy:\n          inputs:\n            - service_account_json_key_path: '$SERVICE_ACCOUNT_JSON_URL'\n            - package_name: '$PACKAGE_NAME'\n            - expansionfile_path: '$EXPANSION_FILE_PATH'\n            - track: '$TRACK'\n            - whatsnews_dir: '$WHATS_NEW_DIR_PATH'\n            - mapping_file: '$MAPPING_PATH'\n      - path::./sync:\n          inputs:\n            - service_account_json_key_path: '$SERVICE_ACCOUNT_JSON_URL'\n            - package_name: '$PACKAGE_NAME'\n            - metadata_dir_path: '$METADATA_DIR_PATH'\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1563347865, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "workflows.yml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../utility`, &embedded.EmbeddedBox{
		Name: `../utility`,
		Time: time.Unix(1563347865, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"workflows.yml": file2,
		},
	})
}
