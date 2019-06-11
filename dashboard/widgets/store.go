package widgets

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cyliu0/tidash/pd"
	"github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

type StoreCountBarChart struct {
	BarChart
	storeIDList   []uint64
	StoreCountMap map[uint64]int
}

func NewStoreCountBarChart(updateInterval time.Duration, storeIDList []uint64, title string) *StoreCountBarChart {
	scbc := &StoreCountBarChart{
		storeIDList: storeIDList,
		BarChart: BarChart{
			BarChart: widgets.NewBarChart(),
			Widget: Widget{
				updateInterval: updateInterval,
			},
		},
	}
	storeCount := len(storeIDList)
	scbc.Labels = make([]string, storeCount)
	scbc.Data = make([]float64, storeCount)
	scbc.StoreCountMap = make(map[uint64]int)
	scbc.Title = title
	//scbc.BarColors = []termui.Color{termui.ColorMagenta}
	for i, storeID := range storeIDList {
		scbc.Labels[i] = strconv.FormatUint(storeID, 10)
	}

	scbc.update()
	go func() {
		for range time.NewTicker(scbc.updateInterval).C {
			scbc.update()
		}
	}()
	return scbc
}

func (scbc *StoreCountBarChart) update() {
	for i, storeID := range scbc.storeIDList {
		scbc.Data[i] = float64(scbc.StoreCountMap[storeID])
	}
}

func (scbc *StoreCountBarChart) UpdateStoreCountMap(storeCountMap map[uint64]int) {
	scbc.StoreCountMap = storeCountMap
	scbc.update()
}

type StoreCountPlot struct {
	Plot
	storeID uint64
	store   *pd.TrendStore
}

func NewStoreCountPlot(updateInterval time.Duration, storeID uint64, title string) *StoreCountPlot {
	scp := &StoreCountPlot{
		storeID: storeID,
		Plot: Plot{
			Plot: widgets.NewPlot(),
			Widget: Widget{
				updateInterval: updateInterval,
			},
		},
	}
	scp.LineColors = []termui.Color{termui.ColorRed, termui.ColorGreen}
	scp.DrawDirection = widgets.DrawLeft
	scp.Title = title
	scp.Data = make([][]float64, 2)
	scp.Data[0] = make([]float64, 0)
	scp.Data[1] = make([]float64, 0)

	scp.update()
	go func() {
		for range time.NewTicker(scp.updateInterval).C {
			scp.update()
		}
	}()
	return scp
}

func (scp *StoreCountPlot) update() {
	if scp.store != nil {
		scp.Data[0] = append(scp.Data[0], float64(scp.store.RegionCount))
		if len(scp.Data[0]) > 125 {
			scp.Data[0] = scp.Data[0][1:]
		}
		scp.Data[1] = append(scp.Data[1], float64(scp.store.LeaderCount))
		if len(scp.Data[1]) > 125 {
			scp.Data[1] = scp.Data[1][1:]
		}
		scp.Title = fmt.Sprintf("Store ID: %d - State: %v, Region Count(red): %v, Leader Count(green): %v", scp.store.ID, scp.store.StateName, scp.store.RegionCount, scp.store.LeaderCount)
	}
}

func (scp *StoreCountPlot) UpdateStore(store pd.TrendStore) {
	scp.store = &store
	scp.update()
}
