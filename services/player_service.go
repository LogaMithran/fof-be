package services

import (
	"context"
	"time"
)

func GetTopKScores(context context.Context) {
	valueChan := make(chan interface{})

	//value, err := connectors.SortedSetGet(context, "player-scores", 0, 20)
	time.Sleep(3 * time.Second)

	select {
	case <-context.Done():
		println("Operation timed out")
	case <-valueChan:
		//return value, nil
	}
}
