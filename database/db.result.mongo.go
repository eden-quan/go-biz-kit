package database

import "go.mongodb.org/mongo-driver/mongo"

type DbDeleteResult struct {
	DeletedCount int64
}

func NewDbDeleteResult(result *mongo.DeleteResult) *DbDeleteResult {
	var deleteCount int64 = 0
	if result != nil {
		deleteCount = result.DeletedCount
	}
	return &DbDeleteResult{
		DeletedCount: deleteCount,
	}
}

type DbUpdateResult struct {
	MatchedCount  int64       // The number of documents matched by the filter.
	ModifiedCount int64       // The number of documents modified by the operation.
	UpsertedCount int64       // The number of documents upserted by the operation.
	UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
}

func NewDbUpdateResult(result *mongo.UpdateResult) *DbUpdateResult {
	if result == nil {
		return &DbUpdateResult{}
	}

	return &DbUpdateResult{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    result.UpsertedID,
	}
}
