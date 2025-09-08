package clients

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"time"
)

type MQTT struct {
	Broker   string
	ClientId string
	Topic    string
	Client   mqtt.Client
	logger   slog.Logger
}

func NewMQTTClient(broker, clientId, topic string, logger slog.Logger) MQTT {

	clientOpts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(clientId).
		SetKeepAlive(1 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			logger.Info("mensaje recibido", "msg", string(msg.Payload()))
		})

	client := mqtt.NewClient(clientOpts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("No pudo establecerse conexión con el broker.", "broker", broker, "topic", topic, "clientId", clientId, "error", token.Error)
		panic(token.Error())
	}

	return MQTT{
		Broker:   broker,
		ClientId: clientId,
		Topic:    topic,
		logger:   logger,
		Client:   client,
	}
}

type Message struct {
	R1 string `json:"1"`
	R2 string `json:"2"`
	R3 string `json:"3"`
	R4 string `json:"4"`
	R5 string `json:"5"`
	R6 string `json:"6"`
}
