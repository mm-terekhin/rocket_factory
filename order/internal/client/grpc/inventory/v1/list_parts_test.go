package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	inventoryV1 "github.com/mm-terekhin/rocket_factory/shared/pkg/proto/inventory/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ClientSuite) TestSuccessListParts() {
	partUUID := gofakeit.UUID()
	now := timestamppb.Now()

	s.grpc.On("ListParts", s.ctx, mock.MatchedBy(func(req *inventoryV1.ListPartsRequest) bool {
		return len(req.GetFilter().GetUuids()) == 1 && req.GetFilter().GetUuids()[0] == partUUID
	})).Return(&inventoryV1.ListPartsResponse{
		Parts: []*inventoryV1.Part{
			{
				Uuid: partUUID,
				Info: &inventoryV1.PartInfo{
					Name:  gofakeit.Name(),
					Price: 123.45,
				},
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
	}, nil)

	parts, err := s.client.ListParts(s.ctx, model.PartsFilter{Uuids: []string{partUUID}})

	s.Require().NoError(err)
	s.Require().Len(parts, 1)
	s.Require().Equal(partUUID, parts[0].UUID)
	s.Require().Equal(123.45, parts[0].Info.Price)
}

func (s *ClientSuite) TestFailureListParts() {
	partUUID := gofakeit.UUID()
	serviceErr := gofakeit.Error()

	s.grpc.On("ListParts", s.ctx, mock.Anything).Return(nil, serviceErr)

	_, err := s.client.ListParts(s.ctx, model.PartsFilter{Uuids: []string{partUUID}})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
}
