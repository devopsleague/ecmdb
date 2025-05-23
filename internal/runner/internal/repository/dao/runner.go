package dao

import (
	"context"
	"fmt"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = mongo.ErrNoDocuments

const (
	RunnerCollection = "c_runner"
)

type RunnerDAO interface {
	CreateRunner(ctx context.Context, r Runner) (int64, error)
	Update(ctx context.Context, req Runner) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	Detail(ctx context.Context, id int64) (Runner, error)
	ListRunner(ctx context.Context, offset, limit int64) ([]Runner, error)
	Count(ctx context.Context) (int64, error)
	FindByCodebookUid(ctx context.Context, codebookUid string, tag string) (Runner, error)
	ListByCodebookUids(ctx context.Context, codebookUids []string) ([]Runner, error)
	ListByIds(ctx context.Context, ids []int64) ([]Runner, error)
	ListTagsPipelineByCodebookUid(ctx context.Context) ([]RunnerPipeline, error)
}

func NewRunnerDAO(db *mongox.Mongo) RunnerDAO {
	return &runnerDAO{
		db: db,
	}
}

type runnerDAO struct {
	db *mongox.Mongo
}

func (dao *runnerDAO) ListByIds(ctx context.Context, ids []int64) ([]Runner, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{"id": bson.M{"$in": ids}}

	cursor, err := col.Find(ctx, filter)
	var result []Runner
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解码错误: %w", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("游标遍历错误: %w", err)
	}
	return result, nil
}

func (dao *runnerDAO) ListByCodebookUids(ctx context.Context, codebookUids []string) ([]Runner, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{"codebook_uid": bson.M{"$in": codebookUids}}

	cursor, err := col.Find(ctx, filter)
	var result []Runner
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解码错误: %w", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("游标遍历错误: %w", err)
	}
	return result, nil
}

func (dao *runnerDAO) Detail(ctx context.Context, id int64) (Runner, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{"id": id}

	var result Runner
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return Runner{}, fmt.Errorf("解码错误，%w", err)
	}

	return result, nil
}

func (dao *runnerDAO) Delete(ctx context.Context, id int64) (int64, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{"id": id}

	result, err := col.DeleteOne(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("删除文档错误: %w", err)
	}

	return result.DeletedCount, nil
}

func (dao *runnerDAO) Update(ctx context.Context, req Runner) (int64, error) {
	col := dao.db.Collection(RunnerCollection)
	updateDoc := bson.M{
		"$set": bson.M{
			"name":            req.Name,
			"codebook_secret": req.CodebookSecret,
			"worker_name":     req.WorkerName,
			"topic":           req.Topic,
			"tags":            req.Tags,
			"desc":            req.Desc,
			"variables":       req.Variables,
			"utime":           time.Now().UnixMilli(),
		},
	}
	filter := bson.M{"id": req.Id}
	count, err := col.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return 0, fmt.Errorf("修改文档操作: %w", err)
	}

	return count.ModifiedCount, nil
}

func (dao *runnerDAO) FindByCodebookUid(ctx context.Context, codebookUid string, tag string) (Runner, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{}
	filter["codebook_uid"] = codebookUid
	filter["tags"] = bson.M{
		"$elemMatch": bson.M{"$eq": tag},
	}

	var result Runner
	if err := col.FindOne(ctx, filter).Decode(&result); err != nil {
		return Runner{}, fmt.Errorf("解码错误，%w", err)
	}

	return result, nil
}

func (dao *runnerDAO) CreateRunner(ctx context.Context, r Runner) (int64, error) {
	r.Id = dao.db.GetIdGenerator(RunnerCollection)
	col := dao.db.Collection(RunnerCollection)
	now := time.Now()
	r.Ctime, r.Utime = now.UnixMilli(), now.UnixMilli()

	_, err := col.InsertOne(ctx, r)
	if err != nil {
		return 0, fmt.Errorf("插入数据错误: %w", err)
	}

	return r.Id, nil
}

func (dao *runnerDAO) ListRunner(ctx context.Context, offset, limit int64) ([]Runner, error) {
	col := dao.db.Collection(RunnerCollection)
	filter := bson.M{}
	opts := &options.FindOptions{
		Sort:  bson.D{{Key: "ctime", Value: -1}},
		Limit: &limit,
		Skip:  &offset,
	}

	cursor, err := col.Find(ctx, filter, opts)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询错误, %w", err)
	}

	var result []Runner
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解码错误: %w", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("游标遍历错误: %w", err)
	}
	return result, nil
}

func (dao *runnerDAO) Count(ctx context.Context) (int64, error) {
	col := dao.db.Collection(RunnerCollection)
	filer := bson.M{}

	count, err := col.CountDocuments(ctx, filer)
	if err != nil {
		return 0, fmt.Errorf("文档计数错误: %w", err)
	}

	return count, nil
}

func (dao *runnerDAO) ListTagsPipelineByCodebookUid(ctx context.Context) ([]RunnerPipeline, error) {
	col := dao.db.Collection(RunnerCollection)
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", "$codebook_uid"},
			// 使用 $push 累加器将选择的字段添加到 runners 数组中
			{"runner_tags", bson.D{{"$push", bson.D{
				{"tags", "$tags"},
				{"topic", "$topic"},
			}}}},
		}}},
	}

	cursor, err := col.Aggregate(ctx, pipeline)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, fmt.Errorf("查询错误, %w", err)
	}

	var result []RunnerPipeline
	if err = cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("解码错误: %w", err)
	}
	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("游标遍历错误: %w", err)
	}

	return result, nil
}

type Runner struct {
	Id             int64       `bson:"id"`
	Name           string      `bson:"name"`
	CodebookUid    string      `bson:"codebook_uid"`
	CodebookSecret string      `bson:"codebook_secret"`
	WorkerName     string      `bson:"worker_name"`
	Topic          string      `bson:"topic"`
	Tags           []string    `bson:"tags"`
	Action         uint8       `bson:"action"`
	Desc           string      `bson:"desc"`
	Variables      []Variables `bson:"variables"`
	Ctime          int64       `bson:"ctime"`
	Utime          int64       `bson:"utime"`
}

type Variables struct {
	Key    string `bson:"key"`
	Value  any    `bson:"value"`
	Secret bool   `bson:"secret"`
}

type RunnerPipeline struct {
	CodebookUid string       `bson:"_id"`
	RunnerTags  []RunnerTags `bson:"runner_tags"`
}

type RunnerTags struct {
	Topic string   `bson:"topic"`
	Tags  []string `json:"tags"`
}
