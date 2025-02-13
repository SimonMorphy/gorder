package app

import "github.com/SimonMorphy/gorder/payment/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}

type Queries struct {
}
