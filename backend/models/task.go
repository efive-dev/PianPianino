package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Importance int

const (
	NotSet Importance = iota
	Low
	Medium
	High
)

type Task struct {
	bun.BaseModel `bun:"table:tasks"`

	ID          int64      `bun:"id,pk,autoincrement"`
	UserID      int64      `bun:"user_id,notnull"`
	User        *User      `bun:"rel:belongs-to,join:user_id=id,on_delete:cascade,on_update:cascade"`
	Description string     `bun:"description"`
	Priority    Importance `bun:"importance,notnull,default:0"`
	Completed   bool       `bun:"completed,notnull,default:false"`
	CreatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp"`
}
