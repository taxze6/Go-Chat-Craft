package dao

import (
	"GoChatCraft/global"
	"GoChatCraft/models"
	"errors"
	"go.uber.org/zap"
)

func GetStoryList(userId uint, page int, pageSize int) (*[]models.UserStory, error) {
	offset := (page - 1) * pageSize
	story := make([]models.UserStory, 0)
	if tx := global.DB.Where("owner_id = ?", userId).Offset(offset).Limit(pageSize).Find(&story); tx.RowsAffected == 0 {
		zap.S().Info("story data found")
		return nil, errors.New("story data found")
	}
	return &story, nil
}

func GetUserShowStoryList(userId uint) (*[]models.ResponseUserStory, int, error) {
	story := make([]models.UserStory, 0)
	if tx := global.DB.Where("owner_id = ?", userId).Find(&story); tx.RowsAffected == 0 {
		zap.S().Info("story data found")
		//return nil, 0, errors.New("story data found")
	}
	likeStory := make([]models.UserStoryLike, 0)
	responseStory := make([]models.ResponseUserStory, 0)

	for _, s := range story {
		currentLikeStory := make([]models.UserStoryLike, 0)
		if tx := global.DB.Where("user_story_id = ?", s.ID).Find(&currentLikeStory); tx.RowsAffected == 0 {
			zap.S().Info("story like data found")
			//return nil, 0, errors.New("story like data found")
		} else {
			//...Used to pass the elements of a slice or array to a function one by one.
			likeStory = append(likeStory, currentLikeStory...)
			responseStory = append(responseStory, models.ResponseUserStory{
				Story:      s,
				StoryLikes: &currentLikeStory,
			})
		}
	}
	likeStoryCount := len(likeStory)
	//Get the latest three data of story.
	var latestStories []models.UserStory
	if len(story) <= 3 {
		latestStories = make([]models.UserStory, len(story))
		copy(latestStories, story)
	} else {
		latestStories = make([]models.UserStory, 3)
		copy(latestStories, story[len(story)-3:])
	}
	return &responseStory, likeStoryCount, nil
}

func AddStory(story *models.UserStory) (*models.UserStory, error) {
	tx := global.DB.Create(&story)
	if tx.RowsAffected == 0 {
		zap.S().Info("failed to add a new story")
		return nil, errors.New("failed to add a new story")
	}
	return story, nil
}

func AddStoryLike(likeStory *models.UserStoryLike) error {
	tx := global.DB.Create(&likeStory)
	if tx.RowsAffected == 0 {
		zap.S().Info("failed to add a new story like")
		return errors.New("failed to add a new story like")
	}
	return nil
}

func AddStoryComment(commentStory *models.UserStoryComment) error {
	tx := global.DB.Create(&commentStory)
	if tx.RowsAffected == 0 {
		zap.S().Info("failed to add a new story like")
		return errors.New("failed to add a new story like")
	}
	return nil
}
