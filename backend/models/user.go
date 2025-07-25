package models

import (
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID       int64  `bun:"id,pk,autoincrement"`
	Username string `bun:"username,unique,notnull"`
	Password string `bun:"password,notnull"`
	Tasks    []Task `bun:"tasks,rel:has-many,join:id=user_id"`
}
