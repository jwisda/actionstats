### actionstats
ActionStats is simple statistical calculation package written in Go. It takes an "Action" which includes a name and a time as an input. Upon request it can return the current average time for each action. ActionStats is written to accept concurrent calls to all functions.

### Install Package
    go get github.com/jwisda/actionstats

### Example Implementation:

    package main

    import (
        "fmt"
        "log"

        "github.com/jwisda/actionstats"
    )

    func main() {
        aStat := actionstats.New()
        actionJsonTemplate := "{\"action\":\"%v\", \"time\":%v}"

        aStat.Config.MinTime = 10
        aStat.AddAction(fmt.Sprintf(actionJsonTemplate, "jump", 100))    
        aStat.AddAction(fmt.Sprintf(actionJsonTemplate, "run", 75))    
        aStat.AddAction(fmt.Sprintf(actionJsonTemplate, "jump", 200))    

        log.Printf(aStat.GetStats())
    }

### Overview

#### func New
    returns a new object ActionStats
        
    aStat := actionstats.New()

#### func AddAction
    Adds an action to the current list of tracked actions and returns an error. 
    Errors must be checked on return to ensure an update. 

    sStat := actionstats.New()
    actionJsonTemplate := "{\"action\":\"%v\", \"time\":%v}"
    if err := aStat.AddAction(fmt.Sprintf(actionJsonTemplate, "jump", 100)); err == nil {
    	log.Print("success!!")
    }

#### func GetStats
    Returns the current Stats of the tracked actions in a json serialized string format.

    aStat := actionstats.New()
    stats := aStat.GetStats() 

#### func TakeSnapshot
    Takes a snapshot of the current tracked actions so that it can be persisted later

    aStat := actionstats.New()
    snapshot := aStat.TakeSnapshot() //returns string

#### func LoadSnapShot
    Loads a previously taken snapshot into a new or existing ActionStats object

    aStat := actionstats.New()
    // maybe add some actions here
    snapshot := aStat.TakeSnapshot() //returns string
    
    aStat2 := actionstats.New()
    aStat2.LoadSnapshot(snapshot)

#### Struct ActionStats
    Name string    
    Config type *Config

#### Struct Action
    Action string //name of action
    Time int //

#### Struct Stats
    Action string //name of action
    Avg int //average time

### Config
This package has been tested with the following default settings. 
These config settings can be changed from the defaults but not every possibility has been tested.

#### MinActionLength 
    type int default = 1        
    Action string must be a least be this long
	
#### MaxActionLength 
    type int default = 20       
    Action string length must be less or equal to this number
	
#### MinTime 
    type int default = 0                
    Time must be greater or equal to this, negative time doesn't make sense
	
#### MaxTime 
    type int default = 24 * 3600 * 1000 
    lets just say that time is milli-seconds and that it must be less than one day

#### MaxActions 
    type int default = 1000000
    maximum number of actions stored

#### ActionCutSet 
    type string default = " {}<>\"'`" 
    characters removed from the Action string name
	
#### MakeActionLowerCase 
    type boolean default = true 
    if false will allow mixed case action types


## Original Requirements

### Requirements 
This assignment may be completed in Java, Go, NodeJS, C++, or Python. Be sure to provide clear instructions on how to build and test your code. Please don’t make any assumptions about the environment that your code will be compiled/run in without explicitly stating those assumptions. Please try to limit the setup complexity by avoiding frameworks or libraries that are far from standard or require any advanced setup-- the simpler the better. If you have questions, please reach out. To submit your work please push to a public GitHub project and be sure to document any configuration or run instructions. We’re looking for a solution to the problem as well as attention to detail and code craftsmanship. Good luck and have fun! The assignment is to write a small library class that can perform the following operations:

    Add Action 
    
    addAction (string) returning error This function accepts a json serialized string of the form below and maintains an average time for each action. 3 sample inputs:
        {"action":"jump", "time":100}
        {"action":"run", "time":75}
        {"action":"jump", "time":200} 
    Assume that an end user will be making concurrent calls into this function.

    Statistics 
    
    getStats () returning string Write a second function that accepts no input and returns a serialized json array of the average time for each action that has been provided to the addAction function. Output after the 3 sample calls above would be: 
    
    [ 
        {"action":"jump", "avg":150}, 
        {"action":"run", "avg":75} 
    ] 
    
    Assume that an end user will be making concurrent calls into all functions.

##

