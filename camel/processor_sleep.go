package camel

import "time"

type sleepProcessor struct {
	name string

	duration int64
}

func newSleepProcessor(name string, dur int64) *sleepProcessor {
	return &sleepProcessor{
		name:     name,
		duration: dur,
	}
}

func (p *sleepProcessor) getName() string {
	return p.name
}

func (p *sleepProcessor) Process(_ *Exchange) {
	time.Sleep(time.Duration(p.duration) * time.Millisecond)
}
