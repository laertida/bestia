package main

import (
	"encoding/json"
	"flag"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"lmnl/bestia/clients"
	"lmnl/bestia/helpers"
	"lmnl/bestia/lights"
	"log/slog"
	"os"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"time"
)

const (
	TRIG_PORT      = "23"
	ECHO_PORT      = "24"
	TEST_PORT      = "17"
	LEFT_ADDR      = 0x23
	RIGH_ADDR      = 0x27
	IODIR          = 0x00
	GPIO           = 0x09
	ADDR           = 0x02
	MQTT_CLIENT_ID = "canon.v2"
	MQTT_TOPIC     = "lmnl.mx/canon.v2/lights"
	MQTT_BROKER    = "tcp://192.168.1.68:1883"
)

func main() {

	var lightsOn = flag.Bool("lightsOn", false, "This flags allows to configure if lights will be on or off.")
	var logLevel = flag.String("logLevel", "error", "This flag determinates which is the log level, valid values are debug, info, warn and error.")
	var mqttClientId = flag.String("mqttClient", MQTT_CLIENT_ID, "This flag determinates the client id for mqtt broker connection.")
	var mqttTopic = flag.String("mqttTopic", MQTT_TOPIC, "This flag determinates the topic to suscribe in the mqtt broker.")
	var mqttBroker = flag.String("mqttBroker", MQTT_BROKER, "This flag determinates the url of mqtt broker.")

	flag.Parse()

	level := helpers.GetLogLevel(*logLevel)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	slog.Info("Bestia started with this initial values", "lightsOn", *lightsOn, "mqttClientId", *mqttClientId, "mqttTopic", *mqttTopic, "mqttBroker", *mqttBroker, "logLevel", *logLevel)

	messages := make(chan mqtt.Message)

	// L 0x40 0x20 0x10 0x08 0x04 0x02
	// R 0x02 0x04 0x08 0x10 0x20 0x40
	R_MAP := map[int]int{1: 0x02, 2: 0x04, 3: 0x08, 4: 0x10, 5: 0x20, 6: 0x40}

	// Initialize periph.io
	if _, err := host.Init(); err != nil {
		fmt.Printf("Failed to initialize periph: %v \n", err)
		os.Exit(1)
		return
	}
	bus, err := i2creg.Open("")
	if err != nil {
		fmt.Printf("failed to open I2C bus: %v \n", err)
		//os.Exit(1)
	}

	defer bus.Close()

	right := lights.NewLights("right side", RIGH_ADDR, R_MAP, bus, *logger)

	testMode := gpioreg.ByName("17")
	test := false

	if testMode.Read() == gpio.High {
		test = true
	}

	opts := mqtt.NewClientOptions().
		AddBroker(*mqttBroker).
		SetClientID(*mqttClientId).
		SetKeepAlive(60 * time.Second).
		SetPingTimeout(10 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		if token := client.Subscribe(*mqttTopic, 0, func(_ mqtt.Client, msg mqtt.Message) {
			messages <- msg
		}); token.Wait() && token.Error() != nil {
			logger.Error("Error al suscribirse.", "error", token.Error())
		} else {
			logger.Info("Suscripción correcta.", "broker", *mqttBroker, "topic", *mqttTopic, "clientId", *mqttClientId)
		}
	}()

	// Main goroutine: consume distance and do something
	if test {
		go func() {
			for msg := range messages {

				message := clients.Message{}

				err := json.Unmarshal(msg.Payload(), &message)

				if err != nil {
					logger.Error("There was an error unmarshaling message")
				}
				logger.Debug("Mensaje recibido", "broker", *mqttBroker, "topic", *mqttTopic, "clientId", *mqttClientId, "message", message)

				if message.R1 == "on" {
					right.On(1)
				} else {
					right.Off(1)
				}

				if message.R2 == "on" {
					right.On(2)
				} else {
					right.Off(2)
				}

				if message.R3 == "on" {
					right.On(3)
				} else {
					right.Off(3)
				}

				if message.R4 == "on" {
					right.On(4)
				} else {
					right.Off(4)
				}

				if message.R5 == "on" {
					right.On(5)
				} else {
					right.Off(5)
				}

				if message.R6 == "on" {
					right.On(6)
				} else {
					right.Off(6)
				}
			}
		}()

	} else {
		go func() {

			for {

				if test {
					right.AllOn()
				} else {
					right.Step()
				}
				time.Sleep(1 * time.Second)
			}

		}()

	}

	select {}
}
