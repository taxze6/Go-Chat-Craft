package models

type UserStory struct {
	Model
	OwnerId uint   `json:"owner_id"`
	Content string `json:"content"`
	Media   string `json:"media"`
	Type    int    `json:"type"`
}

func (table *UserStory) UserStoryTableName() string {
	return "user_story"
}
