package repository

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/domain"
	"github.com/Duke1616/ecmdb/internal/attribute/internal/repository/dao"
)

type AttributeRepository interface {
	CreateAttribute(ctx context.Context, req domain.Attribute) (int64, error)
	SearchAttributeByModelUID(ctx context.Context, identifies string) (domain.AttributeProjection, error)
}

type attributeRepository struct {
	dao dao.AttributeDAO
}

func NewAttributeRepository(dao dao.AttributeDAO) AttributeRepository {
	return &attributeRepository{
		dao: dao,
	}
}

func (a *attributeRepository) CreateAttribute(ctx context.Context, req domain.Attribute) (int64, error) {
	return a.dao.CreateAttribute(ctx, dao.Attribute{
		ModelUID:  req.UID,
		Name:      req.Name,
		UID:       req.ModelUID,
		FieldType: req.FieldType,
		Required:  req.Required,
	})
}

// SearchAttributeByModelUID 查询对应模型的字段信息
func (a *attributeRepository) SearchAttributeByModelUID(ctx context.Context, modelUid string) (domain.AttributeProjection, error) {
	attributeList, err := a.dao.SearchAttributeByModelUID(ctx, modelUid)
	if err != nil {
		return domain.AttributeProjection{}, err
	}
	projection := make(map[string]int, 0)

	for _, ca := range attributeList {
		projection[ca.UID] = 1
	}

	return domain.AttributeProjection{
		Projection: projection,
	}, nil
}

func (a *attributeRepository) toDomain(modelDao *dao.Attribute) domain.Attribute {
	return domain.Attribute{
		ID:        modelDao.Id,
		Name:      modelDao.Name,
		UID:       modelDao.UID,
		FieldType: modelDao.FieldType,
		ModelUID:  modelDao.ModelUID,
		Required:  modelDao.Required,
	}
}
