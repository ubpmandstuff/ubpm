/* Package vault provides utilities for working with ubpm vaults,
 * as well as implementations of cryptographic functions.
 */
package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	VaultDir      = ".ubpm"
	VaultFilename = "vault.ubpm.json"
)

// Init initializes a ubpm vault using a provided password and a mountpoint
func Init(mountpoint string, password []byte) (string, error) {
	// path vars
	vaultDir := filepath.Join(mountpoint, VaultDir)
	vaultPath := filepath.Join(vaultDir, VaultFilename)

	// check if vault already exists
	if _, err := os.Stat(vaultPath); err == nil {
		return "", fmt.Errorf("vault already exists at %s", vaultPath)
	}

	// create vault dir
	if err := os.MkdirAll(vaultDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create vault directory: %w", err)
	}

	// generate salt
	salt, err := GenerateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// derive key from given password and salt
	key := DeriveKey(password, salt)

	// create and marshal empty vault data
	emptyVault := VaultData{
		Entries: []Entry{},
	}

	plaintext, err := json.Marshal(emptyVault)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vault data json: %w", err)
	}

	// encrypt empty vault data
	encrypedData, err := Encrypt(plaintext, key)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %w", err)
	}

	// fill output file
	fileContent := VaultFile{
		Version: UbpmVersion,
		Salt:    salt,
		Data:    encrypedData,
	}

	// marshal output file
	output, err := json.MarshalIndent(fileContent, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal vault json: %w", err)
	}

	// write output file to disk/usb
	if err := os.WriteFile(vaultPath, output, 0o600); err != nil {
		return "", fmt.Errorf("failed to write vault file: %w", err)
	}

	return vaultPath, nil
}

// Open opens a ubpm vault and returns it
func Open(path string, password []byte) (*Vault, error) {
	// read raw data and unmarshal it
	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var file VaultFile
	if err := json.Unmarshal(rawData, &file); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// derive key and decrypt data
	key := DeriveKey(password, file.Salt)

	plaintext, err := Decrypt(file.Data, key)
	if err != nil {
		return nil, fmt.Errorf("couldn't decrypt file: invalid password or corrupted file")
	}

	// unmarshal plaintext into data
	var data VaultData
	if err := json.Unmarshal(plaintext, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted data: %w", err)
	}

	return &Vault{
		Path: path,
		Key:  key,
		Salt: file.Salt,
		Data: &data,
	}, nil
}

func GenerateID() {}

// Save encrypts and writes all changes to the vault to disk
func (v *Vault) Save() error {
	// marshal data into json
	plaintext, err := json.Marshal(v.Data)
	if err != nil {
		return err
	}

	// encrypt marshalled json
	encrypedData, err := Encrypt(plaintext, v.Key)
	if err != nil {
		return err
	}

	// fill vault file struct
	fileContent := &VaultFile{
		Version: UbpmVersion,
		Salt:    v.Salt,
		Data:    encrypedData,
	}

	// marshal vault file
	output, err := json.MarshalIndent(fileContent, "", "  ")
	if err != nil {
		return err
	}

	// write output to disk
	if err := os.WriteFile(v.Path, output, 0o600); err != nil {
		return err
	}

	return nil
}

// AddEntry adds an entry from specified args to a ubpm vault
func (v *Vault) AddEntry(website, username, password, notes string) error {
	data := Entry{
		Website:  website,
		Username: username,
		Password: password,
		Notes:    notes,
	}

	v.Data.Entries = append(v.Data.Entries, data)
	return v.Save()
}

// FindEntry attempts to find an entry by id, otherwise returns an error
func (v *Vault) FindEntry(id string) (*Entry, error) {
	for i := range v.Data.Entries {
		if len(id) >= 4 && strings.HasPrefix(v.Data.Entries[i].ID, id) {
			return &v.Data.Entries[i], nil
		}
	}
	return nil, fmt.Errorf("entry not found: %s", id)
}

// RemoveEntry attempts find an entry by id and remove it, otherwise returns an error
func (v *Vault) RemoveEntry(id string) error {
	for i := range v.Data.Entries {
		if len(id) >= 4 && strings.HasPrefix(v.Data.Entries[i].ID, id) {
			v.Data.Entries = append(v.Data.Entries[:i], v.Data.Entries[i+1:]...)
			return v.Save()
		}
	}
	return fmt.Errorf("entry not found: %s", id)
}
