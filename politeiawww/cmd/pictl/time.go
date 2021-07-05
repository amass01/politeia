// Copyright (c) 2020-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"time"
)

const (
	// timeFormat contains the reference time format that is used
	// throughout this CLI tool. This format is how timestamps are
	// printed when we want to print the human readable version.
	//
	// Mon Jan 2 15:04:05 2006
	timeFormat = "01/02/2006 3:04pm"

	// locationName is the name of the time zone location that is used
	// in the human readable timestamps.
	locationName = "Local"

	// userTimeFormat contains the reference time format that is expected
	// from the user when a date value is provided.
	//
	// Jan 2 2006
	userTimeFormat = "01/02/2006"
)

// timestampFromUnix converts a unix timestamp into a human readable timestamp
// string formatted according to the timeFormat global variable.
func timestampFromUnix(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	return t.Format(timeFormat)
}

func dateFromUnix(unixTime int64) string {
	t := time.Unix(unixTime, 0)
	return t.Format(userTimeFormat)
}

// unixFromTimestamp converts a human readable timestamp string formatted
// according to the timeFormat global variable into a unix timestamp.
func unixFromTimestamp(timestamp string) (int64, error) {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return 0, err
	}
	t, err := time.ParseInLocation(userTimeFormat, timestamp, location)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
