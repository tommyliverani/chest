package chest

import (
	"chest/internal/common"
	"chest/internal/factory"
	"chest/internal/jewel"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

const AES_KIND = "aes"

type AesChest struct {
	baseChest
	Salt            string `json:"salt"`
	EncryptedJewels string `json:"encrypted_jewels"`
}

func init() {
	factory.RegisterChestCreator(AES_KIND, CreateAesChest)
	factory.RegisterChestParser(AES_KIND, ParseAesChest)
}

func CreateAesChest(name string, description string) (factory.Chest, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	ac := &AesChest{
		baseChest: baseChest{
			Id:          common.GenerateChestID(),
			Name:        name,
			Kind:        AES_KIND,
			Description: description,
		},
		Salt: hex.EncodeToString(salt),
	}
	password := common.ReadSecret("Insert the chest password: ")
	confirm := common.ReadSecret("Confirm the chest password: ")
	if password != confirm {
		return nil, fmt.Errorf("passwords do not match")
	}
	// encrypt an empty jewel list so password can be verified on Open
	passwordJewel := &passwordJewelAdapter{password: password}
	if err := ac.encryptAndStore([]json.RawMessage{}, passwordJewel); err != nil {
		return nil, fmt.Errorf("failed to initialize chest: %w", err)
	}
	return ac, nil
}

func ParseAesChest(data json.RawMessage) (factory.Chest, error) {
	var ac AesChest
	if err := json.Unmarshal(data, &ac); err != nil {
		return nil, fmt.Errorf("failed to parse AesChest: %w", err)
	}
	return &ac, nil
}

func (a *AesChest) GetKeyJewelKind() string { return jewel.KEY }

func (a *AesChest) GetEmoji() string { return "🔒" }

func (a *AesChest) ToJson() (json.RawMessage, error) { return json.Marshal(a) }

func (a *AesChest) Delete() error { return nil }

func (a *AesChest) Close() error { return nil }

func (a *AesChest) Open(keyJewel factory.Jewel) error {
	_, err := a.decryptJewels(keyJewel)
	if err != nil {
		return fmt.Errorf("wrong password for chest '%s'", a.Name)
	}
	return nil
}

func (a *AesChest) Edit(keyJewel factory.Jewel) error {
	chestField := common.SelectField("which field do you want to edit?", []string{"description", "name"})
	if chestField == "name" {
		newName := common.ReadField("Insert new name: ")
		a.Name = newName
	}
	if chestField == "description" {
		newDescription := common.ReadField("Insert new description: ")
		a.Description = newDescription
	}
	return nil
}

func (a *AesChest) GetJewels(keyJewel factory.Jewel) ([]factory.Jewel, error) {
	rawJewels, err := a.getDecryptedRawJewels(keyJewel)
	if err != nil {
		return nil, err
	}
	var jewels []factory.Jewel
	for _, raw := range rawJewels {
		j, err := factory.ParseJewel(raw)
		if err != nil {
			return nil, err
		}
		jewels = append(jewels, j)
	}
	return jewels, nil
}

func (a *AesChest) AddJewel(jewelToAdd factory.Jewel, keyJewel factory.Jewel) error {
	rawJewels, err := a.getDecryptedRawJewels(keyJewel)
	if err != nil {
		return err
	}
	for _, raw := range rawJewels {
		name, err := common.GetNameFromJson(raw)
		if err != nil {
			return err
		}
		if name == jewelToAdd.GetName() {
			return fmt.Errorf("jewel with name '%s' already exists in chest", jewelToAdd.GetName())
		}
	}
	jewelBytes, err := jewelToAdd.ToJson()
	if err != nil {
		return fmt.Errorf("failed to marshal jewel: %w", err)
	}
	rawJewels = append(rawJewels, json.RawMessage(jewelBytes))
	return a.encryptAndStore(rawJewels, keyJewel)
}

func (a *AesChest) RemoveJewel(jewelToRemove factory.Jewel, keyJewel factory.Jewel) error {
	rawJewels, err := a.getDecryptedRawJewels(keyJewel)
	if err != nil {
		return err
	}
	jewelName := jewelToRemove.GetName()
	for i, raw := range rawJewels {
		name, err := common.GetNameFromJson(raw)
		if err != nil {
			return err
		}
		if name == jewelName {
			rawJewels = append(rawJewels[:i], rawJewels[i+1:]...)
			return a.encryptAndStore(rawJewels, keyJewel)
		}
	}
	return fmt.Errorf("jewel with name '%s' not found in chest", jewelName)
}

func (a *AesChest) UpdateJewel(jewelName string, jewelToUpdate factory.Jewel, keyJewel factory.Jewel) error {
	rawJewels, err := a.getDecryptedRawJewels(keyJewel)
	if err != nil {
		return err
	}
	jsonJewel, err := jewelToUpdate.ToJson()
	if err != nil {
		return err
	}
	for i, raw := range rawJewels {
		name, err := common.GetNameFromJson(raw)
		if err != nil {
			return err
		}
		if name == jewelName {
			rawJewels[i] = json.RawMessage(jsonJewel)
			return a.encryptAndStore(rawJewels, keyJewel)
		}
	}
	return fmt.Errorf("jewel with name '%s' not found in chest", jewelName)
}

// --- helpers ---

func (a *AesChest) getDecryptedRawJewels(keyJewel factory.Jewel) ([]json.RawMessage, error) {
	return a.decryptJewels(keyJewel)
}

func (a *AesChest) deriveKey(password string) ([]byte, error) {
	salt, err := hex.DecodeString(a.Salt)
	if err != nil {
		return nil, fmt.Errorf("invalid salt: %w", err)
	}
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(password))
	return h.Sum(nil), nil
}

