package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"hehan.net/my/stockcmd/global"

	"hehan.net/my/stockcmd/util"

	mapset "github.com/deckarep/golang-set"
	"github.com/olekukonko/tablewriter"

	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/stat"
	"hehan.net/my/stockcmd/store"
)

var StatCmd = &cobra.Command{
	Use:     "stat",
	Short:   "Show stat of stocks",
	Example: "stat",
	RunE:    statCmdF,
}

var RPSCmd = &cobra.Command{
	Use:     "rps",
	Short:   "Show rps of stocks",
	Example: "rps",
	RunE:    rpsCmdF,
}

var MyStatCmd = &cobra.Command{
	Use:     "mystat",
	Short:   "Show stat of my stocks",
	Example: "mystat",
	RunE:    myStatCmdF,
}

var FetchDataCmd = &cobra.Command{
	Use:     "fetchData",
	Short:   "fetch hq and record",
	Example: "fetchData",
	RunE:    fetchDataCmdF,
}

type StatChg struct {
	Code   string
	Name   string
	MaxChg float64
	MinChg float64
}

var periodVar int
var includeST bool
var groupVar string
var sortByIncrease bool

var daysVar int
var showNumVar int

func getStatChgs(basics []*store.StockBasic) []*StatChg {
	chgs := make([]*StatChg, 0, 512)
	endDay := now.BeginningOfDay()
	var wg sync.WaitGroup
	startDay := endDay.AddDate(0, 0, -150)
	for _, basic := range basics {
		wg.Add(1)
		go func(basic *store.StockBasic) {
			defer wg.Done()

			if !includeST && (strings.Contains(basic.Name, "ST") || strings.Contains(basic.Name, "退")) {
				return
			}

			code := basic.Code
			sChg := &StatChg{
				Code: basic.Code,
				Name: basic.Name,
			}
			df, err := store.GetRecords(code, startDay, endDay)
			if err != nil {
				logger.SugarLog.Errorf("get records for [%s] error [%v]", code, err)
				return
			}
			if df.NRows() == 0 {
				logger.SugarLog.Errorf("get records return zero rows for [%s]", code)
				return
			}
			max, min := stat.GetMaxMin(df, periodVar)
			if max == 0.00 || min == 0.00 {
				logger.SugarLog.Infof("failed to get max min for [%s]", code)
				return
			}
			price := store.GetHQ(code)
			if price == 0.00 {
				logger.SugarLog.Infof("failed to get price for [%s]", code)
				return
			}

			sChg.MaxChg = util.RoundChgRate((price - max) / max)
			sChg.MinChg = util.RoundChgRate((price - min) / min)
			if sChg.MaxChg == -100.0 || sChg.MinChg == -100.0 {
				return
			}
			chgs = append(chgs, sChg)
		}(basic)
	}
	wg.Wait()

	if sortByIncrease {
		sort.SliceStable(chgs, func(i, j int) bool {
			return chgs[i].MinChg > chgs[j].MinChg
		})
	} else {
		sort.SliceStable(chgs, func(i, j int) bool {
			return chgs[i].MaxChg < chgs[j].MaxChg
		})
	}
	return chgs
}

func statCmdF(cmd *cobra.Command, args []string) error {
	basics := store.GetBasics()
	chgs := getStatChgs(basics)
	printTable(chgs)
	return nil
}

func rpsCmdF(cmd *cobra.Command, args []string) error {
	basics := store.GetBasics()
	rpss := stat.GetRPS(basics, daysVar)
	printTableRPS(rpss, showNumVar)
	return nil
}

func myStatCmdF(cmd *cobra.Command, args []string) error {
	groupNames := store.ListGroup()
	basics := make([]*store.StockBasic, 0, 32)
	codeSet := mapset.NewSet()
	for _, name := range groupNames {
		if groupVar != "" && name != groupVar {
			continue
		}

		group := store.GetGroup(name)
		for code, name := range group.Codes {
			if !codeSet.Contains(code) {
				basics = append(basics, &store.StockBasic{
					Code: code,
					Name: name,
				})
				codeSet.Add(code)
			}
		}

		if groupVar == name {
			break
		}
	}

	chgs := getStatChgs(basics)
	printTable(chgs)
	return nil
}

