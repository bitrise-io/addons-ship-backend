package templates

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "email/confirmation.html",
		FileModTime: time.Unix(1562846211, 0),

		Content: string("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html>\n  <head></head>\n  <body>\n    <p>Hey {{ Name }}!</p>\n    <p>Ship wants to send you notifications about the activity of this app:</p>\n    <p>{{ AppTitle }}</p>\n    <p>Notification Settings:</p>\n    <p>New Version: {{ NewVersion }}</p>\n    <p>Successful Publish: {{ SuccessfulPublish }}</p>\n    <p>Failed Publish: {{ FailedPublish }}</p>\n    <a href=\"{{ URL }}\">Confirm Notifications</a>\n\n    <p>If you don’t want to get notifications from this app, just forget this email.</p>\n  </body>\n</html>\n"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "mail.html",
		FileModTime: time.Unix(1561461512, 0),

		Content: string("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html>\n  <head> </head>\n  <body>\n    <p>\n      Hello {{ Name }}\n      <a href=\"{{ URL }}\">Confirm email address</a>\n    </p>\n  </body>\n</html>\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "rice-box.go",
		FileModTime: time.Unix(1562847905, 0),

		Content: string(""),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "templates.go",
		FileModTime: time.Unix(1561461512, 0),

		Content: string("package templates\n\nimport (\n\t\"text/template\"\n\n\trice \"github.com/GeertJohan/go.rice\"\n\t\"github.com/bitrise-io/go-utils/templateutil\"\n\t\"github.com/pkg/errors\"\n)\n\n// Get ...\nfunc Get(templateFileName string, data map[string]interface{}) (string, error) {\n\ttemplateBox, err := rice.FindBox(\"\")\n\tif err != nil {\n\t\treturn \"\", errors.WithStack(err)\n\t}\n\n\ttmpContent, err := templateBox.String(templateFileName)\n\tif err != nil {\n\t\treturn \"\", errors.WithStack(err)\n\t}\n\n\tbody, err := templateutil.EvaluateTemplateStringToString(tmpContent, nil, template.FuncMap(data))\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn body, nil\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1562843918, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file4, // "mail.html"
			file5, // "rice-box.go"
			file6, // "templates.go"

		},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "email",
		DirModTime: time.Unix(1562843927, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "email/confirmation.html"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "email"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(``, &embedded.EmbeddedBox{
		Name: ``,
		Time: time.Unix(1562843918, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":      dir1,
			"email": dir2,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"email/confirmation.html": file3,
			"mail.html":               file4,
			"rice-box.go":             file5,
			"templates.go":            file6,
		},
	})
}
