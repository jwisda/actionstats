package actionstats

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"sort"
	"strings"
	"sync"
)

type ActionStats struct {
	//convenience placeholder for the name of this collection of stats
	Name string

	//contains all the configurations of the component
	Config *Config

	//maximum stored time
	bigMaxInt *big.Int

	//ActionMux Keep updates and reads in sync
	actionMux *sync.Mutex //used to make sure updates to actionTally are managed one at a time

	//ActionTally //raw data to calc results
	actionTally map[string]tallyStats
}

//New function that returns a new object of type ActionStats
func New(name string) ActionStats {
	var config Config
	var actMux sync.Mutex
	bigMaxInt := big.NewInt(math.MaxInt64)
	actionTally := make(map[string]tallyStats)
	as := ActionStats{name, &config, bigMaxInt, &actMux, actionTally}
	as.setDefaults()

	return as
}

func (act ActionStats) AddAction(actionJson string) error {
	var err error
	var actionResult Action

	//json unmarshall will check that action is a string and time is a number < maxInt
	if jsonErr := json.Unmarshal([]byte(actionJson), &actionResult); jsonErr != nil {
		//return jsonErr
		return errors.New(fmt.Sprintf("ActionStats: Action %v Time %v is invalid", actionResult.Action, actionResult.Time))
	}

	// Add conditional checks here
	// is action valid? spaces, empty, too short, too long, cut invalid chars, make action uppercase
	actionKey := strings.Trim(actionResult.Action, act.Config.ActionCutSet)

	if act.Config.MakeActionLowerCase {
		actionKey = strings.ToLower(actionKey)
	}

	if len(actionKey) < act.Config.MinActionLength || len(actionKey) > act.Config.MaxActionLength {
		return errors.New(fmt.Sprintf("ActionStats: Action Key '%v' is invalid", actionResult.Action))
	}

	// is time negative or greater than maxTime
	if actionResult.Time < act.Config.MinTime || actionResult.Time > act.Config.MaxTime {
		return errors.New(fmt.Sprintf("ActionStats: Action %v Time %v is invalid", actionResult.Action, actionResult.Time))
	}

	act.actionMux.Lock()
	defer act.actionMux.Unlock()

	if actionItem, ok := act.actionTally[actionKey]; !ok {
		if len(act.actionTally)+1 < act.Config.MaxActions {
			act.actionTally[actionKey] = tallyStats{actionResult.Time, 1}
		} else {
			return errors.New(fmt.Sprintf("ActionStats: Action %v MaxActions %v are exceeded", actionResult.Action, act.Config.MaxActions))
		}
	} else {
		total := int64(0)
		var mathErr error
		//will tally time exceed maxInt if we add this new time?
		if mathErr, total = act.addAndCheckOverflow(actionItem.totalTime, actionResult.Time); mathErr != nil {
			return errors.New(fmt.Sprintf("ActionStats: Action %v exceeds MaxInt64", actionResult.Action))
		}
		actionItem.totalTime = total
		actionItem.count += 1 //should we check for maxInt64?
		act.actionTally[actionKey] = tallyStats{actionItem.totalTime, actionItem.count}
	}

	return err
}

func (act ActionStats) GetStats() string {
	statResults := "[]" //empty json array
	stats := make([]Stats, 0, len(act.actionTally))

	act.actionMux.Lock() //no updates while calculating results
	for act, tally := range act.actionTally {
		stats = append(stats, Stats{act, tally.totalTime / int64(tally.count)})
	}
	act.actionMux.Unlock() //go ahead and unlock before formating results

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Action < stats[j].Action
	})

	if jsonBytes, err := json.Marshal(stats); err == nil {
		statResults = string(jsonBytes)
	} else {
		log.Printf("ActionStats: getStats() returned an empty set Error: %v", err)
	}

	return statResults
}
