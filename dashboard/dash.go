package dashboard

import (
	"fmt"
	"sort"
	"time"

	"github.com/cyliu0/tidash/dashboard/widgets"
	"github.com/cyliu0/tidash/pd"
	"github.com/gizak/termui"
	"github.com/sirupsen/logrus"
)

type Dash struct {
	grid                *termui.Grid
	storeIDList         []uint64
	updateInterval      time.Duration
	regionCountBarChart *widgets.StoreCountBarChart
	leaderCountBarChart *widgets.StoreCountBarChart
	storeCountPlots     map[uint64]*widgets.StoreCountPlot
}

func InitDash(updateInterval time.Duration) {
	if err := termui.Init(); err != nil {
		logrus.Fatalf("failed to initialize termui: %v", err)
	}
	defer termui.Close()

	dash := &Dash{
		updateInterval:  updateInterval,
		grid:            termui.NewGrid(),
		storeIDList:     make([]uint64, 0),
		storeCountPlots: make(map[uint64]*widgets.StoreCountPlot, 0),
	}
	dash.newStoreCountPlots()
	dash.newRegionCountBarChart()
	dash.newLeaderCountBarChart()
	dash.setupGrid()
	dash.eventLoop()
}

func (dash *Dash) eventLoop() {
	drawTicker := time.NewTicker(dash.updateInterval).C
	uiEvents := termui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(termui.Resize)
				termWidth, termHeight := payload.Width, payload.Height
				dash.grid.SetRect(0, 0, termWidth, termHeight)
				termui.Clear()
				termui.Render(dash.grid)
			}
		case <-drawTicker:
			dash.updateWidgets()
			termui.Render(dash.grid)
		}
	}
}

func (dash *Dash) updateWidgets() {
	trend, err := pd.GPDClinet.GetTrend()
	if err != nil {
		logrus.Fatalf("pd.GPDClinet.GetTrend failed, err: %v", err)
	}
	regionCountMap := make(map[uint64]int, 0)
	leaderCountMap := make(map[uint64]int, 0)
	for _, s := range trend.Stores {
		regionCountMap[s.ID] = s.RegionCount
		leaderCountMap[s.ID] = s.LeaderCount
		dash.storeCountPlots[s.ID].UpdateStore(s)
	}
	dash.regionCountBarChart.UpdateStoreCountMap(regionCountMap)
	dash.leaderCountBarChart.UpdateStoreCountMap(leaderCountMap)
}

func (dash *Dash) newStoreCountPlots() {
	trend, err := pd.GPDClinet.GetTrend()
	if err != nil {
		logrus.Fatalf("pd.GPDClinet.GetTrend failed, err: %v", err)
	}
	for _, store := range trend.Stores {
		dash.storeIDList = append(dash.storeIDList, store.ID)
		p := widgets.NewStoreCountPlot(
			dash.updateInterval,
			store.ID,
			fmt.Sprintf("Store ID: %d - State: %v Red: region count, Green: leader count", store.ID, store.StateName),
		)
		p.UpdateStore(store)
		dash.storeCountPlots[store.ID] = p
	}
	sort.Slice(dash.storeIDList, func(i, j int) bool { return dash.storeIDList[i] < dash.storeIDList[j] })
}

func (dash *Dash) newRegionCountBarChart() {
	dash.regionCountBarChart = widgets.NewStoreCountBarChart(
		dash.updateInterval,
		dash.storeIDList,
		"Stores' Region Count",
	)
}

func (dash *Dash) newLeaderCountBarChart() {
	dash.leaderCountBarChart = widgets.NewStoreCountBarChart(
		dash.updateInterval,
		dash.storeIDList,
		"Stores' Leader Count",
	)
}

func (dash *Dash) setupGrid() {
	storeCount := len(dash.storeIDList)
	gridItems := make([]interface{}, storeCount)
	i := 0
	for _, storeID := range dash.storeIDList {
		gridItems[i] = termui.NewRow(1.0/float64(storeCount), dash.storeCountPlots[storeID])
		i++
	}

	termWidth, termHeight := termui.TerminalDimensions()
	dash.grid.SetRect(0, 0, termWidth, termHeight)
	dash.grid.Set(
		termui.NewCol(3.0/12,
			termui.NewRow(1.0/2, dash.regionCountBarChart),
			termui.NewRow(1.0/2, dash.leaderCountBarChart),
		),
		termui.NewCol(9.0/12,
			gridItems...,
		),
	)
	dash.updateWidgets()
	termui.Render(dash.grid)
}
