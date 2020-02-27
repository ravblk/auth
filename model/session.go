package model

type (
	SessionClient struct {
		SessionHash string `json:"session_hash" db:"session_hash" valid:"required"`
		UserID      string `json:"user_id" db:"user_id" valid:"required"`
	}
	Session struct {
		SessionClient
		IP        string `db:"ip" valid:"required, ip"`
		UserAgent string `db:"user_agent" valid:""`
	}
)
