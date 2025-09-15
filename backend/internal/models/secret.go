package models

import (
	"time"

	"gorm.io/gorm"
)

/*
AES-GCM is an authenticated encryption algorithm (AEAD). It gives you both confidentiality (encryption) and integrity (tamper detection).
When you encrypt with AES-GCM, you get three outputs:

1. Ciphertext: encrypted data.
2. Nonce (IV):

Random value, usually 12 bytes (96 bits).
Must be unique per key for security.
Think of it as a “randomizer” to make identical plaintext encryptions different each time.

3. Tag (authentication tag / MAC):

Typically 16 bytes.
Ensures the ciphertext hasn’t been modified.
Without the correct key + nonce, the tag verification will fail.

encrypt flow
----------------
@startuml
actor "App" as App
participant "DEK Generator" as DEK
participant "AES-GCM (DEK)" as AES_DEK
participant "AES-GCM (KEK[v])" as AES_KEK
database "Database" as DB

App -> DEK: Generate 32-byte Data Encryption Key (DEK)
App -> AES_DEK: Encrypt secret with DEK + random nonce
AES_DEK -> App: Ciphertext + Nonce + Tag

App -> AES_KEK: Wrap DEK with KEK[v] + random nonce
AES_KEK -> App: Wrapped_DEK + WrapNonce + WrapTag

App -> DB: Store {Ciphertext, Nonce, Tag,\nWrapped_DEK, WrapNonce, WrapTag, kek_version}
@enduml


decrypt flow
----------------
@startuml
actor "App" as App
database "Database" as DB
participant "AES-GCM (KEK[v])" as AES_KEK
participant "AES-GCM (DEK)" as AES_DEK

App -> DB: Load record {Ciphertext, Nonce, Tag,\nWrapped_DEK, WrapNonce, WrapTag, kek_version}

App -> AES_KEK: Unwrap DEK with KEK[kek_version] + WrapNonce
AES_KEK -> App: DEK

App -> AES_DEK: Decrypt Ciphertext with DEK + Nonce
AES_DEK -> App: Plaintext Secret (in memory only!)

note right of App
  Zero DEK + plaintext buffers ASAP
  to reduce memory exposure
end note
@enduml
*/
// Secret stores encrypted secrets scoped by realm
type Secret struct {
	ID         string         `json:"id" gorm:"primaryKey;type:text"`
	RealmID    string         `json:"realm_id" gorm:"not null;type:text;index;uniqueIndex:uq_credential_realm_name"`
	Name       string         `json:"name" gorm:"not null;type:text;uniqueIndex:uq_credential_realm_name"`
	Group      string         `json:"group" gorm:"not null;type:text"`
	Desc       string         `json:"desc" gorm:"not null;type:text"`
	Path       string         `json:"path" gorm:"not null;type:text"`
	CipherAlg  string         `json:"cipher_alg" gorm:"not null;type:text"`
	CipherText string         `json:"cipher_text" gorm:"not null;type:text"`
	Nonce      string         `json:"nonce" gorm:"not null;type:text"`
	AuthTag    string         `json:"auth_tag" gorm:"not null;type:text"`
	WrappedDEK string         `json:"wrapped_dek" gorm:"not null;type:text"`
	KEKVersion int            `json:"kek_version" gorm:"not null"`
	CreatedBy  string         `json:"created_by" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedBy  string         `json:"updated_by" gorm:"type:text"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
