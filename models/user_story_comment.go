package models

type UserStoryComment struct {
	Model
	UserStoryId    uint   `json:"user_story_id"`
	CommentOwnerId uint   `json:"comment_owner_id"`
	CommentContent string `json:"comment_content"`
	Type           int    `json:"type"`
}

type ResponseUserStoryComment struct {
	StoryComment *UserStoryComment `json:"story_comment"`
	UserAvatar   string            `json:"user_avatar"`
}

func (table *UserStoryComment) UserStoryCommentTableName() string {
	return "user_story_comment"
}