func (a *AesChest) encryptAndStore(rawJewels []json.RawMessage, keyJewel factory.Jewel) error {
	password, err := getPasswordFromJewel(keyJewel)
	if err != nil {
		return err
	}
	key, err := a.deriveKey(password)
	if err != nil {
		return err
	}
	plaintext, err := json.Marshal(rawJewels)
	if err != nil {
		return fmt.Errorf("failed to marshal jewels: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed to create GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	a.EncryptedJewels = base64.StdEncoding.EncodeToString(ciphertext)
	return nil
}

func (a *AesChest) decryptJewels(keyJewel factory.Jewel) ([]json.RawMessage, error) {
	password, err := getPasswordFromJewel(keyJewel)
	if err != nil {
		return nil, err
	}
	key, err := a.deriveKey(password)
	if err != nil {
		return nil, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(a.EncryptedJewels)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted jewels: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed (wrong password?): %w", err)
	}
	var rawJewels []json.RawMessage
	if err := json.Unmarshal(plaintext, &rawJewels); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted jewels: %w", err)
	}
	return rawJewels, nil
}

func getPasswordFromJewel(j factory.Jewel) (string, error) {
	data, err := j.ToJson()
	if err != nil {
		return "", err
	}
	var helper struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(data, &helper); err != nil {
		return "", err
	}
	if helper.Key == "" {
		return "", fmt.Errorf("key jewel has no key field")
	}
	return helper.Key, nil
}

// passwordJewelAdapter wraps a raw password string to satisfy factory.Jewel
// used only during chest creation before a real key jewel exists
type passwordJewelAdapter struct {
	password string
}

func (p *passwordJewelAdapter) ToJson() (json.RawMessage, error) {
	return json.Marshal(struct {
		Key string `json:"key"`
	}{Key: p.password})
}

func (p *passwordJewelAdapter) GetName() string        { return "" }
func (p *passwordJewelAdapter) GetKind() string        { return "" }
func (p *passwordJewelAdapter) GetEmoji() string       { return "" }
func (p *passwordJewelAdapter) GetDescription() string { return "" }
func (p *passwordJewelAdapter) Edit() error            { return nil }
func (p *passwordJewelAdapter) Print()                 {}
func (p *passwordJewelAdapter) Copy()                  {}
func (p *passwordJewelAdapter) Use()                   {}
