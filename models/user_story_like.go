package models

type UserStoryLike struct {
	Model
	UserStoryId uint `json:"user_story_id"`
	LikeOwnerId uint `json:"like_owner_id"`
}

func (table *UserStoryLike) UserStoryLikeTableName() string {
	return "user_story_like"
}
