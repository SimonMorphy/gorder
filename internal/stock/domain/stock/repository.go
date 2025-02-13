package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/SimonMorphy/gorder/common/genproto/orderpb"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFountError struct {
	MissingIds []string
}

func (n NotFountError) Error() string {
	return fmt.Sprintf("not found in stock:%s ", strings.Join(n.MissingIds, ","))
}
