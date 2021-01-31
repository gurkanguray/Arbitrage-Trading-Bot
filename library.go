package main

import (
	"fmt"
	"strconv"
)

// SingleMarketActivity is a defined interface for whitebit API
type SingleMarketActivity interface {
	// SingleMarketActivity calling methods for Result structure data
	getBid() float64
	getAsk() float64
	getOpen() float64
	getHigh() float64
	getLow() float64
	getLast() float64
	getVolume() float64
	getChange() float64
}

// EmbeddedObject exported type declaration for using functions
type EmbeddedObject map[string]interface{}

// Result is a structure declaration Result data of whitebit API
type Result struct {
	// defining struct variables
	bid    float64
	ask    float64
	open   float64
	high   float64
	low    float64
	last   float64
	volume float64
	change float64
}

// GetBid takes "bid" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetBid() float64 {
	bid := embedded["bid"].(string)
	b, err := strconv.ParseFloat(bid, 64)
	if err != nil {
		fmt.Println(b)
	}
	return b
}

// GetAsk takes "ask" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetAsk() float64 {
	ask := embedded["ask"].(string)
	a, err := strconv.ParseFloat(ask, 64)
	if err != nil {
		fmt.Println(a)
	}
	return a
}

// GetOpen takes "open" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetOpen() float64 {
	open := embedded["open"].(string)
	o, err := strconv.ParseFloat(open, 64)
	if err != nil {
		fmt.Println(o)
	}
	return o
}

// GetHigh takes "high" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetHigh() float64 {
	high := embedded["high"].(string)
	h, err := strconv.ParseFloat(high, 64)
	if err != nil {
		fmt.Println(h)
	}
	return h
}

// GetLow takes "low" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetLow() float64 {
	low := embedded["low"].(string)
	l, err := strconv.ParseFloat(low, 64)
	if err != nil {
		fmt.Println(l)
	}
	return l
}

// GetLast takes "last" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetLast() float64 {
	last := embedded["last"].(string)
	l, err := strconv.ParseFloat(last, 64)
	if err != nil {
		fmt.Println(l)
	}
	return l
}

// GetVolume takes "volume" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetVolume() float64 {
	volume := embedded["volume"].(string)
	v, err := strconv.ParseFloat(volume, 64)
	if err != nil {
		fmt.Println(v)
	}
	return v
}

// GetChange takes "change" value of the embedded variable which returns as float64
func (embedded EmbeddedObject) GetChange() float64 {
	change := embedded["change"].(string)
	c, err := strconv.ParseFloat(change, 64)
	if err != nil {
		fmt.Println(c)
	}
	return c
}