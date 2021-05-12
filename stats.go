package main

import (
	"math"
	"strconv"
	"time"
)

type unixTime struct {
	time.Time
}

type Stats struct {
	Start         unixTime `json:"start"`
	End           unixTime `json:"end"`
	Local         Local    `json:"local"`
	Remote        Remote   `json:"remote"`
	CPR           CPR      `json:"cpr"`
	AltSuppressed int      `json:"altitude_suppressed"`
	Tracks        Tracks   `json:"tracks"`
	Messages      int      `json:"messages"`
}

type Local struct {
	SamplesProcessed int     `json:"sample_processed"`
	SamplesDropped   int     `json:"samples_dropped"`
	ModeAC           int     `json:"modeac"`
	ModeS            int     `json:"modes"`
	Bad              int     `json:"bad"`
	UnknownICAO      int     `json:"unknown_icao"`
	Accepted         []int   `json:"accepted"`
	Signal           float32 `json:"signal"`
	Noise            float32 `json:"noise"`
	PeakSignal       float32 `json:"peak_signal"`
	StrongSignals    int     `json:"strong_signals"`
}

type Remote struct {
	ModeAC      int   `json:"modeac"`
	Modes       int   `json:"modes"`
	Bad         int   `json:"bad"`
	UnknownICAO int   `json:"unknown_icao"`
	Accepted    []int `json:"accepted"`
	Requests    int   `json:"http_requests"`
}

type CPR struct {
	Surface        int `json:"surface"`
	Airborne       int `json:"airborne"`
	GlobalOK       int `json:"global_ok"`
	GlobalBad      int `json:"global_bad"`
	GlobalBadRange int `json:"global_range"`
	GlobalBadSpeed int `json:"global_speed"`
	GlobalSkipped  int `json:"global_skipped"`
	LocalOK        int `json:"local_ok"`
	LocalAcftRel   int `json:"local_aircraft_relative"`
	LocalRecvRel   int `json:"local_receiver_relative"`
	LocalSkip      int `json:"local_skipped"`
	LocalSkipRange int `json:"local_range"`
	LocalSkipSpeed int `json:"local_speed"`
	Filtered       int `json:"filtered"`
}

type Tracks struct {
	All           int `json:"all"`
	SingleMessage int `json:"single_message"`
}

func (u *unixTime) UnmarshalJSON(data []byte) error {
	n, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}

	s := int64(math.Floor(n))
	ns := int64((n - float64(s)) * 1e9)

	u.Time = time.Unix(s, ns)
	return nil
}
