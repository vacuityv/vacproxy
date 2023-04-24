package service

import (
	"errors"
	"sync"
	"time"
)

const (
	workerBits   uint8 = 10
	maxWorker    int64 = -1 ^ (-1 << workerBits)
	sequenceBits uint8 = 12
	sequenceMask int64 = -1 ^ (-1 << sequenceBits)
	timeShift    uint8 = workerBits + sequenceBits
	workerShift  uint8 = sequenceBits
)

type Snowflake struct {
	workerId      int64
	sequence      int64
	lastTimestamp int64
	mutex         sync.Mutex
}

func NewSnowflake(workerId int64) (*Snowflake, error) {
	if workerId < 0 || workerId > maxWorker {
		return nil, errors.New("worker id out of range")
	}
	return &Snowflake{
		workerId:      workerId,
		sequence:      0,
		lastTimestamp: -1,
	}, nil
}

func (sf *Snowflake) NextId() (int64, error) {
	sf.mutex.Lock()
	defer sf.mutex.Unlock()
	timestamp := time.Now().UnixNano() / 1000000
	if timestamp < sf.lastTimestamp {
		return 0, errors.New("time is moving backwards")
	}
	if timestamp == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & sequenceMask
		if sf.sequence == 0 {
			timestamp = sf.nextTimestamp()
		}
	} else {
		sf.sequence = 0
	}
	sf.lastTimestamp = timestamp
	return (timestamp << timeShift) | (sf.workerId << workerShift) | sf.sequence, nil
}

func (sf *Snowflake) nextTimestamp() int64 {
	timestamp := time.Now().UnixNano() / 1000000
	for timestamp <= sf.lastTimestamp {
		timestamp = time.Now().UnixNano() / 1000000
	}
	return timestamp
}
