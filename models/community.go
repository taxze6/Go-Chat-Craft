package models

import (
	"GoChatCraft/global"
	"errors"
)

type Community struct {
	Model
	Name    string
	OwnerId uint
	Type    int
	Image   string
	Desc    string
}

func (r *Community) ComTableName() string {
	return "community"
}

func FindUsersId(groupId uint) (*[]uint, error) {
	relation := make([]Relation, 0)
	if tx := global.DB.Where("target_id = ? and type = 2", groupId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("no member information found")
	}
	userIDs := make([]uint, 0)
	for _, v := range relation {
		userId := v.OwnerId
		userIDs = append(userIDs, userId)
	}
	return &userIDs, nil
}
