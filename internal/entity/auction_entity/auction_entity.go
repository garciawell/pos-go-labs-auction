package auction_entity

import (
	"context"
	"os"
	"time"

	"github.com/garciawell/labs-auction-expert/internal/internal_error"

	"github.com/google/uuid"
)

func CreateAuction(
	productName, category, description string,
	condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (au *Auction) Validate() *internal_error.InternalError {
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 2 ||
		len(au.Description) <= 10 && (au.Condition != New &&
			au.Condition != Refurbished &&
			au.Condition != Used) {
		return internal_error.NewBadRequestError("invalid auction object")
	}

	return nil
}

func (au *Auction) VerifyAuctionExpires() bool {
	durationAuction := os.Getenv("TIME_AUCTION")
	duration, err := time.ParseDuration(durationAuction)
	if err != nil {
		return false
	}
	if time.Now().Unix() > au.Timestamp.Add(duration).Unix() && au.Status == Active {
		return true
	}
	return false
}

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota + 1
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionEntity *Auction) *internal_error.InternalError

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]Auction, *internal_error.InternalError)

	FindAuctionById(
		ctx context.Context, id string) (*Auction, *internal_error.InternalError)

	UpdateStatusAuction(
		ctx context.Context,
		id string,
		status AuctionStatus) *internal_error.InternalError
}
