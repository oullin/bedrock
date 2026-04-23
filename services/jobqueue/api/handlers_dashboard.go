package api

import (
	"net/http"

	"github.com/bedrock/packages/jobqueue"
)

// dashboardStats ports DashboardStatsController::index. Upstream JobQueue
// returns the current snapshot counters, the number of master supervisors,
// and whether every one of them is paused.
func dashboardStats(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snapshot, _ := opts.Repository.Latest(r.Context())
		stats := jobqueue.StatsForSnapshot(snapshot)
		supervisors := opts.Supervisors.All()

		paused := len(supervisors) > 0

		for _, supervisor := range supervisors {
			if supervisor.Status != "paused" {
				paused = false

				break
			}
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"jobsPerMinute":          stats.Throughput,
			"processes":              len(supervisors),
			"queueWithMaxRuntime":    longestWaitQueue(snapshot),
			"queueWithMaxThroughput": highestThroughputQueue(snapshot),
			"failedJobs":             stats.Failed,
			"pendingJobs":            stats.Pending,
			"status":                 supervisorStatus(paused, supervisors),
			"wait":                   map[string]int64{"longestWait": int64(stats.LongestWait.Seconds())},
		})
	}
}

func supervisorStatus(paused bool, supervisors []MasterSupervisor) string {
	if len(supervisors) == 0 {
		return "inactive"
	}

	if paused {
		return "paused"
	}

	return "running"
}

func longestWaitQueue(snapshot jobqueue.Snapshot) string {
	var (
		name    string
		longest = int64(-1)
	)

	for _, queue := range snapshot.Queues {
		seconds := int64(queue.Wait.Seconds())

		if seconds > longest {
			longest = seconds
			name = queue.Name
		}
	}

	return name
}

func highestThroughputQueue(snapshot jobqueue.Snapshot) string {
	var (
		name string
		best = -1
	)

	for _, queue := range snapshot.Queues {
		if queue.Throughput > best {
			best = queue.Throughput
			name = queue.Name
		}
	}

	return name
}

// masterSupervisors ports MasterSupervisorController::index. The upstream
// controller returns the collection of master supervisors reported by
// MasterSupervisor::all().
func masterSupervisors(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		supervisors := opts.Supervisors.All()

		if supervisors == nil {
			supervisors = []MasterSupervisor{}
		}

		writeJSON(w, http.StatusOK, supervisors)
	}
}
