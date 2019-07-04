package templates

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "mail.html",
		FileModTime: time.Unix(1561461512, 0),

		Content: string("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n<html>\n  <head> </head>\n  <body>\n    <p>\n      Hello {{ Name }}\n      <a href=\"{{ URL }}\">Confirm email address</a>\n    </p>\n  </body>\n</html>\n"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "rice-box.go",
		FileModTime: time.Unix(1562152052, 0),

		Content: string(""),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "templates.go",
		FileModTime: time.Unix(1561461512, 0),

		Content: string("package templates\n\nimport (\n\t\"text/template\"\n\n\trice \"github.com/GeertJohan/go.rice\"\n\t\"github.com/bitrise-io/go-utils/templateutil\"\n\t\"github.com/pkg/errors\"\n)\n\n// Get ...\nfunc Get(templateFileName string, data map[string]interface{}) (string, error) {\n\ttemplateBox, err := rice.FindBox(\"\")\n\tif err != nil {\n\t\treturn \"\", errors.WithStack(err)\n\t}\n\n\ttmpContent, err := templateBox.String(templateFileName)\n\tif err != nil {\n\t\treturn \"\", errors.WithStack(err)\n\t}\n\n\tbody, err := templateutil.EvaluateTemplateStringToString(tmpContent, nil, template.FuncMap(data))\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn body, nil\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1561701795, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "mail.html"
			file3, // "rice-box.go"
			file4, // "templates.go"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(``, &embedded.EmbeddedBox{
		Name: ``,
		Time: time.Unix(1561701795, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"mail.html":    file2,
			"rice-box.go":  file3,
			"templates.go": file4,
		},
	})
}
