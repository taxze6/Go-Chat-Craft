package models

type UserStory struct {
	Model
	OwnerId uint   `json:"owner_id"`
	Content string `json:"content"`
	Media   string `json:"media"`
	Type    int    `json:"type"`
}

type ResponseUserStory struct {
	Story         UserStory           `json:"story"`
	StoryLikes    *[]UserStoryLike    `json:"story_likes"`
	StoryComments *[]UserStoryComment `json:"story_comments"`
}
type UserShowStoryListResponse struct {
	StoryList *[]ResponseUserStory
	Count     int
}

func (table *UserStory) UserStoryTableName() string {
	return "user_story"
}
