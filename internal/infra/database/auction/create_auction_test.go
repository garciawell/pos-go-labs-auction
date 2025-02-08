package auction

import (
	"context"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMonitorExpiredAuctions(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should successfully update expired auctions", func(mt *mtest.T) {
		// Criando um mock de cursor com um leilão ativo
		firstBatch := mtest.CreateCursorResponse(1, "auctions.auctions", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "auction1"},
			{Key: "product_name", Value: "Product 1"},
			{Key: "category", Value: "Category 1"},
			{Key: "description", Value: "Description 1"},
			{Key: "condition", Value: auction_entity.New},
			{Key: "status", Value: auction_entity.Active},
			{Key: "timestamp", Value: time.Now().Unix()},
		})

		// Resposta de sucesso para UpdateOne (atualização do status do leilão)
		updateSuccess := mtest.CreateSuccessResponse()

		// Adicionando as respostas mockadas
		mt.AddMockResponses(firstBatch, updateSuccess)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		go MonitorExpiredAuctions(ctx, mt.DB)

		time.Sleep(6 * time.Second)

		// Criando um mock para verificar se o status foi atualizado
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "auctions.auctions", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: "auction1"},
			{Key: "status", Value: auction_entity.Completed},
		}))

		// Verificando se o leilão foi atualizado corretamente
		updatedAuction := mt.Coll.FindOne(ctx, bson.M{"_id": "auction1"})
		var result AuctionEntityMongo
		err := updatedAuction.Decode(&result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Status != auction_entity.Completed {
			t.Errorf("expected status %v, got %v", auction_entity.Completed, result.Status)
		}
	})
}
