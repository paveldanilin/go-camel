package camel

import "time"

type sleepProcessor struct {
	id string

	duration int64
}

func newSleepProcessor(id string, dur int64) *sleepProcessor {
	return &sleepProcessor{
		id:       id,
		duration: dur,
	}
}

func (p *sleepProcessor) getId() string {
	return p.id
}

func (p *sleepProcessor) Process(_ *Exchange) {
	time.Sleep(time.Duration(p.duration) * time.Millisecond)
}
