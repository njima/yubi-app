package entity

import "github.com/uptrace/bun"

type TaskCategoryType struct {
	bun.BaseModel `bun:"table:task_category_type,alias:tct"`

	ID   string `bun:"id,pk,type:varchar(36)"`
	Slug string `bun:"slug,type:varchar(100),notnull"`
	Name string `bun:"name,type:varchar(255),notnull"`
}

type TaskTag struct {
	bun.BaseModel `bun:"table:task_tag,alias:tt"`

	ID             string `bun:"id,pk,type:varchar(36)"`
	Name           string `bun:"name,type:varchar(255),notnull"`
	CategoryTypeID string `bun:"category_type_id,type:varchar(36),notnull"`

	// relations
	CategoryType *TaskCategoryType `bun:"rel:belongs-to,join:category_type_id=id"`
}

type TaskTagAssignment struct {
	bun.BaseModel `bun:"table:task_tag_assignment,alias:tta"`

	TaskID string `bun:"task_id,type:varchar(36),notnull"`
	TagID  string `bun:"tag_id,type:varchar(36),notnull"`

	// relations
	Tag *TaskTag `bun:"rel:belongs-to,join:tag_id=id"`
}

var TaskCategoryTypeTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*TaskCategoryType)(nil)).
		IfNotExists()
}

var TaskTagTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*TaskTag)(nil)).
		IfNotExists().
		ForeignKey(`("category_type_id") REFERENCES "task_category_type" ("id")`)
}

var TaskTagAssignmentTableCreator = func(db *bun.DB) *bun.CreateTableQuery {
	return db.NewCreateTable().
		Model((*TaskTagAssignment)(nil)).
		IfNotExists().
		ForeignKey(`("task_id") REFERENCES "task" ("id_natural") ON DELETE CASCADE`).
		ForeignKey(`("tag_id") REFERENCES "task_tag" ("id") ON DELETE CASCADE`)
}
