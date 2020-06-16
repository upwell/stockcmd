package cmd

import (
	"os"
	"sort"
	"sync"

	"hehan.net/my/stockcmd/logger"

	"github.com/olekukonko/tablewriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/stat"
	"hehan.net/my/stockcmd/store"
)

var ShowCmd = &cobra.Command{
	Use:     "show [group]",
	Short:   "Show stocks stat of group",
	Example: "show mygroup",
	Args:    cobra.MinimumNArgs(1),
	RunE:    showCmdF,
}

func showCmdF(cmd *cobra.Command, args []string) error {
	gName := args[0]
	g := store.GetGroup(gName)
	if g == nil {
		return errors.Errorf("Group [%s] not exist", gName)
	}

	retries := 0
	rets := make([]*stat.DailyStat, 0, 32)
	for retries < 3 {
		var wg sync.WaitGroup
		var statErr error
		for code, _ := range g.Codes {
			wg.Add(1)
			go func(code string) {
				defer wg.Done()
				ds, err := stat.GetDailyState(code)
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

	sort.Slice(rets, func(i, j int) bool {
		return rets[i].ChgToday > rets[j].ChgToday
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetReflowDuringAutoWrap(false)
	table.SetHeader(stat.Fields(stat.DailyStat{}))

	for _, ds := range rets {
		table.Append(ds.Row())
	}
	table.Render()

	return nil
}

func init() {
	rootCmd.AddCommand(ShowCmd)
}
