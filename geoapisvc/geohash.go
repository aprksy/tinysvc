package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	base32 = "0123456789bcdefghjkmnpqrstuvwxyz"

	msgInvalidLatLng  = "invalid lat/lng value"
	msgInvalidGeohash = "invalid geohash"
)

type LatLng struct {
	Lat float64
	Lng float64
}

type GeohashBound struct {
	SouthWest LatLng
	NorthEast LatLng
}

func GeohashEncode(lat, lng float64, precision int) (string, error) {
	prec := 12
	index := 0
	bit := 0
	evenBit := true
	result := ""
	latMin, latMax := -90.0, 90.0
	lngMin, lngMax := -180.0, 180.0

	if precision > 0 || precision < 13 {
		prec = precision
	}

	if lat < latMin || lat > latMax {
		return "", errors.New(msgInvalidLatLng)
	}
	if lng < lngMin || lng > lngMax {
		return "", errors.New(msgInvalidLatLng)
	}

	for len(result) < prec {
		if evenBit {
			lngMid := (lngMin + lngMax) / 2
			if lng >= lngMid {
				index = index*2 + 1
				lngMin = lngMid
			} else {
				index = index * 2
				lngMax = lngMid
			}
		} else {
			latMid := (latMin + latMax) / 2
			if lat >= latMid {
				index = index*2 + 1
				latMin = latMid
			} else {
				index = index * 2
				latMax = latMid
			}
		}
		evenBit = !evenBit

		bit++
		if bit == 5 {
			result = result + string([]rune(base32)[index])
			bit = 0
			index = 0
		}
	}

	return result, nil
}

func GeohashDecode(hash string) (LatLng, error) {
	result := LatLng{}
	bound, err := GeohashGetBound(hash)

	if err != nil {
		return result, err
	}

	latMin := bound.SouthWest.Lat
	lngMin := bound.SouthWest.Lng
	latMax := bound.NorthEast.Lat
	lngMax := bound.NorthEast.Lng

	lat := (latMin + latMax) / 2
	lng := (lngMin + lngMax) / 2

	result.Lat = toFixed(lat, 5)
	result.Lng = toFixed(lng, 5)
	return result, nil
}

func GeohashGetBound(hash string) (GeohashBound, error) {
	result := GeohashBound{}
	err := errors.New(msgInvalidGeohash)
	if len(hash) == 0 {
		return result, err
	}

	ghash := []byte(strings.ToLower(hash))
	evenBit := true
	latMin, latMax := -90.0, 90.0
	lngMin, lngMax := -180.0, 180.0

	for i := 0; i < len(ghash); i++ {
		chr := string(ghash[i])
		index := strings.Index(base32, chr)

		if index == -1 {
			return result, err
		}

		for n := 4; n >= 0; n-- {
			bitN := index >> n & 1
			if evenBit {
				lngMid := (lngMin + lngMax) / 2
				if bitN == 1 {
					lngMin = lngMid
				} else {
					lngMax = lngMid
				}
			} else {
				latMid := (latMin + latMax) / 2
				if bitN == 1 {
					latMin = latMid
				} else {
					latMax = latMid
				}
			}
			evenBit = !evenBit
		}
	}

	result.SouthWest.Lat = latMin
	result.SouthWest.Lng = lngMin
	result.NorthEast.Lat = latMax
	result.NorthEast.Lng = lngMax

	return result, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

/*
Geohash9 2x2 & 4x4 (usecase 9)
*/

var (
	suffix2x2 = []string{"0123", "4567", "89bc", "defg", "hjkm", "npqr", "stuv", "wxyz"}
	suffix4x4 = []string{"0123456789bcdefg", "hjkmnpqrstuvwxyz"}
)

func Geohash9_ext(hash string) (t2x2, t4x4 string, t2x2Center, t4x4Center LatLng, err error) {
	base := hash[:8]
	sfx := hash[8:9]
	sfxIndex := -1

	// t4x4
	for i, suffixStr := range suffix4x4 {
		if strings.Contains(suffixStr, sfx) {
			sfxIndex = i
			break
		}
	}
	t4x4 = fmt.Sprintf("%s-%d", base, sfxIndex)

	// find southwest
	SW := suffix4x4[sfxIndex][3:4]
	bound, _ := GeohashGetBound(base + SW)
	t4x4Center = bound.NorthEast

	// t2x2
	for i, suffixStr := range suffix2x2 {
		if strings.Contains(suffixStr, sfx) {
			sfxIndex = i
			break
		}
	}
	t2x2 = fmt.Sprintf("%s%d", t4x4, sfxIndex%4)

	// find southwest
	SW = suffix2x2[sfxIndex][:1]
	bound, _ = GeohashGetBound(base + SW)
	t2x2Center = bound.NorthEast
	return
}
