package models

import (
	"github.com/bitrise-io/api-utils/structs"
	"github.com/pkg/errors"
)

// UpdatableModelService ...
type UpdatableModelService struct{}

// UpdateData ...
func (u *UpdatableModelService) UpdateData(object interface{}, whiteList []string) (map[string]interface{}, error) {
	if len(whiteList) < 1 {
		return nil, errors.New("No attributes to update")
	}
	whiteList = append(whiteList, "UpdatedAt")
	updateData := map[string]interface{}{}
	for _, attribute := range whiteList {
		dbFieldName, err := structs.GetFieldNameByAttributeNameAndTag(object, attribute, "json")
		if err != nil {
			return nil, errors.WithStack(err)
		}
		fieldValue, err := structs.GetValueByAttributeName(object, attribute)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		updateData[dbFieldName] = fieldValue
	}
	return updateData, nil
}
