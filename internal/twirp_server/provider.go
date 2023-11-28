package twirp_server

import (
	"context"
	"log/slog"

	"github.com/charliemcelfresh/charlie-go/internal/config"

	pb "github.com/charliemcelfresh/charlie-go/rpc/charlie-go"
	"github.com/twitchtv/twirp"
)

type Repository interface {
	CreateItem(ctx context.Context, name string) (Item, error)
	GetItem(ctx context.Context, itemID string) (Item, error)
}

type provider struct {
	Repository Repository
	Logger     *slog.Logger
}

func NewProvider() provider {
	return provider{
		NewRepository(config.GetDB()),
		config.GetLogger(),
	}
}

func (p provider) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*pb.Item, error) {
	item, err := p.Repository.CreateItem(ctx, req.Name)
	if err != nil {
		return &pb.Item{}, twirp.NewError(twirp.FailedPrecondition, "cannot create item")
	}
	toReturn := &pb.Item{
		Id:        item.ID,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return toReturn, nil
}

func (p provider) GetItem(ctx context.Context, req *pb.GetItemRequest) (*pb.Item, error) {
	itemToReturn := &pb.Item{}
	item, err := p.Repository.GetItem(ctx, req.Id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return &pb.Item{}, twirp.NewError(twirp.NotFound, "item not found")
	}
	if err != nil {
		return &pb.Item{}, twirp.NewError(twirp.Internal, err.Error())
	}
	itemToReturn = &pb.Item{
		Id:        item.ID,
		Name:      item.Name,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return itemToReturn, nil
}
