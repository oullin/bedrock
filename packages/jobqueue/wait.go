package jobqueue

import (
	"math"
	"sort"
	"time"
)

// QueueWaitTime is the estimated time required to clear one queue.
type QueueWaitTime struct {
	Queue    string
	Duration time.Duration
}

// WaitTimes estimates queue clear times from snapshot throughput.
func WaitTimes(snapshot Snapshot, queues ...string) []QueueWaitTime {
	filter := make(map[string]struct{}, len(queues))

	for _, queue := range queues {
		filter[queue] = struct{}{}
	}

	times := make([]QueueWaitTime, 0, len(snapshot.Queues))

	for _, queue := range snapshot.Queues {
		if len(filter) > 0 {
			if _, ok := filter[queue.Name]; !ok {
				continue
			}
		}

		times = append(times, QueueWaitTime{
			Queue:    queue.Name,
			Duration: queueClearDuration(queue),
		})
	}

	sort.SliceStable(times, func(i, j int) bool {
		return times[i].Queue < times[j].Queue
	})

	return times
}

func queueClearDuration(queue QueueStatus) time.Duration {
	if queue.Pending <= 0 {
		return 0
	}

	if queue.Wait > 0 {
		return queue.Wait
	}

	if queue.Throughput <= 0 {
		return 0
	}

	minutes := math.Ceil(float64(queue.Pending) / float64(queue.Throughput))

	return time.Duration(minutes) * time.Minute
}
