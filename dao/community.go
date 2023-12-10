package dao

import (
	"GoChatCraft/global"
	"GoChatCraft/models"
	"errors"
)

func CreateCommunity(community models.Community) (int, error) {
	com := models.Community{}
	if tx := global.DB.Where("name = ?", community.Name).First(&com); tx.RowsAffected == 1 {
		return -1, errors.New("the group record already exists")
	}
	tx := global.DB.Begin()
	if t := tx.Create(&community); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("failed to create group record")
	}
	relation := models.Relation{}
	relation.OwnerId = community.OwnerId
	relation.TargetId = community.ID
	relation.Type = 2
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("failed to create group record")
	}
	tx.Commit()
	return 0, nil
}

func GetCommunityList(ownerId uint) (*[]models.Community, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = 2", ownerId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("the group record does not exist")
	}
	communityID := make([]uint, 0)
	for _, v := range relation {
		cid := v.TargetId
		communityID = append(communityID, cid)
	}
	community := make([]models.Community, 0)
	if tx := global.DB.Where("id in ?", communityID).Find(&community); tx.RowsAffected == 0 {
		return nil, errors.New("failed to retrieve group data")
	}
	return &community, nil
}

// JoinCommunity Search and join the group based on the group nickname.
func JoinCommunity(ownerId uint, cname string) (int, error) {
	community := models.Community{}
	if tx := global.DB.Where("name = ?", cname).First(&community); tx.RowsAffected == 0 {
		return -1, errors.New("the group record does not exist")
	}
	relation := models.Relation{}
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 2", ownerId, community.ID).First(&relation); tx.RowsAffected == 1 {
		return -1, errors.New("the group has already been joined")
	}
	relation = models.Relation{}
	relation.OwnerId = ownerId
	relation.TargetId = community.ID
	relation.Type = 2

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		return -1, errors.New("failed to join")
	}
	return 0, nil
}
