package rka

import (
	"errors"
	"github.com/karalabe/hid"
)

type Device struct {
	info hid.DeviceInfo
	dev  *hid.Device

	Name string
	PID  uint16
}

func (device *Device) Open() error {
	if device.dev != nil {
		return errors.New("roccat.device: Device handle is not null, but Open() called")
	}

	var err error
	device.dev, err = device.info.Open()

	return err
}

// Sends a feature report to the device.
func (device *Device) Send(data []byte) (int, error) {
	if device.dev == nil {
		return -1, errors.New("roccat.device: cannot Send() to Device with null handle")
	}

	return device.dev.SendFeatureReport(data)
}

func (device *Device) Close() error {
	if device.dev == nil {
		return errors.New("roccat.device: Close() called for Device with null handle")
	}

	err := device.dev.Close()
	device.dev = nil

	return err
}

// Find all ROCCAT (vid 0x1e7d) products with the given product ID.
func FindAll(product uint16) []Device {
	deviceInfo := hid.Enumerate(0x1e7d, product)

	if len(deviceInfo) == 0 {
		return []Device{}
	}

	devices := make([]Device, len(deviceInfo))

	for i, info := range deviceInfo {
		devices[i] = Device{
			info: info,
			Name: info.Product,
			PID:  product,
		}
	}

	return devices
}
