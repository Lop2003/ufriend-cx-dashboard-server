package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ufriend-cx-dashboard-server/internal/follow_up/query/dto"
)

type FollowUpQueryRepository interface {
	FindAll(ctx context.Context, filter *dto.FollowUpFilter) ([]*dto.FollowUpResponse, int64, error)
}

// IMongoCollection defines the subset of mongo.Collection methods used by the repository for testing
type IMongoCollection interface {
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error)
}

type mongoCollectionWrapper struct {
	coll *mongo.Collection
}

func (w *mongoCollectionWrapper) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (*mongo.Cursor, error) {
	return w.coll.Aggregate(ctx, pipeline, opts...)
}

type followUpQueryRepository struct {
	collection IMongoCollection
}

func NewFollowUpQueryRepository(db *mongo.Database) FollowUpQueryRepository {
	return &followUpQueryRepository{
		collection: &mongoCollectionWrapper{coll: db.Collection("follow_ups")},
	}
}
