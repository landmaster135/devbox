package domain

import "time"

// Session はエージェントのセッション情報を表す。
type Session struct {
	UUID         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Branch       string
	CWD          string
	Conversation string
}
