package cmd

import (
	"os"
	"sort"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/olekukonko/tablewriter"
	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/stat"
	"hehan.net/my/stockcmd/task"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/store"
)

var ShowCmd = &cobra.Command{
	Use:     "show [group]",
	Short:   "Show stocks stat of group",
	Example: "show mygroup",
	Args:    cobra.MinimumNArgs(1),
	RunE:    showCmdF,
}

var showPeriodVar int
var intervalVar int

// showOnce wrapper for show daily stat once
func showOnce(codes []string) {
	retries := 0
	rets := make([]*stat.DailyStat, 0, 32)
	for retries < 3 {
		var wg sync.WaitGroup
		var statErr error
		for _, code := range codes {
			wg.Add(1)
			go func(code string) {
				defer wg.Done()
				ds, err := stat.GetDailyState(code, showPeriodVar)
				if err == nil {
					rets = append(rets, ds)
				} else {
					if errors.Is(err, store.ErrDBColNotMatch) {
						statErr = err
						logger.SugarLog.Infof("db fields changed, clean history data and get again")
					} else {
						logger.SugarLog.Errorf("get daily state error [%v]", err)
					}
				}
			}(code)
		}
		wg.Wait()

		retries++
		if statErr != nil {
			store.RecreateDailyBucket()
		} else {
			break
		}
	}

	//fmt.Printf("[%s] since t2\n", time.Since(t2))
	//t3 := time.Now()

	sort.Slice(rets, func(i, j int) bool {
		return rets[i].ChgToday > rets[j].ChgToday
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetReflowDuringAutoWrap(false)
	fieldNames := stat.Fields(stat.DailyStat{})

	excludeFields := mapset.NewSet("avg_200", "PB")
	filterNames := make([]string, 0)
	for _, field := range fieldNames {
		if !excludeFields.Contains(field) {
			filterNames = append(filterNames, field)
		}
	}
	table.SetHeader(filterNames)

	for _, ds := range rets {
		table.Append(ds.Row())
	}
	table.Render()
}

func showCmdF(cmd *cobra.Command, args []string) error {
	//t1 := time.Now()

	gName := args[0]
	g := store.GetGroup(gName)
	codes := make([]string, 0, 32)
	if g == nil {
		if gName == "all" {
			codeSet := store.GetAllStockCodes()
			for code := range codeSet.Iter() {
				codes = append(codes, code.(string))
			}
		} else {
			return errors.Errorf("Group [%s] not exist", gName)
		}
	} else {
		for code, _ := range g.Codes {
			codes = append(codes, code)
		}
	}

	task.CheckAllStockDividendDay(false)

	if intervalVar > 0 {
		timer := time.NewTicker(time.Duration(intervalVar) * time.Second)

		for {
			select {
			case <-timer.C:
				showOnce(codes)
			}
		}
	} else {
		showOnce(codes)
	}

	//startDBStat := store.DB.Stats()
	//fmt.Printf("[%s] since t1\n", time.Since(t1))
	//t2 := time.Now()

	//fmt.Printf("[%s] since t3\n", time.Since(t3))

	//endDBStat := store.DB.Stats()
	//diffDBStat := endDBStat.Sub(&startDBStat)
	//json.NewEncoder(os.Stderr).Encode(diffDBStat)

	return nil
}

func init() {
	ShowCmd.Flags().IntVarP(&showPeriodVar, "period", "p", 120, "get the <period> days of stat")
	ShowCmd.Flags().IntVarP(&intervalVar, "interval", "i", 0, "when set to greater than 0,"+
		" it will continue update the daily stat periodically by this `interval`")
	rootCmd.AddCommand(ShowCmd)
}
