package vault

import "time"

const UbpmVersion = 1

// type VaultFile defines the default, locked state of a ubpm vault
type VaultFile struct {
	Version int    `json:"version"`
	Salt    []byte `json:"salt"`
	Data    []byte `json:"data"`
}

// type Vault defines the decrypted state of a ubpm vault
type Vault struct {
	Path string
	Key  []byte
	Salt []byte
	Data *VaultData
}

// type VaultData defines the decrypted data of a vault
type VaultData struct {
	Entries []Entry `json:"entries"`
}

// type Entry describes a vault entry
type Entry struct {
	ID        string    `json:"id"`
	Website   string    `json:"website"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
