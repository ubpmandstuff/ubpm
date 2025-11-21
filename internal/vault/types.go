package vault

import "time"

const UbpmVersion = 1

type Vault struct {
	Version int    `json:"version"`
	Salt    []byte `json:"salt"`
	Data    []byte `json:"data"`
}

type VaultData struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	ID        string    `json:"id"`
	URIs      []string  `json:"uris"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
