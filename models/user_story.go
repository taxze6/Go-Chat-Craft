package models

type UserStory struct {
	Model
	OwnerId uint   `json:"owner_id"`
	Content string `json:"content"`
	Media   string `json:"media"`
	Type    int    `json:"type"`
}

type UserShowStoryListResponse struct {
	StoryList *[]UserStory
	Count     int
}

func (table *UserStory) UserStoryTableName() string {
	return "user_story"
}
