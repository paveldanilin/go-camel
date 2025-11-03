package exchange

import (
	"sync"
	"time"
)

// The MessageHistory maintains a list of all components that the message passed through.
// Every component that processes the message (including the originator) adds one entry to the list.
// The MessageHistory should be part of the message header because it contains system-specific control information.
// Keeping this information in the header separates it from the message body that contains application specific data.
type MessageHistory struct {
	mu      sync.RWMutex
	records []*MessageHistoryRecord
}

func NewMessageHistory() *MessageHistory {
	return &MessageHistory{records: []*MessageHistoryRecord{}}
}

func (mh *MessageHistory) AddRecord(r *MessageHistoryRecord) {
	mh.mu.Lock()
	defer mh.mu.Unlock()
	mh.records = append(mh.records, r)
}

func (mh *MessageHistory) Records() []*MessageHistoryRecord {
	mh.mu.RLock()
	defer mh.mu.RUnlock()

	// make shallow copy
	records := make([]*MessageHistoryRecord, len(mh.records))
	for i, rec := range mh.records {
		records[i] = rec
	}

	return records
}

type MessageHistoryRecord struct {
	time        time.Time
	elapsedTime int64
	routeName   string
	stepName    string
}

func NewMessageHistoryRecord(routeName, stepName string) *MessageHistoryRecord {
	return &MessageHistoryRecord{
		routeName:   routeName,
		stepName:    stepName,
		time:        time.Now(),
		elapsedTime: -1,
	}
}

func (mhr *MessageHistoryRecord) ElapsedTime() int64 {
	return mhr.elapsedTime
}

func (mhr *MessageHistoryRecord) UpdateElapsedTime() {
	if mhr.elapsedTime < 0 {
		mhr.elapsedTime = time.Since(mhr.time).Milliseconds()
	}
}

func (mhr *MessageHistoryRecord) Time() time.Time {
	return mhr.time
}

func (mhr *MessageHistoryRecord) RouteName() string {
	return mhr.routeName
}

func (mhr *MessageHistoryRecord) StepName() string {
	return mhr.stepName
}

func (mhr *MessageHistoryRecord) Message() *Message {
	return nil
}
