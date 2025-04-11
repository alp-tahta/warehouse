package service

import (
	"testing"

	"github.com/alp-tahta/warehouse/internal/barcode"
	"github.com/alp-tahta/warehouse/internal/logger"
	"github.com/alp-tahta/warehouse/internal/model"
	"github.com/alp-tahta/warehouse/internal/repository"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryI(ctrl)
	mockBarcoder := barcode.NewMockBarcoder(ctrl)
	mockLogger := logger.Init()

	service := New(mockLogger, mockBarcoder, mockRepo)

	req := model.CreateOrderRequest{}
	mockRepo.EXPECT().CreateOrder(req).Return(nil)

	err := service.CreateOrder(req)
	assert.NoError(t, err)
}

func TestUpdateBarcodeStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryI(ctrl)
	mockBarcoder := barcode.NewMockBarcoder(ctrl)
	mockLogger := logger.Init()

	service := New(mockLogger, mockBarcoder, mockRepo)

	id := "barcode123"
	mockRepo.EXPECT().CheckIfBarcodeMarked(id).Return(false, nil)
	mockBarcoder.EXPECT().ResolveBarcode(id).Return(model.BarcodeFields{}, nil)
	mockRepo.EXPECT().IncreaseShelfOccupancy(model.BarcodeFields{}).Return(nil)
	mockRepo.EXPECT().MarkBarcode(id).Return(nil)

	err := service.UpdateBarcodeStatus(id)
	assert.NoError(t, err)
}

func TestGetShelvesDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepositoryI(ctrl)
	mockBarcoder := barcode.NewMockBarcoder(ctrl)
	mockLogger := logger.Init()

	service := New(mockLogger, mockBarcoder, mockRepo)

	expectedShelves := []model.ShelfInformationWithCustomer{}
	mockRepo.EXPECT().GetShelvesDetails().Return(expectedShelves, nil)

	shelves, err := service.GetShelvesDetails()
	assert.NoError(t, err)
	assert.Equal(t, expectedShelves, shelves)
}
