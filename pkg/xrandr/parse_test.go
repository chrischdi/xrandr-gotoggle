package xrandr

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_parseMonitorLine(t *testing.T) {

	tests := []struct {
		line    string
		want    Monitor
		wantErr bool
	}{
		{
			"eDP-1 connected primary 1920x1080+0+0 (normal left inverted right x axis y axis) 344mm x 194mm",
			Monitor{
				ID:          "eDP-1",
				Connected:   true,
				Resolution:  Resolution{Width: 1920, Height: 1080, Interlaced: false},
				Modes:       nil,
				Primary:     true,
				Size:        Size{Width: 344, Height: 194},
				Position:    Position{0, 0},
				Orientation: "normal",
			},
			false,
		},
		{
			"eDP-1 connected primary 1080x1920+0+0 right (normal left inverted right x axis y axis) 344mm x 194mm",
			Monitor{
				ID:          "eDP-1",
				Connected:   true,
				Resolution:  Resolution{Width: 1920, Height: 1080, Interlaced: false},
				Modes:       nil,
				Primary:     true,
				Size:        Size{Width: 344, Height: 194},
				Position:    Position{0, 0},
				Orientation: "right",
			},
			false,
		},
		{
			"eDP-1 connected primary 1920x1080+0+0 inverted (normal left inverted right x axis y axis) 344mm x 194mm",
			Monitor{
				ID:          "eDP-1",
				Connected:   true,
				Resolution:  Resolution{Width: 1920, Height: 1080, Interlaced: false},
				Modes:       nil,
				Primary:     true,
				Size:        Size{Width: 344, Height: 194},
				Position:    Position{0, 0},
				Orientation: "inverted",
			},
			false,
		},
		{
			`DP-1 connected 3440x1440+0+0 (normal left inverted right x axis y axis) 800mm x 335mm`,
			Monitor{
				ID:          "DP-1",
				Connected:   true,
				Resolution:  Resolution{Width: 3440, Height: 1440, Interlaced: false},
				Modes:       nil,
				Primary:     false,
				Size:        Size{Width: 800, Height: 335},
				Position:    Position{0, 0},
				Orientation: "normal",
			},
			false,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			got, err := parseMonitorLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMonitorLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("parseMonitorLine() = \n got %+v\nwant %+v", *got, tt.want)
			}
		})
	}
}
