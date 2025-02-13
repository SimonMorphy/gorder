package models

import "fmt"

type RabbitMQVoucher struct {
	Host     string
	Port     string
	User     string
	Password string
}

func (v RabbitMQVoucher) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s", v.User, v.Password, v.Host, v.Port)
}
