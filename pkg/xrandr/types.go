package xrandr

import "fmt"

type Resolution struct {
	Width      int  `json:"width"`
	Height     int  `json:"height"`
	Interlaced bool `json:"interlaced"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Mode xrandr output mode
type Mode struct {
	Resolution   Resolution    `json:"resolution"`
	RefreshRates []RefreshRate `json:"refresh_rates"`
}

// Monitor all the info of xrandr output
type Monitor struct {
	ID         string     `json:"id"`
	Modes      []Mode     `json:"modes"`
	Primary    bool       `json:"primary"`
	Size       Size       `json:"size"`
	Connected  bool       `json:"connected"`
	Resolution Resolution `json:"resolution"`
	Position   Position   `json:"position"`
}

// Screen all the info of xrandr screen
type Screen struct {
	No                int        `json:"no"`
	CurrentResolution Resolution `json:"current_resolution"`
	MinResolution     Resolution `json:"min_resolution"`
	MaxResolution     Resolution `json:"max_resolution"`
	Monitors          []Monitor  `json:"monitors"`
}

// RefreshRateValue refresh rate value
type RefreshRateValue float32

// RefreshRate mode refresh rate
type RefreshRate struct {
	Value     RefreshRateValue `json:"value"`
	Current   bool             `json:"current"`
	Preferred bool             `json:"preferred"`
}

// Screens slice of screens
type Screens []Screen

func (r Resolution) String() string {
	interlaced := ""
	if r.Interlaced {
		interlaced = "i"
	}
	return fmt.Sprintf("%dx%d%s", r.Width, r.Height, interlaced)
}

func (p Position) String() string {
	return fmt.Sprintf("%dx%d", p.X, p.Y)
}
