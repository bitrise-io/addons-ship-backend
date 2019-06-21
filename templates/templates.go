package templates

import (
	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/bitrise-io/go-utils/templateutil"
	"github.com/pkg/errors"
)

// Get ...
func Get(templateFileName string, data map[string]interface{}) (string, error) {
	templateBox, err := rice.FindBox("")
	if err != nil {
		return "", errors.WithStack(err)
	}

	tmpContent, err := templateBox.String(templateFileName)
	if err != nil {
		return "", errors.WithStack(err)
	}

	body, err := templateutil.EvaluateTemplateStringToString(tmpContent, nil, template.FuncMap(data))
	if err != nil {
		return "", err
	}
	// t, err := template.ParseFiles(templateFileName)
	// if err != nil {
	// 	return "", err
	// }
	// buf := new(bytes.Buffer)
	// if err = t.Execute(buf, data); err != nil {
	// 	return "", err
	// }
	// return buf.String(), nil
	return body, nil
}
