package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/garciawell/labs-auction-expert/internal/entity/auction_entity"
	"github.com/garciawell/labs-auction-expert/internal/internal_error"
	"github.com/stretchr/testify/assert"
)

type MockAuctionRepository struct {
	Auctions []auction_entity.Auction
}

func (m *MockAuctionRepository) CreateAuction(ctx context.Context, auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	m.Auctions = append(m.Auctions, *auctionEntity)
	return nil
}

func (m *MockAuctionRepository) FindAuctions(ctx context.Context, status auction_entity.AuctionStatus, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	return m.Auctions, nil
}

func (m *MockAuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	for i := range m.Auctions {
		if m.Auctions[i].Id == id {
			return &m.Auctions[i], nil //  Retorna referência correta
		}
	}
	return nil, internal_error.NewInternalServerError("Auction not found")
}

func (m *MockAuctionRepository) UpdateStatusAuction(ctx context.Context, id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	for i := range m.Auctions {
		if m.Auctions[i].Id == id {
			m.Auctions[i].Status = status //  Atualiza diretamente a referência
			return nil
		}
	}
	return internal_error.NewInternalServerError("Auction not found")
}

func TestMonitorExpiredAuctions(t *testing.T) {
	os.Setenv("TIME_AUCTION", "1m")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockRepo := &MockAuctionRepository{
		Auctions: []auction_entity.Auction{
			{
				Id:          "1",
				ProductName: "Test Product",
				Category:    "Test Category",
				Description: "Test Description",
				Condition:   auction_entity.New,
				Status:      auction_entity.Active,
				Timestamp:   time.Now().Add(-10 * time.Minute),
			},
		},
	}

	go MonitorExpiredAuctions(ctx, mockRepo)

	time.Sleep(6 * time.Second)

	assert.Equal(t, auction_entity.Completed, mockRepo.Auctions[0].Status)
}
