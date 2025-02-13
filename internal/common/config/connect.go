package config

import (
	"github.com/SimonMorphy/gorder/common/broker"
	"github.com/SimonMorphy/gorder/common/config/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var RabbitMQ *amqp091.Channel

func InitConfig() (func() error, error) {
	//初始化viper
	if err := NewViperConfig(); err != nil {
		logrus.Fatal(err)
		return func() error {
			return nil
		}, err
	}
	//初始化RabbitMQ
	var rabbitMQConfig models.RabbitMQVoucher
	if err := viper.Sub("rabbit-mq").Unmarshal(&rabbitMQConfig); err != nil {
		logrus.Fatal(err)
		return func() error {
			return nil
		}, err
	}
	conn, f2 := broker.Connect(&rabbitMQConfig)
	RabbitMQ = conn
	return func() error {
		_ = f2()
		_ = conn.Close()
		return nil
	}, nil

	//初始化MongoDB

}
