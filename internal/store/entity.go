package store

import (
	"fmt"
	"log/slog"

	"github.com/bradenhc/kolob/internal/crypto"
	"github.com/bradenhc/kolob/internal/model"
)

// A GroupEntity bundles together encrypted group information and queryable metadata about the
// group. The metadata allows the entity to be easily fetched based on search criteria without
// needing to decrypt the group information.
type GroupEntity struct {
	Id            model.Uuid
	GroupHash     crypto.DataHash
	PassSalt      crypto.Salt
	PassHash      crypto.PassHash
	CreatedAt     int64
	UpdatedAt     int64
	EncryptedKey  []byte
	EncryptedData []byte
}

// NewGroupEntity creates a new entity of group inforation to persist in some way. It performs the
// obfuscation of queryable metadata about a group, encrypts the group information, and combines the
// metadata with the encrypted group information.
func NewGroupEntity(g *model.Group, pass crypto.Password, dkey crypto.Key) (GroupEntity, error) {
	// Generate a second key using the user password that will encrypt the data key (pass key)
	psalt, err := crypto.NewSalt()
	if err != nil {
		var e GroupEntity
		return e, fmt.Errorf("failed to create salt for group: %v", err)
	}
	slog.Info("Deriving password key")
	pkey := crypto.NewDerivedKey(pass, psalt)

	// Encyrypt the data key using the pass key before we store it in the database
	slog.Info("Encrypting data key")
	ekey, err := crypto.Encrypt(pkey, dkey)
	if err != nil {
		var e GroupEntity
		return e, fmt.Errorf("failed to encrypt data key: %v", err)
	}

	// Hash the group id so we can use it for authentication without leaking information
	ghash := crypto.HashData(g.Gid())

	// Hash the password so we can use it for authentication
	phash, err := crypto.HashPassword(pass)
	if err != nil {
		var e GroupEntity
		return e, fmt.Errorf("failed to generated group password hash: %v", err)
	}

	// Encrypt the group information to protect privacy
	slog.Info("Encrypting group information")
	edata, err := crypto.Encrypt(dkey, g.Table().Bytes)
	if err != nil {
		var e GroupEntity
		return e, fmt.Errorf("failed to encrypt group data before storing in database: %v", err)
	}

	// Create and return the group entity that bundles all the information we persist
	return GroupEntity{
		Id:            model.Uuid(g.Id()),
		GroupHash:     ghash,
		PassSalt:      psalt,
		PassHash:      phash,
		CreatedAt:     g.Created(),
		UpdatedAt:     g.Updated(),
		EncryptedKey:  ekey,
		EncryptedData: edata,
	}, nil
}
