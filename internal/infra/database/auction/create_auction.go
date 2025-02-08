package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

func (ar *AuctionRepository) UpdateStatusAuction(
	ctx context.Context,
	id string,
	status auction_entity.AuctionStatus) *internal_error.InternalError {
	_, err := ar.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		logger.Error("Error trying to update auction", err)
		return internal_error.NewInternalServerError("Error trying to update auction")
	}
	return nil
}

func MonitorExpiredAuctions(ctx context.Context, database *mongo.Database) {
	for {
		time.Sleep(5 * time.Second)
		auctionRepo := NewAuctionRepository(database)
		getAllAuctions, err := auctionRepo.FindAuctions(ctx, auction_entity.AuctionStatus(0), "", "")
		if err != nil {
			logger.Error("Error trying to find auctions in monite Expires", err)
			return
		}

		for _, auction := range getAllAuctions {
			if auction.VerifyAuctionExpires() {
				fmt.Println("Auction: ", auction.Status, auction.ProductName)
				auctionRepo.UpdateStatusAuction(ctx, auction.Id, auction_entity.Completed)
			}
		}
	}
}
