package auction_entity

import (
	"os"
	"testing"
	"time"

	"github.com/garciawell/labs-auction-expert/internal/entity/auction_entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVerifyAuctionExpires(t *testing.T) {
	// Set up environment variable
	os.Setenv("TIME_AUCTION", "1h")
	defer os.Unsetenv("TIME_AUCTION")

	tests := []struct {
		name     string
		auction  auction_entity.Auction
		expected bool
	}{
		{
			name: "Auction not expired",
			auction: auction_entity.Auction{
				Id:          uuid.New().String(),
				ProductName: "Product 1",
				Category:    "Category 1",
				Description: "Description 1",
				Condition:   auction_entity.New,
				Status:      auction_entity.Active,
				Timestamp:   time.Now().Add(-30 * time.Minute),
			},
			expected: false,
		},
		{
			name: "Auction expired",
			auction: auction_entity.Auction{
				Id:          uuid.New().String(),
				ProductName: "Product 2",
				Category:    "Category 2",
				Description: "Description 2",
				Condition:   auction_entity.Used,
				Status:      auction_entity.Active,
				Timestamp:   time.Now().Add(-2 * time.Hour),
			},
			expected: true,
		},
		{
			name: "Auction not active",
			auction: auction_entity.Auction{
				Id:          uuid.New().String(),
				ProductName: "Product 3",
				Category:    "Category 3",
				Description: "Description 3",
				Condition:   auction_entity.Refurbished,
				Status:      auction_entity.Completed,
				Timestamp:   time.Now().Add(-2 * time.Hour),
			},
			expected: false,
		},
		{
			name: "Invalid duration",
			auction: auction_entity.Auction{
				Id:          uuid.New().String(),
				ProductName: "Product 4",
				Category:    "Category 4",
				Description: "Description 4",
				Condition:   auction_entity.New,
				Status:      auction_entity.Active,
				Timestamp:   time.Now().Add(-2 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Invalid duration" {
				os.Setenv("TIME_AUCTION", "invalid")
			} else {
				os.Setenv("TIME_AUCTION", "1h")
			}
			result := tt.auction.VerifyAuctionExpires()
			assert.Equal(t, tt.expected, result)
		})
	}
}
