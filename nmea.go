package nmea

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	EmptyTime  = "000000"
	OneMin     = 0.0166666666667
	Knot       = 1.8519993258722454
	TimeLayout = "020106150405.00"
)

type GPSMessage struct {
	Coords struct {
		Lat float64
		Lng float64
	}
	Kph       float32
	Track     float32
	Timestamp time.Time
}

func checksum(N *string) (checksum uint64) {
	for _, r := range *N {
		checksum ^= uint64(r)
	}
	return
}

func convertToDecimal(C string) (decimal float64) {
	split := strings.Index(C, ".") - 2
	degrees, _ := strconv.ParseFloat(C[:split], 64)
	minutes, _ := strconv.ParseFloat(C[split:], 64)
	decimal = degrees + (minutes * OneMin)
	return
}

func Parse(N string) (Msg GPSMessage, err error) {
	if !strings.HasPrefix(N, "$") {
		err = errors.New("Invalid NMEA sentence")
		return
	}
	data_str := N[1 : len(N)-3]
	csum, _ := strconv.ParseUint(N[len(N)-2:], 16, 16)
	if csum != checksum(&data_str) {
		err = errors.New("Checksum invalid")
		return
	}
	data := strings.Split(data_str, ",")
	if data[0] != "GPRMC" {
		err = errors.New("Unsupported sentence")
		return
	}
	if data[1] == EmptyTime || data[9] == EmptyTime {
		err = errors.New("Sentence invalid. No timestamp.")
		return
	}
	Msg.Timestamp, _ = time.Parse(TimeLayout, data[9]+data[1])
	Msg.Coords.Lat = convertToDecimal(data[3])
	if data[4] == "S" {
		Msg.Coords.Lat = -Msg.Coords.Lat
	}
	Msg.Coords.Lng = convertToDecimal(data[5])
	if data[6] == "W" {
		Msg.Coords.Lng = -Msg.Coords.Lng
	}
	kph, _ := strconv.ParseFloat(data[7], 32)
	Msg.Kph = float32(kph * Knot)
	track, _ := strconv.ParseFloat(data[8], 32)
	Msg.Track = float32(track)
	return
}
