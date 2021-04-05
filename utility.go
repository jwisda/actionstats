package actionstats

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
)

func (act ActionStats) setDefaults() {
	act.Config.MinActionLength = 1        //action name must be a least be this long
	act.Config.MaxActionLength = 20       //action name must be less or equal to this number
	act.Config.MinTime = 0                //time must be greater or equal to this
	act.Config.MaxTime = 24 * 3600 * 1000 //lets just say that time is milli-seconds and that it must be less than one day
	act.Config.MaxActions = 1000000       //arbitrary but a big number
	act.Config.ActionCutSet = " {}<>\"'`" //seems like a good set of unwanted chars
	act.Config.MakeActionLowerCase = true //don't allow mixed case actions
}

//ClearDataSet: clear out the current collection
func (act ActionStats) ClearDataSet() {
	//we'll let the GC collect and repurpose the old allocations
	act.actionMux.Lock()
	act.actionTally = make(map[string]tallyStats)
	act.actionMux.Unlock()
}

//TakeSnapshot: serialize current state of data
//this could be save a variety of ways db, file, to persist data across systems
//outside the direct spec but allows convenient testing and makes me feel better
//about writing a component that actually seems useful
func (act ActionStats) TakeSnapshot() string {
	ssResults := "[]" //snapshotresults, empty json array
	ss := make([]Snapshot, 0, len(act.actionTally))

	act.actionMux.Lock() //no updates while calculating results
	for act, tally := range act.actionTally {
		ss = append(ss, Snapshot{act, tally.totalTime, tally.count})
	}
	act.actionMux.Unlock() //go ahead and unlock

	if jsonBytes, err := json.Marshal(ss); err == nil {
		ssResults = string(jsonBytes)
	} else {
		log.Printf("ActionStats: takeSnapshot() returned an empty set Error: %v", err)
	}

	return ssResults
}

//load snapshot from previously serialized snapshot
//outside the direct spec but helps improve testing significantly
func (act ActionStats) LoadSnapshot(ss string) error {
	var err error
	snapshotResult := make([]Snapshot, 0, 100)

	if jsonErr := json.Unmarshal([]byte(ss), &snapshotResult); jsonErr != nil {
		//return jsonErr
		return errors.New(fmt.Sprintf("ActionStats: LoadSnapShot Failed %v", jsonErr))
	}

	act.ClearDataSet()
	act.actionMux.Lock() //no updates while updating data
	defer act.actionMux.Unlock()

	for _, snap := range snapshotResult {
		act.actionTally[snap.Action] = tallyStats{snap.TotalTime, snap.Count}
	}

	return err
}

//addAndCheckMaxInt: helper function to check overflow
func (act ActionStats) addAndCheckOverflow(a, b int64) (error, int64) {
	var err error
	bigA := big.NewInt(int64(a))
	bigB := big.NewInt(int64(b))
	if bigA.Add(bigA, bigB).Cmp(act.bigMaxInt) > 0 {
		return errors.New(fmt.Sprint("ActionStats: Add exceeds MaxInt64")), 0
	}
	return err, a + b
}
