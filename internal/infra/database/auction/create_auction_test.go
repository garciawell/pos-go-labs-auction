package auction

import (
	"context"
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
	return []auction_entity.Auction{
		{
			Id:          "1",
			ProductName: "Test Product",
			Category:    "Test Category",
			Description: "Test Description",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now().Add(-10 * time.Minute),
		},
	}, nil
}

func (m *MockAuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	for _, auction := range m.Auctions {
		if auction.Id == id {
			return &auction, nil
		}
	}
	return nil, internal_error.NewInternalServerError("Auction not found")
}

func (m *MockAuctionRepository) UpdateStatusAuction(ctx context.Context, id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	for i, auction := range m.Auctions {
		if auction.Id == id {
			m.Auctions[i].Status = status
			return nil
		}
	}
	return internal_error.NewInternalServerError("Auction not found")
}

func TestMonitorExpiredAuctions(t *testing.T) {
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

	assert.Equal(t, auction_entity.Active, mockRepo.Auctions[0].Status)
}
