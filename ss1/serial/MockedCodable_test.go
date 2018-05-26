package serial_test

import "github.com/inkyblackness/hacked/ss1/serial"

type MockedCodable struct {
	calledCoder serial.Coder
}

func (codable *MockedCodable) Code(coder serial.Coder) {
	codable.calledCoder = coder
}
