package models

import (
	"fmt"
	"strings"
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

	ID          int64      `bun:"id,pk,autoincrement" json:"id"`
	UserID      int64      `bun:"user_id,notnull" json:"user_id"`
	User        *User      `bun:"rel:belongs-to,join:user_id=id,on_delete:cascade,on_update:cascade" json:"user,omitempty"`
	Description string     `bun:"description" json:"description"`
	Priority    Importance `bun:"importance,notnull,default:0" json:"priority"`
	Completed   bool       `bun:"completed,notnull,default:false" json:"completed"`
	CreatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time  `bun:",nullzero,notnull,default:current_timestamp" json:"updated_at"`
}

func (i *Importance) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`) // remove quotes

	switch strings.ToLower(s) {
	case "low":
		*i = Low
	case "normal", "medium":
		*i = Medium
	case "high":
		*i = High
	case "", "notset":
		*i = NotSet
	default:
		return fmt.Errorf("invalid priority value: %s", s)
	}

	return nil
}

func (i Importance) MarshalJSON() ([]byte, error) {
	var s string
	switch i {
	case Low:
		s = "low"
	case Medium:
		s = "normal"
	case High:
		s = "high"
	default:
		s = "notset"
	}
	return []byte(`"` + s + `"`), nil
}
