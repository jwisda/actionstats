package actionstats_test

import (
	"actionstats"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

func TestHappyPath(t *testing.T) {
	as := actionstats.New()

	if err := as.AddAction("{\"action\":\"jump\", \"time\": 100}"); err != nil {
		t.Error("failed to addAction ")
	}

	if err := as.AddAction("{\"action\":\"run\", \"time\": 75}"); err != nil {
		t.Error("failed to addAction ")
	}

	if err := as.AddAction("{\"action\":\"jump\", \"time\": 200}"); err != nil {
		t.Error("failed to addAction ")
	}

	test := as.GetStats()
	fmt.Printf("happy results %v\n", test)
	if test != `[{"action":"jump","avg":150},{"action":"run","avg":75}]` {
		t.Error("Happy Path test stats failed")
	}
}

func TestNegTime(t *testing.T) {
	as := actionstats.New()
	err := as.AddAction("{\"action\":\"jump\", \"time\": -100}")
	if !ErrorContains(err, "ActionStats: Action jump Time -100 is invalid") {
		t.Error("Negative Time Failed test ")
	}
}

func TestBadJson(t *testing.T) {
	as := actionstats.New()
	err := as.AddAction("{\"action\":120, \"time\": \"-100\"}")
	if !ErrorContains(err, "ActionStats: Action  Time 0 is invalid") {
		t.Error("Bad Json Test failed ")
	}
}

func TestNameTooLong(t *testing.T) {
	as := actionstats.New()
	err := as.AddAction("{\"action\":\"THIS IS OVER TWENTY CHARACTERS\", \"time\": 100}")
	if !ErrorContains(err, "ActionStats: Action Key 'THIS IS OVER TWENTY CHARACTERS' is invalid") {
		t.Error("Name Too Long failed test ")
	}
}

func TestNameTooShort(t *testing.T) {
	as := actionstats.New()
	err := as.AddAction("{\"action\":\"\", \"time\": 100}")
	if !ErrorContains(err, "ActionStats: Action Key '' is invalid") {
		t.Error("Name Too Short failed test ")
	}
}

func TestTimeTooBig(t *testing.T) {
	as := actionstats.New()
	test := fmt.Sprintf("{\"action\":\"jump\", \"time\": %v }", int(math.Pow(2, 60)-1))
	err := as.AddAction(test)
	if !ErrorContains(err, "ActionStats: Action jump Time 1152921504606846976 is invalid") {
		t.Error("Time too big failed test ")
	}
}

//TestTimeOverflow this test should hit the overflow test for integers
func TestTimeOverflow(t *testing.T) {
	as := actionstats.New()
	as.Config.MaxTime = int64(math.Pow(2, 61) - 1)
	test := fmt.Sprintf("{\"action\":\"jump\", \"time\": %v }", int(math.Pow(2, 61)-1))
	for i := 0; i < 4; i++ {
		err := as.AddAction(test)
		if err != nil {
			log.Printf("TimeOverflow: %v", err)
			if !ErrorContains(err, "ActionStats: Action jump exceeds MaxInt64") {
				t.Error("TimeOverflow failed test ")
			}
		}
	}
}

func TestMaxActions(t *testing.T) {
	as := actionstats.New()
	as.Config.MaxActions = 8
	for i := 0; i < 9; i++ {
		test := fmt.Sprintf("{\"action\":\"test%v\", \"time\": %v }", i, 10)
		err := as.AddAction(test)
		if err != nil {
			log.Printf("MaxActions: %v", err)
			if !ErrorContains(err, "MaxActions") {
				t.Error("MaxActions failed test ")
			}
		}
	}
}

func TestConcurrency(t *testing.T) {
	as := actionstats.New()
	testChars := []byte("asdfghjklqwertyuiopzxcvbnm")

	var wg sync.WaitGroup //wait for it...

	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		action := ""
		for i := 0; i < 3; i++ { //26^3 = 17576 string combos
			action = action + string(testChars[rand.Intn(len(testChars))])
		}
		time := rand.Intn(10000)
		actionStr := fmt.Sprintf("{\"action\":\"%v\", \"time\":%v}", action, time)

		go func(actionstr string, wg *sync.WaitGroup) { //make concurrent
			defer wg.Done()
			err := as.AddAction(actionstr)
			if err != nil && !ErrorContains(err, "ActionStats: Action jump exceeds MaxInt64") {
				t.Error("Concurrency test failed")
			}
		}(actionStr, &wg)

		if i%1000 == 0 {
			wg.Add(1)
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				stat := as.GetStats()

				statResult := make([]actionstats.Stats, 0, 100)
				if jsonErr := json.Unmarshal([]byte(stat), &statResult); jsonErr != nil {
					//return jsonErr
					t.Error(fmt.Sprintf("Concurrency stat results failed to unmarshall %v", jsonErr))
				}
				//log.Printf("Length of Stat results %v %v\n", len(statResult), i)

			}(i, &wg)
		}
	}

	wg.Wait()             //wait for all processes to complete
	stat := as.GetStats() //one final GetStats just to make sure everything is working

	statResult := make([]actionstats.Stats, 0, 100)
	if jsonErr := json.Unmarshal([]byte(stat), &statResult); jsonErr != nil {
		//return jsonErr
		t.Error(fmt.Sprintf("Concurrency stat results failed to unmarshall %v", jsonErr))
	}

	fmt.Printf("Final length of stat results %v \n", len(statResult))
}

func TestEmptySerialization(t *testing.T) {
	as := actionstats.New()
	//snapshot := as.TakeSnapshot()
	snapshot := "[[[[{{}}]]]]"
	err := as.LoadSnapshot(snapshot)
	if err != nil && !ErrorContains(err, "ActionStats: LoadSnapShot Failed invalid character '{' ") {
		t.Error("Failed snapshot test")
	}
}

func TestEmptyStats(t *testing.T) {
	as := actionstats.New()
	testStat := as.GetStats()
	if testStat != "[]" {
		t.Error("failed EmptyStats test")
	}
}

func TestSerialization(t *testing.T) {
	testChars := []byte("asdfghjklqwertyuiopzxcvbnm!@#$%^&*(_+-={}[]:;<>?,./")
	as := actionstats.New()

	for i := 0; i < 100000; i++ {
		action := ""
		for i := 0; i < 3; i++ {
			action = action + string(testChars[rand.Intn(len(testChars))])
		}
		time := rand.Intn(10000) - 100
		data := fmt.Sprintf("{\"action\":\"%v\", \"time\":%v}", action, time)
		err := as.AddAction(data)
		if err != nil {
			//log.Printf("error trapped on TestSerializtion %v", err)
		}
	}

	stats1 := as.GetStats()
	snapshot := as.TakeSnapshot()

	//create new object
	as2 := actionstats.New()

	err := as2.LoadSnapshot(snapshot)
	if err != nil {
		t.Error("Failed snapshot test")
	}
	stats2 := as2.GetStats()

	if stats1 != stats2 {
		t.Error("Test Serialization failed to reload into new object")
	}

	//use existing object
	err = as.LoadSnapshot(snapshot)
	if err != nil {
		t.Error("Failed snapshot test")
	}
	stats2 = as.GetStats()

	if stats1 != stats2 {
		t.Error("Test Serialization failed to reload into original object")
	}
}

func ErrorContains(err error, errString string) bool {
	if err == nil {
		return errString == ""
	}
	if errString == "" {
		return false
	}
	return strings.Contains(err.Error(), errString)
}
