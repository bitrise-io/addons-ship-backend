package models

import (
	"github.com/bitrise-io/api-utils/structs"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// UpdatabeModelService ...
type UpdatabeModelService struct{}

// WhiteListedUpdate ...
func (u *UpdatabeModelService) WhiteListedUpdate(db *gorm.DB, object interface{}, whiteList []string) (interface{}, []error, error) {
	if len(whiteList) < 1 {
		return nil, nil, errors.New("No attributes to update")
	}

	updateData := map[string]interface{}{}
	for _, attribute := range whiteList {
		dbFieldName, err := structs.GetFieldNameByAttributeNameAndTag(object, attribute, "db")
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		fieldValue, err := structs.GetValueByAttributeName(object, attribute)
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		updateData[dbFieldName] = fieldValue
	}

	result := db.Model(object).Updates(updateData)
	verrs := result.GetErrors()
	if len(verrs) > 0 {
		return nil, verrs, nil
	}
	if result.Error != nil {
		return nil, nil, result.Error
	}
	return object, nil, nil
}
