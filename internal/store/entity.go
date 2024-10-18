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

// Decrypt decrypts the entity's data buffer and returns a Group object using the underlying buffer.
func (e *GroupEntity) Decrypt(k crypto.Key) (*model.Group, error) {
	data, err := crypto.Decrypt(k, e.EncryptedData)
	if err != nil {
		return nil, err
	}

	return model.GetRootAsGroup(data, 0), nil
}

func (e *GroupEntity) Update(key crypto.Key, gid, name, desc []byte) (*model.Group, error) {
	prev, err := e.Decrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt group prior to update: %v", err)
	}

	next := model.GroupCloneWithUpdates(prev, gid, name, desc)

	e.GroupHash = crypto.HashData(next.Gid())
	e.UpdatedAt = next.Updated()

	edata, err := crypto.Encrypt(key, next.Table().Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt updated entity: %v", err)
	}
	e.EncryptedData = edata

	return next, nil
}

type MemberEntity struct {
	Id            model.Uuid
	UsernameHash  crypto.DataHash
	PassHash      crypto.PassHash
	CreatedAt     int64
	UpdatedAt     int64
	EncryptedData []byte
}

func NewMemberEntity(m *model.Member, pass crypto.Password, key crypto.Key) (MemberEntity, error) {
	// Hash the username so that we can store it in the DB without leaking information and use it
	// for fast lookup later
	uhash := crypto.HashData(m.Uname())

	// Hash the password so we can use it for authentication
	phash, err := crypto.HashPassword(pass)
	if err != nil {
		var e MemberEntity
		return e, fmt.Errorf("failed to hash member password: %v", err)
	}

	// Encrypt the member edata prior to storing it in the DB
	edata, err := crypto.Encrypt(key, m.Table().Bytes)
	if err != nil {
		var e MemberEntity
		return e, fmt.Errorf("failed to encrypt member data before storing: %v", err)
	}

	// Create and return the member entity that bundles all persisted information
	return MemberEntity{
		Id:            model.Uuid(m.Id()),
		UsernameHash:  uhash,
		PassHash:      phash,
		CreatedAt:     m.Created(),
		UpdatedAt:     m.Updated(),
		EncryptedData: edata,
	}, nil
}

func (e *MemberEntity) Decrypt(k crypto.Key) (*model.Member, error) {
	data, err := crypto.Decrypt(k, e.EncryptedData)
	if err != nil {
		return nil, err
	}

	return model.GetRootAsMember(data, 0), nil
}

func (e *MemberEntity) Update(k crypto.Key, uname, name []byte) (*model.Member, error) {
	prev, err := e.Decrypt(k)
	if err != nil {
		return nil, err
	}

	next := model.CloneMemberWithUpdates(prev, uname, name)

	e.UsernameHash = crypto.HashData(uname)
	e.UpdatedAt = next.Updated()
	edata, err := crypto.Encrypt(k, next.Table().Bytes)
	if err != nil {
		return nil, err
	}
	e.EncryptedData = edata

	return next, nil
}
