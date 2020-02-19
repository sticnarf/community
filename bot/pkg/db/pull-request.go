package db

import "time"

// PullRequest struct
type PullRequest struct {
	ID        int       `gorm:"id"`
	PrID      int       `gorm:"pr_id"`
	Owner     string    `gorm:"owner"`
	Repo      string    `gorm:"repo"`
	Title     string    `gorm:"title"`
	Label     string    `gorm:"label"`
	Merge     bool      `gorm:"merge"`
	CreatedAt time.Time `gorm:"created_at"`
}
