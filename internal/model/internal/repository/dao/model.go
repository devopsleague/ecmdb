package dao

import (
	"context"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type ModelDAO interface {
	CreateModelGroup(ctx context.Context, mg ModelGroup) (int64, error)
}

func NewModelDAO(client *mongo.Client) ModelDAO {
	return &modelDAO{
		db: client.Database("cmdb"),
	}
}

type modelDAO struct {
	db *mongo.Database
}

func (m *modelDAO) CreateModelGroup(ctx context.Context, mg ModelGroup) (int64, error) {
	now := time.Now()
	mg.Ctime, mg.Utime = now.UnixMilli(), now.UnixMilli()
	mg.Id = mongox.GetDataID(m.db, "c_model_group")

	col := m.db.Collection("c_model_group")
	_, err := col.InsertMany(ctx, []interface{}{mg})

	if err != nil {
		return 0, err
	}

	return mg.Id, nil
}

type ModelGroup struct {
	Id    int64  `bson:"id"`
	Name  string `bson:"name"`
	Ctime int64  `bson:"ctime"`
	Utime int64  `bson:"utime"`
}