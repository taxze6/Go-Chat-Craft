package models

type Relation struct {
	Model
	OwnerId  uint   //Users corresponding to the relationship.
	TargetId uint   //It corresponds to whom
	Type     int    //Relationship types: 1 represents friendship, 2 represents group relationship.
	Desc     string //Description
}

func (r *Relation) RelTableName() string {
	return "relation"
}
