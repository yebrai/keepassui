package secretsreader

import (
	"keepassui/internal/secretsdb"
	"slices"
)

//go:generate mockgen -destination=../mocks/secretsreader/mock_secretsreader.go -source=./secretsreader.go

type DefaultSecretsReader struct {
	UriID          string
	ContentInBytes []byte
	Password       string
}

var loadedDB *secretsdb.SecretsDB

func CreateDefaultSecretsReader(uriID string, contentInBytes []byte, password string) DefaultSecretsReader {
	return DefaultSecretsReader{
		UriID:          uriID,
		ContentInBytes: contentInBytes,
		Password:       password,
	}
}

type SecretReader interface {
	ReadEntriesFromContentGroupedByPath() error
	GetUriID() string
	GetFirstPath() string
	GetEntriesForPath(path string) []secretsdb.SecretEntry
	WriteDBBytes() ([]byte, error)
	AddSecretEntry(secretEntry secretsdb.SecretEntry)
	ModifySecretEntry(originalTitle, originalGroup string, originalIsGroup bool, secretEntry secretsdb.SecretEntry)
	DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool
}

func (dsr DefaultSecretsReader) ReadEntriesFromContentGroupedByPath() error {
	secretsDB, err := secretsdb.ReadSecretsDBFromDBBytes(dsr.ContentInBytes, dsr.Password)
	if err == nil {
		loadedDB = &secretsDB
	}
	return err
}

func (dsr DefaultSecretsReader) GetUriID() string {
	return dsr.UriID
}

func (dsr DefaultSecretsReader) GetFirstPath() string {
	return loadedDB.PathsInOrder[0]
}

func (dsr DefaultSecretsReader) GetEntriesForPath(path string) []secretsdb.SecretEntry {
	return loadedDB.EntriesByPath[path]
}

func (dsr DefaultSecretsReader) WriteDBBytes() ([]byte, error) {
	return loadedDB.WriteDBBytes(dsr.Password)
}

func (dsr DefaultSecretsReader) AddSecretEntry(secretEntry secretsdb.SecretEntry) {
	loadedDB.AddSecretEntry(secretEntry)
}

func (dsr DefaultSecretsReader) ModifySecretEntry(originalTitle, originalGroup string, originalIsGroup bool, secretEntry secretsdb.SecretEntry) {
	entries := loadedDB.EntriesByPath[originalGroup]
	i := slices.IndexFunc(entries, func(se secretsdb.SecretEntry) bool {
		return se.Title == originalTitle && se.IsGroup == originalIsGroup
	})
	if i != -1 {
		entries[i].Group = secretEntry.Group
		entries[i].Title = secretEntry.Title
		entries[i].Username = secretEntry.Username
		entries[i].Password = secretEntry.Password
		entries[i].Notes = secretEntry.Notes
	}
}

func (dsr DefaultSecretsReader) DeleteSecretEntry(secretEntry secretsdb.SecretEntry) bool {
	return loadedDB.DeleteSecretEntry(secretEntry)
}
