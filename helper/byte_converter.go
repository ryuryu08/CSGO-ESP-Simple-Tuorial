package helper

import (
	"bytes"
	"encoding/binary"
	"math"
)

func Float32fromBytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Float32SliceFromBytes(data []byte) []float32 {
	bytesBuffer := bytes.NewBuffer(data)
	var array [16]float32
	err := binary.Read(bytesBuffer, binary.LittleEndian, &array)
	if err != nil {
		panic(err.Error())
	}
	return array[:]
}