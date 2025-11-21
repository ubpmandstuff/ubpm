/* Package vault provides utilities for working with ubpm vaults,
 * as well as implementations of cryptographic functions.
 */
package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	VaultDir      = ".ubpm"
	VaultFilename = "vault.ubpm.json"
)

// InitVault initializes a ubpm vault using a provided password and a mountpoint
func InitVault(mountpoint string, passwd []byte) (string, error) {
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
	key := DeriveKey(passwd, salt)

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
	fileContent := Vault{
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
