package mhz19b

import (
	"errors"
	"time"

	"github.com/tarm/serial"
)

const messageLength = 9

type Sensor struct {
	device string
	stream *serial.Port
}

func New(device string) *Sensor {
	return &Sensor{
		device: device,
	}
}

func (s *Sensor) Connect(set2kDetectionRange, disableABC bool) error {
	config := &serial.Config{
		Name:        s.device,
		Baud:        9600,
		ReadTimeout: time.Second,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		return errors.Join(err, errors.New("serial.OpenPort"))
	}

	s.stream = stream

	if set2kDetectionRange || disableABC {
		// Wait for sensor readiness
		for {
			_, err = s.ReadCO2()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			break
		}

		err1 := s.Set2kDetectionRange()
		err2 := s.DisableABC()

		if err1 != nil || err2 != nil {
			return errors.Join(err1, err2)
		}
	}

	return nil
}

func (s *Sensor) ReadCO2() (uint16, error) {
	_, err := s.stream.Write([]byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79})
	if err != nil {
		return 0, errors.Join(err, errors.New("s.stream.Write"))
	}

	data, err := s.read()
	if err != nil {
		return 0, errors.Join(err, errors.New("s.read"))
	}

	return 256*uint16(data[2]) + uint16(data[3]), nil
}

func (s *Sensor) Set2kDetectionRange() error {
	// First ensure device is ready for commands
	_, err := s.ReadCO2()
	if err != nil {
		return errors.Join(err, errors.New("s.ReadCO2"))
	}

	// Set detection range 2000ppm
	var high, low uint8

	high = 2000 / 256
	low = 2000 % 256

	cmd := []byte{0xFF, 0x01, 0x99, high, low, 0x00, 0x00, 0x00, 0x00}
	cmd[len(cmd)-1] = checksum(cmd)

	_, err = s.stream.Write(cmd)
	if err != nil {
		return errors.Join(err, errors.New("s.stream.Write"))
	}

	return nil
}

// DisableABC - disables Automatic Baseline Correction (ABC logic function)
// ABC logic function refers to that sensor itself do zero point judgment and automatic calibration procedure
// intelligently after a continuous operation period. The automatic calibration cycle is every 24 hours after powered
// on. The zero point of automatic calibration is 400ppm. From July 2015, the default setting is with built-in
// automatic calibration function if no special request.
// This function is usually suitable for indoor air quality monitor such as offices, schools and homes, not suitable for
// greenhouse, farm and refrigeratory where this function should be off. Please do zero calibration timely, such as
// manual or commend calibration.
func (s *Sensor) DisableABC() error {
	// First ensure device is ready for commands
	_, err := s.ReadCO2()
	if err != nil {
		return errors.Join(err, errors.New("s.ReadCO2"))
	}

	// Disable ABC. This command has no reply
	_, err = s.stream.Write([]byte{0xFF, 0x01, 0x79, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return errors.Join(err, errors.New("s.stream.Write"))
	}

	return nil
}

func (s *Sensor) read() ([]byte, error) {
	buffer := make([]byte, messageLength)

	n, err := s.stream.Read(buffer)
	if err != nil {
		return nil, errors.Join(err, errors.New("s.stream.Read"))
	}

	err = checkMessage(buffer[:n])
	if err != nil {
		return nil, errors.Join(err, errors.New("checkMessage"))
	}

	return buffer, nil
}

func checkMessage(data []byte) error {
	if len(data) != messageLength {
		return errors.New("unexpected reply length")
	}

	if data[0] != 0xFF {
		return errors.New("unexpected 1st byte in reply")
	}

	if data[messageLength-1] != checksum(data) {
		return errors.New("bad checksum")
	}

	return nil
}

func checksum(data []byte) uint8 {
	sum := uint8(0)

	for i := 1; i < len(data)-1; i++ {
		sum += data[i]
	}

	sum = 0xFF - sum
	sum++

	return sum
}
