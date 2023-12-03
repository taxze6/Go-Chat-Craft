package dao

import (
	"GoChatCraft/global"
	"GoChatCraft/models"
	"errors"
	"go.uber.org/zap"
)

func FriendList(userId uint) (*[]models.UserBasic, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = 1", userId).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("relation data found")
		return nil, errors.New("no friend relationship found")
	}
	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetId)
	}
	user := make([]models.UserBasic, 0)
	if tx := global.DB.Where("id in ?", userID).Find(&user); tx.RowsAffected == 0 {
		zap.S().Info("no friend relationship found in the relation data")
		return nil, errors.New("no friends found")
	}
	return &user, nil
}

// AddFriend Add friend using QR code on mobile device.
func AddFriend(userID, TargetId uint) (int, error) {
	if userID == TargetId {
		return -2, errors.New("the userID and TargetID are equal")
	}
	//Querying user by ID
	targetUser, err := FindUserId(TargetId)
	if err != nil {
		return -1, errors.New("no user found")
	}
	if targetUser.ID == 0 {
		zap.S().Info("no user found")
		return -1, errors.New("no user found")
	}
	relation := models.Relation{}
	//The purpose of these two query statements is to ensure that when adding friends,
	//the friend relationship can be checked correctly regardless of who initiates the request.
	//If only one query is performed, some situations may be missed,
	//such as when the friend relationship already exists but the query conditions do not match.
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 1", userID, TargetId).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("the friend exists")
		return 0, errors.New("the friend exists")
	}

	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 1", TargetId, userID).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("the friend exists")
		return 0, errors.New("the friend exists")
	}

	tx := global.DB.Begin()
	relation.OwnerId = userID
	relation.TargetId = targetUser.ID
	relation.Type = 1
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("failed to create friend record")
		//Transaction rollback
		tx.Rollback()
		return -1, errors.New("failed to create friend record")
	}
	relation = models.Relation{}
	relation.OwnerId = TargetId
	relation.TargetId = userID
	relation.Type = 1

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("failed to create friend record")

		//Transaction rollback
		tx.Rollback()
		return -1, errors.New("failed to create friend record")
	}

	tx.Commit()
	return 1, nil
}

func AddFriendByName(userId uint, targetName string) (int, error) {
	user, err := FindUserByName(targetName)
	if err != nil {
		return -1, errors.New("the user does not exist")
	}
	if user.ID == 0 {
		zap.S().Info("user not found")
		return -1, errors.New("the user does not exist")
	}
	return AddFriend(userId, user.ID)
}
