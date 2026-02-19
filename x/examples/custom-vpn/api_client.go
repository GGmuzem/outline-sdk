package main

import (
	"time"
)

type Server struct {
	ID        string
	Country   string
	Flag      string // Emoji or icon name
	Config    string // ss:// key
	Latency   int    // in ms
	IsPremium bool
}

type UserPlan string

const (
	PlanFree    UserPlan = "Free"
	PlanPremium UserPlan = "Premium"
)

type UserInfo struct {
	ID           string
	Email        string
	Plan         UserPlan
	ExpiryDate   time.Time
	ActiveServer *Server
}

// FetchServerList simulates an API call to get available VPN servers
func FetchServerList() []Server {
	return []Server{
		{"us-1", "USA", "🇺🇸", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 45, false},
		{"de-1", "Germany", "🇩🇪", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 60, false},
		{"jp-1", "Japan", "🇯🇵", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 150, true},
		{"uk-1", "UK", "🇬🇧", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 55, true},
		{"nl-1", "Netherlands", "🇳🇱", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 50, true},
		{"sg-1", "Singapore", "🇸🇬", "ss://YWVzLTEyOC1nY206dGVzdA@127.0.0.1:0", 120, true},
	}
}

// GetUserInfo simulates an API call to get the current user's profile and subscription
func GetUserInfo() UserInfo {
	return UserInfo{
		ID:         "user_123",
		Email:      "client@drfrake.com",
		Plan:       PlanFree, // Change to PlanPremium to test premium features
		ExpiryDate: time.Now().AddDate(0, 1, 0),
	}
}
