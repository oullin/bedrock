package horizon

import "time"

// DashboardStats summarizes the current queue state for a Horizon dashboard.
type DashboardStats struct {
	Queues      int
	Pending     int
	Reserved    int
	Processing  int
	Failed      int
	Throughput  int
	LongestWait time.Duration
	Paused      bool
}

// StatsForSnapshot aggregates queue counters from a snapshot.
func StatsForSnapshot(snapshot Snapshot) DashboardStats {
	stats := DashboardStats{
		Queues: len(snapshot.Queues),
		Paused: len(snapshot.Queues) > 0,
	}

	for _, queue := range snapshot.Queues {
		stats.Pending += queue.Pending
		stats.Reserved += queue.Reserved
		stats.Processing += queue.Processing
		stats.Failed += queue.Failed
		stats.Throughput += queue.Throughput

		if queue.Wait > stats.LongestWait {
			stats.LongestWait = queue.Wait
		}

		if !queue.Paused {
			stats.Paused = false
		}
	}

	return stats
}
