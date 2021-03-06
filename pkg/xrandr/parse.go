package xrandr

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func parseScreens(s string) ([]Screen, error) {
	var screens []Screen

	ls := bufio.NewScanner(strings.NewReader(s))
	ls.Split(bufio.ScanLines)

	var err error
	var screen *Screen
	var monitor *Monitor
	var mode *Mode

	for ls.Scan() {
		l := ls.Text()
		if isScreenLine(l) {
			if screen != nil {
				screens = append(screens, *screen)
			}

			screen, err = parseScreenLine(l)
			if err != nil {
				return nil, err
			}

			continue
		}

		if isMonitorLine(l) {
			if monitor != nil {
				screen.Monitors = append(screen.Monitors, *monitor)
			}

			monitor, err = parseMonitorLine(l)
			if err != nil {
				return nil, err
			}

			continue
		}

		if monitor != nil && !monitor.Connected {
			continue
		}

		mode, err = parseModeLine(l)
		if err != nil {
			return nil, err
		}

		monitor.Modes = append(monitor.Modes, *mode)
	}

	screen.Monitors = append(screen.Monitors, *monitor)
	screens = append(screens, *screen)

	return screens, nil
}

func parseScreenLine(line string) (*Screen, error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "Screen") {
		return nil, fmt.Errorf("invalid screen line: %s", line)
	}

	re := regexp.MustCompile(`Screen \d+`)
	screenStr := re.FindString(line)
	if screenStr == "" {
		return nil, fmt.Errorf("unexpected screen line format: %s", line)
	}
	no, err := strconv.Atoi(strings.Split(screenStr, " ")[1])
	if err != nil {
		return nil, fmt.Errorf("error parsing screen number: %s", err)
	}

	screen := &Screen{
		No: no,
	}

	parseScreenResolution := func(s, typ string) (*Resolution, error) {
		if !strings.HasPrefix(s, typ) {
			return nil, fmt.Errorf("expected to start with %s", typ)
		}
		s = strings.Replace(s, typ, "", -1)
		s = strings.TrimSpace(s)
		resolution, err := parseResolution(s)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s screen resolution: %s", typ, err)
		}

		return resolution, nil
	}

	for _, typ := range []string{"minimum", "current", "maximum"} {
		re = regexp.MustCompile(fmt.Sprintf(`%s \d+\s*x\s*\d+`, typ))
		resStr := re.FindString(line)
		if resStr == "" {
			return nil, fmt.Errorf("%s resolution could not be found: %s", typ, line)
		}

		resolution, err := parseScreenResolution(resStr, typ)
		if err != nil {
			return nil, err
		}

		switch typ {
		case "minimum":
			screen.MinResolution = *resolution
		case "current":
			screen.CurrentResolution = *resolution
		case "maximum":
			screen.MaxResolution = *resolution
		}
	}

	return screen, nil
}

func parseMonitorLine(line string) (*Monitor, error) {
	line = strings.TrimSpace(line)
	tokens := strings.SplitN(line, " ", 2)
	if len(tokens) != 2 {
		return nil, fmt.Errorf("invalid monitor line format: %s", line)
	}

	id := tokens[0]
	monitor := Monitor{
		ID: id,
	}

	monitor.Connected = strings.Contains(line, " connected ")
	if !monitor.Connected {
		return &monitor, nil
	}

	primary := strings.Contains(line, "primary")
	re := regexp.MustCompile(`\d+mm\s*x\s*\d+mm`)
	sizeStr := re.FindString(tokens[1])
	if sizeStr == "" {
		return &monitor, nil
	}

	sizeStr = strings.Replace(sizeStr, "mm", "", 2)
	size, err := parseResolution(sizeStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing monitor size: %s", err)
	}

	re = regexp.MustCompile(`\d+\s*x\s*\d+\+\d+\+\d+`)
	resStr := re.FindString(tokens[1])
	if resStr == "" {
		return nil, fmt.Errorf("could not determine monitor resolution and position, expected WxH+X+Y: %s", line)
	}

	resolution, position, err := parseResolutionWithPosition(resStr)
	if err != nil {
		return nil, fmt.Errorf("could not determine monitor resolution and position: %s", err)
	}

	orientation := "normal"
	if strings.Contains(line, "left (") {
		orientation = "left"
	} else if strings.Contains(line, "right (") {
		orientation = "right"
	}
	if orientation != "normal" {
		resolution.Width, resolution.Height = resolution.Height, resolution.Width
	}

	monitor.Primary = primary
	monitor.Size = Size{Height: size.Height, Width: size.Width}
	monitor.Position = *position
	monitor.Resolution = *resolution
	monitor.Orientation = orientation

	return &monitor, nil
}

