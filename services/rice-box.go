package services

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "workflows.yml",
		FileModTime: time.Unix(1562152026, 0),

		Content: string("format_version: '7'\ndefault_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git\nworkflows:\n  resign_archive_app_store:\n  steps:\n    - script@1.1.5:\n        inputs:\n          - content: |-\n              #!/usr/bin/env bash\n              set -ex\n\n              echo \"Hello lil' friend!\"\n        title: Say hello to my lil' friend\n    # - certificate-and-profile-installer@1.10.1: {}\n    # - script@1.1.5:\n    #     inputs:\n    #       - content: |-\n    #           #!/usr/bin/env bash\n    #           # fail if any commands fails\n    #           set -e\n    #           # debug log\n    #           set -x\n    #           brew install jq\n    #     title: Install jq\n    # # - script:\n    # #     title: Auth to Bitrise API\n    # #     inputs:\n    # #       - content: |\n    # #           #!/usr/bin/env bash\n    # #           set -ex\n    # #           curl -H \"Authorization: ${BITRISE_ACCESS_TOKEN}\" https://api.bitrise.io/v0.1/me\n    # - script:\n    #     title: Get artifact from the build\n    #     inputs:\n    #       - content: |\n    #           #!/usr/bin/env bash\n    #           set -ex\n    #           download_url=$(curl -X GET \"https://api.bitrise.io/v0.1/apps/${BITRISE_APP_SLUG}/builds/${BITRISE_BUILD_SLUG}/artifacts/${BITRISE_ARTIFACT_SLUG}\" -H \"accept: application/json\" -H \"Authorization: ${BITRISE_ACCESS_TOKEN}\" | jq -r '.data.expiring_download_url')\n    #           envman add --key BITRISE_DOWNLOAD_URL --value $download_url\n    # - resource-archive@2.0.1:\n    #     inputs:\n    #       - extract_to_path: './'\n    #       - archive_url: '$BITRISE_DOWNLOAD_URL'\n    # - export-xcarchive@1.0.1:\n    #     inputs:\n    #       - archive_path: unarchived/Xcode-10_default.xcarchive\n    #       - upload_bitcode: 'no'\n    #       - compile_bitcode: 'no'\n    #       - export_method: app-store\n    # - deploy-to-bitrise-io:\n    #     inputs:\n    #       - notify_user_groups: none\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1561701795, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "workflows.yml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../utility`, &embedded.EmbeddedBox{
		Name: `../utility`,
		Time: time.Unix(1561701795, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"workflows.yml": file2,
		},
	})
}
