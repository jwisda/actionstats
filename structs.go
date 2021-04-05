package actionstats

//Action: This is the format of the incoming string
//I made this public because it's convenient for someone to use this struct when writing an app that works against actionstats
type Action struct {
	Action string `json:"action"`
	Time   int64  `json:"time"`
}

//Stats: This is the required format for the output
//public because it's convenient for someone to use this struct when writing an app that works against actionstats
type Stats struct {
	Action string `json:"action"`
	Avg    int64  `json:"avg"`
}

//tallyStats: internal tally of the incoming Actions
type tallyStats struct {
	totalTime int64
	count     int
}

//snapshot: a snap shot of the current data which can be persisted on any platform, db, file, key-value
//I made this public because it's convenient for someone to use this struct when writing an app that works against actionstats
type Snapshot struct {
	Action    string `json:"action"`
	TotalTime int64  `json:"totaltime"`
	Count     int    `json:"count"`
}

//Config: used to manage configuration
//I made this public because it's convenient for someone update the config directly if they have particular needs
//and can't work with the default config. Ideally this config would be stored with the serialized data but that
//is an issue for another time and far outside the spec
type Config struct {
	MinActionLength     int    //action name must be a least be this long
	MaxActionLength     int    //action name must be less or equal to this number
	MinTime             int64  //time must be greater or equal to this
	MaxTime             int64  //time must be less than this
	MaxActions          int    //just setting an upper limit
	ActionCutSet        string //set of unwanted chars in string format
	MakeActionLowerCase bool   //make the action string lower case
}
