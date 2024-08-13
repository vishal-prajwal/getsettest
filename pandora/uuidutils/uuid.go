package uuidutils

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
)

func UUIDFromStringOrNil(str string) uuid.UUID {
	id, _ := uuid.Parse(strings.TrimSpace(str))
	return id
}

func FromStringToPtr(str string) *uuid.UUID {
	if str == "" {
		return nil
	}

	return UUIDFromNilStringOrNil(&str)
}

func PtrToString(id *uuid.UUID) string {
	if id == nil {
		return ""
	}

	return id.String()
}

func UUIDFromNilStringOrNil(str *string) *uuid.UUID {
	if str == nil {
		return nil
	}

	id, _ := uuid.Parse(strings.TrimSpace(*str))

	return &id
}

func NilUUIDToNilString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	return aws.String(id.String())
}

func UuidSliceToStringSlice(uuids []uuid.UUID) []string {
	var result []string

	for _, id := range uuids {
		result = append(result, id.String())
	}

	return result
}

func InUUIDSlice(needle uuid.UUID, slice []uuid.UUID) bool {
	for _, v := range slice {
		if v == needle {
			return true
		}
	}

	return false
}

func UniqueUUIDSlice(slice []uuid.UUID) []uuid.UUID {
	var result []uuid.UUID

	unique := make(map[uuid.UUID]bool)

	for _, id := range slice {
		if unique[id] {
			continue
		}

		unique[id] = true

		result = append(result, id)
	}

	return result
}

func StringSliceToUUID(slice []string) []uuid.UUID {
	var items []uuid.UUID

	for _, id := range slice {
		if parsedId, err := uuid.Parse(id); err == nil {
			items = append(items, parsedId)
		}
	}

	return items
}
