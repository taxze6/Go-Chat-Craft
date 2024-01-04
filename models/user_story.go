package models

type UserStory struct {
	Model
	OwnerId uint   `json:"owner_id"`
	Content string `json:"content"`
	media   string `json:"media"`
	Type    int    `json:"type"`
}

type UserStoryLike struct {
	Model
	UserStoryId uint `json:"user_story_id"`
	LikeOwnerId uint `json:"like_owner_id"`
}

type UserStoryComment struct {
	Model
	UserStoryId    uint   `json:"user_story_id"`
	CommentOwnerId uint   `json:"comment_owner_id"`
	CommentContent string `json:"comment_content"`
	Type           int    `json:"type"`
}

func (table *UserStory) UserStoryTableName() string {
	return "user_story"
}

func (table *UserStoryLike) UserStoryLikeTableName() string {
	return "user_story_like"
}

func (table *UserStoryComment) UserStoryCommentTableName() string {
	return "user_story_comment"
}
