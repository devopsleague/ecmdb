package integration

import (
	"context"
	"github.com/Duke1616/ecmdb/internal/relation/internal/integration/startup"
	"github.com/Duke1616/ecmdb/internal/relation/internal/repository/dao"
	"github.com/Duke1616/ecmdb/internal/relation/internal/web"
	"github.com/Duke1616/ecmdb/pkg/ginx/test"
	"github.com/Duke1616/ecmdb/pkg/mongox"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

type HandlerRMTestSuite struct {
	suite.Suite

	dao    dao.RelationModelDAO
	db     *mongox.Mongo
	server *gin.Engine
}

func (s *HandlerRMTestSuite) TearDownSuite() {
	_, err := s.db.Collection(dao.ModelRelationCollection).DeleteMany(context.Background(), bson.M{})
	require.NoError(s.T(), err)
	_, err = s.db.Collection("c_id_generator").DeleteMany(context.Background(), bson.M{})
	require.NoError(s.T(), err)
}

func (s *HandlerRMTestSuite) TearDownTest() {
	_, err := s.db.Collection(dao.ModelRelationCollection).DeleteMany(context.Background(), bson.M{})
	require.NoError(s.T(), err)
	_, err = s.db.Collection("c_id_generator").DeleteMany(context.Background(), bson.M{})
	require.NoError(s.T(), err)
}

func (s *HandlerRMTestSuite) SetupSuite() {
	handler, err := startup.InitRMHandler()
	require.NoError(s.T(), err)
	server := gin.Default()
	handler.RegisterRoute(server)

	s.db = startup.InitMongoDB()
	s.dao = dao.NewRelationModelDAO(s.db)
	s.server = server
}

func (s *HandlerRMTestSuite) TestCreate() {
	testCase := []struct {
		name string
		req  web.CreateModelRelationReq

		wantCode int
		wantResp test.Result[int64]
	}{
		{
			name: "创建成功",
		},
	}

	for _, tc := range testCase {
		s.T().Run(tc.name, func(t *testing.T) {

		})
	}
}

func TestRMHandler(t *testing.T) {
	suite.Run(t, new(HandlerRMTestSuite))
}
