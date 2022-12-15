package util

import (
	"necross.it/backend/database"
)

type Service struct {
	Server database.Server
}

func (s *Service) EntryExists(dest interface{}, tableName string, selector string, identifier string) bool {
	err := s.Server.DB.Get(dest, "SELECT * FROM ? WHERE ? = ?", tableName, selector, identifier)
	return err == nil
}

func (s *Service) GetEntryCount(tableName string) (int, error) {
	var count int
	err := s.Server.DB.Get(&count, "SELECT COUNT(*) FROM ?", tableName)
	return count, err
}
