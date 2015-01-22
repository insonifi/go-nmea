package nmea

import (
	"nmea"
	"testing"
)

const (
	NMEA_Valid         = "$GPRMC,124617.00,A,5657.27002,N,02403.80194,E,0.020,51.63,301212,212*2e"
	NMEA_NoPrefix      = "GPRMC,124617.00,A,5657.27002,N,02403.80194,E,0.020,51.63,301212,212*2e"
	NMEA_InvalidPrefix = "!GPRMC,124617.00,A,5657.27002,N,02403.80194,E,0.020,51.63,301212,212*2e"
	NMEA_ChecksumErr   = "$GPRMC,124617.00,A,5657.27002,N,02403.80194,E,0.020,51.63,301212,212*22"
	NMEA_TimeErr       = "$GPRMC,000000,A,5657.27002,N,02403.80194,E,0.020,51.63,000000,212*2e"
	TimeString         = "2012-12-30 12:46:17 +0000 UTC"
	Lat                = 56.95450033333524
	Lng                = 24.063365666666794
	Kph                = 0.037039984
	Track              = 51.63
)

func TestChecksum(t *testing.T) {
	_, err := nmea.Parse(NMEA_ChecksumErr)
	if err == nil {
		t.Error("Checksum incorrect")
	}
}
func TestTimeErr(t *testing.T) {
	_, err := nmea.Parse(NMEA_ChecksumErr)
	if err == nil {
		t.Error("Empty time not detected")
	}
}
func TestParse(t *testing.T) {
	n, err := nmea.Parse(NMEA_Valid)
	if err != nil {
		t.Error(err)
	}
	if n.Coords.Lat != Lat || n.Coords.Lng != Lng {
		t.Error("Coordinates incorrect")
	}
	if n.Kph != Kph {
		t.Error("Speed incorrect")
	}
	if n.Track != Track {
		t.Error("Track incorrect")
	}
	if n.Timestamp.String() != TimeString {
		t.Error("Timestamp incorrect")
	}
}
func TestPrefix(t *testing.T) {
	_, err := nmea.Parse(NMEA_NoPrefix)
	if err == nil {
		t.Error("Absent prefix not detected")
	}
	_, err = nmea.Parse(NMEA_InvalidPrefix)
	if err == nil {
		t.Error("Invalid prefix not detected")
	}
}
