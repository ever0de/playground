package scheduler

import "sort"

// ─────────────────────────────────────────────
// FIFO (First In First Out / FCFS)
// ─────────────────────────────────────────────
// 도착 순서대로 실행하며, 한 번 시작하면 완료될 때까지 실행합니다 (비선점).

type FIFO struct{}

func (f *FIFO) Name() string { return "FIFO" }

func (f *FIFO) Schedule(jobs []Job) SimResult {
	js := copyJobs(jobs)
	// 도착 시각 기준 정렬, 동일하면 ID 순
	sort.Slice(js, func(i, k int) bool {
		if js[i].ArrivalTime == js[k].ArrivalTime {
			return js[i].ID < js[k].ID
		}
		return js[i].ArrivalTime < js[k].ArrivalTime
	})

	var (
		results  []Result
		timeline []TimelineEntry
		now      int
	)

	for _, j := range js {
		// CPU가 작업 도착 전에 비어 있으면 유휴 구간 기록
		if now < j.ArrivalTime {
			timeline = append(timeline, TimelineEntry{now, j.ArrivalTime, "IDLE"})
			now = j.ArrivalTime
		}

		firstRun := now
		end := now + j.BurstTime
		timeline = append(timeline, TimelineEntry{now, end, j.ID})
		now = end

		results = append(results, Result{
			JobID:          j.ID,
			ArrivalTime:    j.ArrivalTime,
			CompletionTime: now,
			TurnaroundTime: now - j.ArrivalTime,
			ResponseTime:   firstRun - j.ArrivalTime,
			WaitingTime:    firstRun - j.ArrivalTime,
		})
	}

	avgT, avgR, avgW := calcStats(results, js)
	return SimResult{results, avgT, avgR, avgW, timeline}
}
