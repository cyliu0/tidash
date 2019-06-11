package widgets

import (
	"time"

	"github.com/gizak/termui/widgets"
)

type Widget struct {
	updateInterval time.Duration
}

type BarChart struct {
	*widgets.BarChart
	Widget
}

type SparklineGroup struct {
	*widgets.SparklineGroup
	Widget
}

type Plot struct {
	*widgets.Plot
	Widget
}
