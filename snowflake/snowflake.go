package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	epoch             = int64(1230771600000)
	timestampBits     = uint(41)
	datacenteridBits  = uint(5)
	workeridBits      = uint(5)
	sequenceBits      = uint(12)
	timestampMax      = int64(-1 ^ (-1 << timestampBits))
	datacenteridMax   = int64(-1 ^ (-1 << datacenteridBits))
	workeridMax       = int64(-1 ^ (-1 << workeridBits))
	sequenceMask      = int64(-1 ^ (-1 << sequenceBits))
	workeridShift     = sequenceBits
	datacenteridShift = sequenceBits + workeridBits
	timestampShift    = sequenceBits + workeridBits + datacenteridBits
)

type Snowflake struct {
	sync.Mutex
	timestamp    int64
	workerid     int64
	datacenterid int64
	sequence     int64
}

var (
	current = &Snowflake{}
)

func WithBackground(sf *Snowflake) *Snowflake {
	current = sf
	return current
}

func Background() *Snowflake {
	return current
}

func NewSnowflake(datacenterid, workerid int64) *Snowflake {
	return &Snowflake{
		timestamp:    0,
		datacenterid: datacenterid & datacenteridMax,
		workerid:     workerid & workeridMax,
		sequence:     0,
	}
}

func (s *Snowflake) NextID() int64 {
	s.Lock()
	now := time.Now().UTC().UnixNano() / 1000000
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UTC().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		s.Unlock()
		panic(fmt.Errorf("epoch must be between 0 and %d", timestampMax-1))
		return 0
	}
	s.timestamp = now
	r := int64((t) << timestampShift | (s.datacenterid << datacenteridShift) | (s.workerid << workeridShift) | (s.sequence))
	s.Unlock()
	return r
}

func NextID() int64 {
	return Background().NextID()
}