func parseModeLine(line string) (*Mode, error) {
	line = strings.TrimSpace(line)
	mode := Mode{}

	ws := bufio.NewScanner(strings.NewReader(line))
	ws.Split(bufio.ScanWords)
	for ws.Scan() {
		w := ws.Text()
		if strings.Contains(w, "x") {
			res, err := parseResolution(w)
			if err != nil {
				return nil, err
			}

			mode.Resolution = *res
			continue
		}

		if w == "+" {
			continue
		}
		rate, err := parseRefreshRate(w)
		if err != nil {
			return nil, err
		}

		mode.RefreshRates = append(mode.RefreshRates, *rate)
	}

	return &mode, nil
}

func parseResolution(s string) (*Resolution, error) {
	if !strings.Contains(s, "x") {
		return nil, fmt.Errorf("invalid size format; expected format WxH but got %s", s)
	}

	var interlaced bool
	if strings.HasSuffix(s, "i") {
		interlaced = true
		s = strings.TrimSuffix(s, "i")
	}

	res := strings.Split(s, "x")
	width, err := strconv.Atoi(strings.TrimSpace(res[0]))
	if err != nil {
		return nil, fmt.Errorf("could not parse mode width size (%s): %s", s, err)
	}

	height, err := strconv.Atoi(strings.TrimSpace(res[1]))
	if err != nil {
		return nil, fmt.Errorf("could not parse mode height size (%s): %s", s, err)
	}

	return &Resolution{
		Width:      int(width),
		Height:     int(height),
		Interlaced: interlaced,
	}, nil
}

func parseResolutionWithPosition(s string) (*Resolution, *Position, error) {
	tokens := strings.SplitN(s, "+", 2)
	size, err := parseResolution(tokens[0])
	if err != nil {
		return nil, nil, fmt.Errorf("invalid resolution with position format; expected WxH+X+Y, got %s: %s", s, err)
	}

	tokens = strings.Split(tokens[1], "+")
	if len(tokens) != 2 {
		return nil, nil, fmt.Errorf("invalid position format; expected X+Y, got %s", tokens)
	}

	x, err := strconv.Atoi(strings.TrimSpace(tokens[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid position X: %s", err)
	}

	y, err := strconv.Atoi(strings.TrimSpace(tokens[1]))
	if err != nil {
		return nil, nil, fmt.Errorf("invalid position Y: %s", err)
	}

	position := Position{
		X: x,
		Y: y,
	}

	return size, &position, nil
}

func parseRefreshRate(s string) (*RefreshRate, error) {
	s = strings.TrimSpace(s)
	current := strings.Contains(s, "*")
	preferred := strings.Contains(s, "+")

	s = strings.TrimSpace(strings.Trim(s, "*+ "))
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid rate value (%s): %s", s, err)
	}

	return &RefreshRate{
		Value:     RefreshRateValue(value),
		Current:   current,
		Preferred: preferred,
	}, nil
}

func isScreenLine(l string) bool {
	return strings.HasPrefix(l, "Screen")
}

func isMonitorLine(l string) bool {
	return strings.Contains(l, "connected") || strings.Contains(l, "disconnected")
}