func printTable(chgs []*StatChg) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetReflowDuringAutoWrap(false)

	table.SetHeader([]string{"Code", "Name", "ChgMax", "ChgMin"})

	for idx, chg := range chgs {
		if idx > 50 {
			break
		}
		table.Append([]string{chg.Code, chg.Name, util.Float64String(chg.MaxChg), util.Float64String(chg.MinChg)})
	}
	table.Render()
}

func printTableRPS(rpss []*stat.RPS, numVar int) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetReflowDuringAutoWrap(false)

	table.SetHeader([]string{"Code", "Name", "RPS", "Chg"})

	for idx, rps := range rpss {
		if idx > numVar {
			break
		}
		table.Append([]string{rps.Code, rps.Name, util.Float64String(rps.Value), util.Float64String(rps.Change)})
	}
	table.Render()
}

func fetchDataCmdF(cmd *cobra.Command, args []string) error {
	codes := store.GetCodes()

	fmt.Printf("fetch hq data ...\n")
	start := time.Now()
	var wg sync.WaitGroup
	hqs := make([]*store.StockHQ, 0, 512)
	for _, code := range codes {
		wg.Add(1)
		api := global.GetHQSource()
		go func(code string) {
			defer wg.Done()
			v, err := api.GetHQ(code)
			if err != nil {
				fmt.Printf("failed to get price for [%s] with error [%v]\n", code, err)
				return
			}
			if v.IsSuspend {
				fmt.Printf("%s is suspend\n", code)
			}
			if v.Now == 0.00 && v.Last == 0.00 {
				fmt.Printf("%s now and last is zero\n", code)
			}
			hq := &store.StockHQ{
				Code:  code,
				Price: fmt.Sprintf("%f", v.Now),
			}
			hqs = append(hqs, hq)
		}(code)
	}

	wg.Wait()
	store.BulkWriteHQ(hqs)
	fmt.Printf("fetch hq data done, take [%s], start fetch history records ... \n", time.Since(start))

	// use channel to control the number of concurrent fetch tasks
	fetchCh := make(chan string, 64)
	go func() {
		for {
			code, ok := <-fetchCh
			if !ok {
				break
			}

			go func(code string) {
				_, err := stat.GetDataFrame(global.GetDataSource(), code)
				if err != nil {
					logger.SugarLog.Error(err)
					fetchCh <- code
				} else {
					wg.Done()
				}
			}(code)
		}
	}()

	for _, code := range codes {
		wg.Add(1)
		fetchCh <- code
	}

	wg.Wait()
	close(fetchCh)
	return nil
}

func init() {
	StatCmd.Flags().IntVarP(&periodVar, "period", "p", 30, "get the <period> days of stat")
	StatCmd.Flags().BoolVarP(&includeST, "includeST", "t", false, "exclude the st from results")
	StatCmd.Flags().BoolVarP(&sortByIncrease, "sortByIncrease", "r", false, "sort by the increase rate with descending sort")

	MyStatCmd.Flags().IntVarP(&periodVar, "period", "p", 30, "get the <period> days of stat")
	MyStatCmd.Flags().BoolVarP(&includeST, "includeST", "t", false, "exclude the st from results")
	MyStatCmd.Flags().StringVarP(&groupVar, "group", "g", "", "show stat of specified group")
	MyStatCmd.Flags().BoolVarP(&sortByIncrease, "sortByIncrease", "r", false, "sort by the increase rate with descending sort")

	RPSCmd.Flags().IntVarP(&daysVar, "days", "d", 5, "get <days> days of rps")
	RPSCmd.Flags().IntVarP(&showNumVar, "showNum", "n", 50, "show top <showNum> records of rps rank")

	rootCmd.AddCommand(StatCmd)
	rootCmd.AddCommand(MyStatCmd)
	rootCmd.AddCommand(FetchDataCmd)
	rootCmd.AddCommand(RPSCmd)
}
