package models

type UserCollect struct {
	Model
	CollectOwnerId uint   `json:"collect_owner_id"`
	CollectContent string `json:"collect_content"`
	Type           int    `json:"type"`
}

func (table *UserCollect) UserCollectTableName() string {
	return "user_collect"
}
