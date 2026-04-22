package jobqueue

import "testing"

func TestAutoScalerInventoryParity(t *testing.T) {
	t.Run("AutoScalerTest::test_scaler_attempts_to_get_closer_to_proper_balance_on_each_iteration", func(t *testing.T) {
		snapshot := Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100, Throughput: 10},
			{Name: "mail", Pending: 1, Throughput: 10},
		}}
		options := AutoScaleOptions{MinProcesses: 1, MaxProcesses: 6, MaxShift: 1}

		first := recommendationsByQueue(RecommendProcesses(snapshot, map[string]int{"default": 1, "mail": 1}, options))
		second := recommendationsByQueue(RecommendProcesses(snapshot, first, options))

		if first["default"] != 2 || second["default"] != 3 {
			t.Fatalf("default recommendations = first %d, second %d; want 2 then 3", first["default"], second["default"])
		}
	})

	t.Run("AutoScalerTest::test_balance_stays_even_when_queue_is_empty", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default"},
			{Name: "mail"},
			{Name: "reports"},
		}}, nil, AutoScaleOptions{MinProcesses: 2, MaxProcesses: 9}))

		if recommendations["default"] != 2 || recommendations["mail"] != 2 || recommendations["reports"] != 2 {
			t.Fatalf("recommendations = %#v, want even minimums", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_balancer_assigns_more_processes_on_busy_queue", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100, Throughput: 10},
			{Name: "mail", Pending: 5, Throughput: 10},
		}}, nil, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 6}))

		if recommendations["default"] <= recommendations["mail"] {
			t.Fatalf("recommendations = %#v, want default to receive more processes", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_balancing_a_single_queue_assigns_it_the_min_workers_with_empty_queue", func(t *testing.T) {
		recommendations := RecommendProcesses(Snapshot{Queues: []QueueStatus{{Name: "default"}}}, nil, AutoScaleOptions{MinProcesses: 3, MaxProcesses: 10})

		if len(recommendations) != 1 || recommendations[0].Processes != 3 {
			t.Fatalf("recommendations = %#v, want one queue with min workers", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_will_not_scale_past_max_process_threshold_under_high_load", func(t *testing.T) {
		recommendations := RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100},
			{Name: "mail", Pending: 100},
		}}, nil, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 5, Strategy: BalanceBySize})

		if totalProcesses(recommendations) > 5 {
			t.Fatalf("recommendations = %#v, want total at or below max", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_will_not_scale_below_minimum_worker_threshold", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100},
			{Name: "mail", Pending: 100},
		}}, nil, AutoScaleOptions{MinProcesses: 2, MaxProcesses: 1, Strategy: BalanceBySize}))

		if recommendations["default"] < 2 || recommendations["mail"] < 2 {
			t.Fatalf("recommendations = %#v, want every queue at minimum", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_considers_max_shift_and_attempts_to_get_closer_to_proper_balance_on_each_iteration", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100},
			{Name: "mail", Pending: 1},
		}}, map[string]int{"default": 1, "mail": 5}, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 6, MaxShift: 2, Strategy: BalanceBySize}))

		if recommendations["default"] != 3 || recommendations["mail"] != 3 {
			t.Fatalf("recommendations = %#v, want max shift to limit both queues by 2", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_does_not_permit_going_to_zero_processes_despite_exceeding_max_processes", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 100},
			{Name: "mail", Pending: 100},
		}}, nil, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 1, Strategy: BalanceBySize}))

		if recommendations["default"] == 0 || recommendations["mail"] == 0 {
			t.Fatalf("recommendations = %#v, want no queue at zero", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_assigns_more_processes_to_queue_with_more_jobs_when_using_size_strategy", func(t *testing.T) {
		recommendations := recommendationsByQueue(RecommendProcesses(Snapshot{Queues: []QueueStatus{
			{Name: "default", Pending: 50, Throughput: 100},
			{Name: "mail", Pending: 10, Throughput: 1},
		}}, nil, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 6, Strategy: BalanceBySize}))

		if recommendations["default"] <= recommendations["mail"] {
			t.Fatalf("recommendations = %#v, want size strategy to favor the larger queue", recommendations)
		}
	})

	t.Run("AutoScalerTest::test_scaler_works_with_a_single_process_pool", func(t *testing.T) {
		recommendations := RecommendProcesses(Snapshot{Queues: []QueueStatus{{Name: "default", Pending: 25}}}, nil, AutoScaleOptions{MinProcesses: 1, MaxProcesses: 5, Strategy: BalanceBySize})

		if len(recommendations) != 1 || recommendations[0].Queue != "default" || recommendations[0].Processes != 5 {
			t.Fatalf("recommendations = %#v, want one scaled process pool", recommendations)
		}
	})
}

func recommendationsByQueue(recommendations []ProcessRecommendation) map[string]int {
	byQueue := make(map[string]int, len(recommendations))
	for _, recommendation := range recommendations {
		byQueue[recommendation.Queue] = recommendation.Processes
	}

	return byQueue
}

func totalProcesses(recommendations []ProcessRecommendation) int {
	total := 0
	for _, recommendation := range recommendations {
		total += recommendation.Processes
	}

	return total
}
