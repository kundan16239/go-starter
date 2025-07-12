package repository

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// DocumentToStruct converts a BSON document to a struct
func DocumentToStruct(doc bson.M, target interface{}) error {
	if doc == nil {
		return fmt.Errorf("document is nil")
	}

	data, err := bson.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	err = bson.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to struct: %w", err)
	}

	return nil
}

// StructToDocument converts a struct to a BSON document
func StructToDocument(source interface{}) (bson.M, error) {
	data, err := bson.Marshal(source)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	var doc bson.M
	err = bson.Unmarshal(data, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to document: %w", err)
	}

	return doc, nil
}

// DocumentsToStructs converts a slice of BSON documents to a slice of structs
func DocumentsToStructs(docs []bson.M, target interface{}) error {
	if len(docs) == 0 {
		return nil
	}

	data, err := bson.Marshal(docs)
	if err != nil {
		return fmt.Errorf("failed to marshal documents: %w", err)
	}

	err = bson.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to structs: %w", err)
	}

	return nil
}
