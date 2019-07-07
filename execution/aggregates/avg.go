package aggregates

import (
	"github.com/cube2222/octosql/execution"
	"github.com/pkg/errors"
)

type Average struct {
	averages *execution.HashMap
	counts   *execution.HashMap
}

func NewAverage() *Average {
	return &Average{
		averages: execution.NewHashMap(),
		counts:   execution.NewHashMap(),
	}
}

func (agg *Average) AddRecord(key []interface{}, value interface{}) error {
	var floatValue float64
	switch value := value.(type) {
	case float64:
		floatValue = value
	case int:
		floatValue = float64(value)
	default:
		return errors.Errorf("invalid value to average: %v with type %v", value, execution.GetType(value))
	}

	count, ok, err := agg.counts.Get(key)
	if err != nil {
		return errors.Wrap(err, "couldn't get current element count out of hashmap")
	}

	average, ok, err := agg.averages.Get(key)
	if err != nil {
		return errors.Wrap(err, "couldn't get current average out of hashmap")
	}

	var newAverage float64
	var newCount int
	if ok {
		newCount = count.(int) + 1
		newAverage = (average.(float64)*float64(newCount-1) + floatValue) / float64(newCount)
	} else {
		newCount = 1
		newAverage = floatValue
	}

	err = agg.counts.Set(key, newCount)
	if err != nil {
		return errors.Wrap(err, "couldn't put new element count into hashmap")
	}

	err = agg.averages.Set(key, newAverage)
	if err != nil {
		return errors.Wrap(err, "couldn't put new average into hashmap")
	}

	return nil
}

func (agg *Average) GetAggregated(key []interface{}) (interface{}, error) {
	average, ok, err := agg.averages.Get(key)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get average out of hashmap")
	}

	if !ok {
		return nil, errors.Errorf("average for key not found")
	}

	return average, nil
}

func (agg *Average) String() string {
	return "avg"
}
