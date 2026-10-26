package order

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/mm-terekhin/rocket_factory/order/internal/model"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestSuccessCreate() {
	var (
		expectedUserUUID  = gofakeit.UUID()
		partUUID          = gofakeit.UUID()
		expectedPartUUIDs = []string{partUUID}
	)

	filter := model.PartsFilter{Uuids: expectedPartUUIDs}

	s.inventoryClient.On("ListParts", testCtx(), filter).Return([]model.Part{{UUID: partUUID, Info: model.PartInfo{Price: 100}}}, nil)
	s.orderRepository.On("CreateOrder", testCtx(), mock.MatchedBy(func(o model.Order) bool {
		return o.UserUUID == expectedUserUUID && o.TotalPrice == 100 && o.Status == string(model.StatusPendingPayment)
	})).Return(nil)

	res, _ := s.service.Create(testCtx(), model.CreateOrderData{UserUUID: expectedUserUUID, PartUuids: []string{partUUID}})

	s.Require().Equal(expectedUserUUID, res.UserUUID)
	s.Require().Equal(expectedPartUUIDs, res.PartUuids)
	s.Require().Equal(100.0, res.TotalPrice)
	s.Require().Equal(string(model.StatusPendingPayment), res.Status)
	s.Require().NotEmpty(res.OrderUUID)
}

func (s *ServiceSuite) TestFailureCreate_EmptyUser() {
	_, err := s.service.Create(testCtx(), model.CreateOrderData{})
	s.requireInvalidArgument(err, "user_uuid is required")
}

func (s *ServiceSuite) TestFailureCreate_NoParts() {
	_, err := s.service.Create(testCtx(), model.CreateOrderData{UserUUID: gofakeit.UUID(), PartUuids: []string{}})
	s.requireInvalidArgument(err, "order has no parts")
}

func (s *ServiceSuite) TestFailureCreate_EmptyPartUUID() {
	userUUID := gofakeit.UUID()
	_, err := s.service.Create(testCtx(), model.CreateOrderData{
		UserUUID:  userUUID,
		PartUuids: []string{""},
	})
	s.requireInvalidArgument(err, "part_uuid must not be empty")
}

func (s *ServiceSuite) TestFailureCreate_InventoryError() {
	userUUID := gofakeit.UUID()
	partUUID := gofakeit.UUID()
	partUUIDs := []string{partUUID}

	serviceErr := gofakeit.Error()

	s.inventoryClient.On("ListParts", testCtx(), model.PartsFilter{Uuids: partUUIDs}).Return(make([]model.Part, 0), serviceErr)

	_, err := s.service.Create(testCtx(), model.CreateOrderData{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
}

func (s *ServiceSuite) TestFailureCreate_PartNotFound() {
	userUUID := gofakeit.UUID()
	partUUID := gofakeit.UUID()
	partUUIDs := []string{partUUID}

	serviceErr := &model.InvalidArgumentError{Message: fmt.Sprintf("part not found: %s", partUUID)}

	s.inventoryClient.On("ListParts", testCtx(), model.PartsFilter{Uuids: partUUIDs}).Return([]model.Part{}, serviceErr)

	_, err := s.service.Create(testCtx(), model.CreateOrderData{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
}

func (s *ServiceSuite) TestFailureCreate_RepositoryError() {
	userUUID := gofakeit.UUID()
	partUUID := gofakeit.UUID()
	partUUIDs := []string{partUUID}
	serviceErr := gofakeit.Error()

	s.inventoryClient.On("ListParts", testCtx(), model.PartsFilter{Uuids: partUUIDs}).Return([]model.Part{{UUID: partUUID, Info: model.PartInfo{Price: 345}}}, nil)
	s.orderRepository.On("CreateOrder", testCtx(), mock.MatchedBy(func(o model.Order) bool {
		return o.UserUUID == userUUID && o.TotalPrice == 345
	})).Return(serviceErr)

	_, err := s.service.Create(testCtx(), model.CreateOrderData{
		UserUUID:  userUUID,
		PartUuids: partUUIDs,
	})

	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
}
