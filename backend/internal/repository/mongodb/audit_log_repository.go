package mongodb

import (
	"context"
	"time"

	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	"github.com/n2pluto/cinema-booking-system/internal/domain/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type auditLogRepository struct {
	col *mongo.Collection
}

func NewAuditLogRepository(col *mongo.Collection) repository.AuditLogRepository {
	return &auditLogRepository{col: col}
}

func (r *auditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	log.ID = bson.NewObjectID()
	log.Timestamp = time.Now()
	_, err := r.col.InsertOne(ctx, log)
	return err
}

func (r *auditLogRepository) List(ctx context.Context, f repository.AuditLogFilter) ([]*entity.AuditLog, int64, error) {
	filter := bson.D{}

	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 {
		limit = 20
	}
	skip := int64((page - 1) * limit)

	opts := options.Find().
		SetSort(bson.D{{Key: "timestamp", Value: -1}}).
		SetSkip(skip).
		SetLimit(int64(limit))

	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []*entity.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
