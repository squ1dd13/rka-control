package main

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Lights struct {
	Wheel LED

	LeftRibbon  [4]LED
	RightRibbon [4]LED

	LeftPanel  LED
	RightPanel LED

	Brightness uint8
}

type YAMLLights struct {
	KAimo struct {
		WheelHex string `yaml:"wheel"`

		LRibbonHex [4]string `yaml:"left_ribbon"`
		RRibbonHex [4]string `yaml:"right_ribbon"`

		LPanelHex string `yaml:"left_panel"`
		RPanelHex string `yaml:"right_panel"`

		Bright uint8 `yaml:"brightness"`
	} `yaml:"kone_aimo"`
}

func (yml *YAMLLights) LEDify() Lights {
	wrapper := yml.KAimo

	return Lights{
		Wheel: HexToLED(wrapper.WheelHex),

		LeftRibbon: [4]LED{
			HexToLED(wrapper.LRibbonHex[0]),
			HexToLED(wrapper.LRibbonHex[1]),
			HexToLED(wrapper.LRibbonHex[2]),
			HexToLED(wrapper.LRibbonHex[3]),
		},

		RightRibbon: [4]LED{
			HexToLED(wrapper.RRibbonHex[0]),
			HexToLED(wrapper.RRibbonHex[1]),
			HexToLED(wrapper.RRibbonHex[2]),
			HexToLED(wrapper.RRibbonHex[3]),
		},

		LeftPanel:  HexToLED(wrapper.LPanelHex),
		RightPanel: HexToLED(wrapper.RPanelHex),
		Brightness: wrapper.Bright,
	}
}

func loadLights(filePath string) Lights {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var config YAMLLights

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)

	if err != nil {
		panic(err)
	}

	return config.LEDify()
}

type KoneAIMO struct {
	device Device
	Visual Lights
}

func NewKoneAIMO(device Device) (KoneAIMO, error) {
	if device.PID != 0x2e27 {
		return KoneAIMO{}, errors.New(fmt.Sprintf(
			"kone_aimo: expected PID 0x2e27 for KoneAIMO device, got 0x%x",
			device.PID,
		))
	}

	err := device.Open()
	if err != nil {
		panic(err)
	}

	red := RGB(0xFF0000)
	redGradient := [4]LED{red, red, red, red}

	// TODO: Get the mouse's current lighting setup.
	return KoneAIMO{
		device: device,
		Visual: Lights{
			Wheel:       red,
			LeftRibbon:  redGradient,
			RightRibbon: redGradient,
			LeftPanel:   red,
			RightPanel:  red,
			Brightness:  0xFF,
		},
	}, nil
}

func (mouse *KoneAIMO) Update() error {
	_, err := mouse.device.Send(mouse.Visual.toBytes())
	return err
}

func (mouse *KoneAIMO) Close() error {
	return mouse.device.Close()
}

func (lights *Lights) toBytes() []byte {
	// This first block is almost entirely unknown bytes.
	startBlock := []byte{
		0x06, 0x69, 0x00, 0x06, 0x1F, 0x08, 0x10, 0x18, 0x20, 0x40,
		0x04, 0x03, 0x00, 0x00, 0x08, 0xFF, 0x07, 0x00, 0x01, 0x01,
		lights.Brightness, 0x1D, 0x13, 0xFF, 0x00, 0xFF, 0x59, 0xFF, 0x00, 0x00,
		0xFF, 0xFD, 0xFD, 0x00, 0x00, 0xFF, 0xF4, 0x64, 0x00, 0x00,
		0xFF, 0xF4, 0x00, 0x00, 0x00, 0xFF,
	}

	buffer := bytes.NewBuffer(startBlock)

	lights.Wheel.Write(buffer)

	for _, led := range lights.LeftRibbon {
		led.Write(buffer)
	}

	for _, led := range lights.RightRibbon {
		led.Write(buffer)
	}

	lights.LeftPanel.Write(buffer)
	lights.RightPanel.Write(buffer)

	// These do change but I'm not sure what they do yet.
	// The final two are just some sort of identifier/timestamp/index,
	//  but the mouse doesn't seem to care what it is.
	buffer.Write([]byte{0x02, 0x01, 0xA5, 0x33})

	return buffer.Bytes()
}
