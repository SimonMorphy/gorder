package app

import "github.com/SimonMorphy/gorder/stock/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	CheckItemInStockHandler query.CheckItemInStockHandler
	GetItemHandler          query.GetItemHandler
}
