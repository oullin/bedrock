package jobqueue

import (
	"math"
	"sort"
	"time"
)

// BalanceStrategy controls how queue load is weighted during auto scaling.
type BalanceStrategy string

const (
	// BalanceByTime weights queues by their estimated time to clear.
	BalanceByTime BalanceStrategy = "time"
	// BalanceBySize weights queues by pending job count.
	BalanceBySize BalanceStrategy = "size"
)

// AutoScaleOptions contains process boundaries for queue balancing.
type AutoScaleOptions struct {
	MinProcesses int
	MaxProcesses int
	MaxShift     int
	Strategy     BalanceStrategy
}

// ProcessRecommendation is the target process count for one queue.
type ProcessRecommendation struct {
	Queue     string
	Processes int
}

// RecommendProcesses returns deterministic per-queue process recommendations.
func RecommendProcesses(snapshot Snapshot, current map[string]int, options AutoScaleOptions) []ProcessRecommendation {
	queues := append([]QueueStatus(nil), snapshot.Queues...)
	sort.SliceStable(queues, func(i, j int) bool {
		return queues[i].Name < queues[j].Name
	})

	if len(queues) == 0 {
		return nil
	}

	minProcesses := options.MinProcesses
	if minProcesses <= 0 {
		minProcesses = 1
	}

	desiredTotal := minProcesses * len(queues)
	if hasPendingJobs(queues) && options.MaxProcesses > desiredTotal {
		desiredTotal = options.MaxProcesses
	}
	if options.MaxProcesses > 0 && desiredTotal > options.MaxProcesses {
		desiredTotal = options.MaxProcesses
	}
	if desiredTotal < minProcesses*len(queues) {
		desiredTotal = minProcesses * len(queues)
	}

	targets := allocateProcesses(queues, desiredTotal, minProcesses, options.Strategy)
	recommendations := make([]ProcessRecommendation, 0, len(queues))

	for _, queue := range queues {
		processes := targets[queue.Name]
		if options.MaxShift > 0 && current != nil {
			processes = limitShift(processes, current[queue.Name], options.MaxShift, minProcesses)
		}

		if processes < minProcesses {
			processes = minProcesses
		}

		recommendations = append(recommendations, ProcessRecommendation{
			Queue:     queue.Name,
			Processes: processes,
		})
	}

	return recommendations
}

func allocateProcesses(queues []QueueStatus, total, minProcesses int, strategy BalanceStrategy) map[string]int {
	targets := make(map[string]int, len(queues))
	for _, queue := range queues {
		targets[queue.Name] = minProcesses
	}

	remaining := total - minProcesses*len(queues)
	if remaining <= 0 {
		return targets
	}

	weights := make([]float64, len(queues))
	totalWeight := 0.0
	for i, queue := range queues {
		weights[i] = queueWeight(queue, strategy)
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		for i := 0; i < remaining; i++ {
			targets[queues[i%len(queues)].Name]++
		}

		return targets
	}

	type remainder struct {
		index int
		value float64
	}

	remainders := make([]remainder, 0, len(queues))
	assigned := 0
	for i, weight := range weights {
		share := (weight / totalWeight) * float64(remaining)
		extra := int(math.Floor(share))
		targets[queues[i].Name] += extra
		assigned += extra
		remainders = append(remainders, remainder{index: i, value: share - float64(extra)})
	}

	sort.SliceStable(remainders, func(i, j int) bool {
		if remainders[i].value == remainders[j].value {
			return queues[remainders[i].index].Name < queues[remainders[j].index].Name
		}

		return remainders[i].value > remainders[j].value
	})

	for i := 0; assigned < remaining; i++ {
		targets[queues[remainders[i%len(remainders)].index].Name]++
		assigned++
	}

	return targets
}

func hasPendingJobs(queues []QueueStatus) bool {
	for _, queue := range queues {
		if queue.Pending > 0 || queue.Reserved > 0 {
			return true
		}
	}

	return false
}

func queueWeight(queue QueueStatus, strategy BalanceStrategy) float64 {
	if strategy == BalanceBySize {
		return float64(max(queue.Pending, 0))
	}

	if queue.Wait > 0 {
		return queue.Wait.Seconds()
	}

	if queue.Pending <= 0 {
		return 0
	}

	if queue.Throughput <= 0 {
		return float64(queue.Pending)
	}

	return (time.Duration(queue.Pending) * time.Minute / time.Duration(queue.Throughput)).Seconds()
}

func limitShift(target, current, maxShift, minProcesses int) int {
	if current <= 0 {
		current = minProcesses
	}

	upper := current + maxShift
	lower := current - maxShift
	if lower < minProcesses {
		lower = minProcesses
	}

	if target > upper {
		return upper
	}
	if target < lower {
		return lower
	}

	return target
}
