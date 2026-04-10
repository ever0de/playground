package scheduler

import "fmt"

// ─────────────────────────────────────────────
// 도메인 타입
// ─────────────────────────────────────────────

// Job 은 스케줄 대상 작업을 나타냅니다.
type Job struct {
	ID          string
	ArrivalTime int // 도착 시각
	BurstTime   int // 총 CPU 실행 시간 (ms)

	// I/O 슬라이스: 각 원소는 (CPU burst, I/O duration) 쌍.
	// nil 이면 순수 CPU 작업.
	IOSlices []IOSlice
}

// IOSlice 는 CPU 실행 후 발생하는 I/O 구간입니다.
type IOSlice struct {
	CPUBurst int // 이 슬라이스에서 실행할 CPU 시간
	IOTime   int // 이후 I/O 대기 시간
}

// Result 는 단일 작업의 스케줄 결과입니다.
type Result struct {
	JobID          string
	ArrivalTime    int
	CompletionTime int
	TurnaroundTime int // Completion - Arrival
	ResponseTime   int // 첫 실행 시각 - Arrival
	WaitingTime    int // Turnaround - BurstTime
}

func (r Result) String() string {
	return fmt.Sprintf(
		"Job=%-4s Arrival=%3d Completion=%3d Turnaround=%3d Response=%3d Waiting=%3d",
		r.JobID, r.ArrivalTime, r.CompletionTime,
		r.TurnaroundTime, r.ResponseTime, r.WaitingTime,
	)
}

// SimResult 는 전체 시뮬레이션 결과입니다.
type SimResult struct {
	Results       []Result
	AvgTurnaround float64
	AvgResponse   float64
	AvgWaiting    float64
	Timeline      []TimelineEntry // 시각화용 타임라인
}

// TimelineEntry 는 특정 시각 구간의 CPU 점유 정보입니다.
type TimelineEntry struct {
	Start int
	End   int
	JobID string // "IDLE" 이면 CPU 유휴
}

// Scheduler 는 모든 스케줄링 알고리즘이 구현해야 하는 인터페이스입니다.
type Scheduler interface {
	Name() string
	Schedule(jobs []Job) SimResult
}

func calcStats(results []Result, jobs []Job) (float64, float64, float64) {
	n := float64(len(results))
	if n == 0 {
		return 0, 0, 0
	}
	burstMap := make(map[string]int, len(jobs))
	for _, j := range jobs {
		burstMap[j.ID] = j.BurstTime
	}
	var sumT, sumR, sumW float64
	for _, r := range results {
		sumT += float64(r.TurnaroundTime)
		sumR += float64(r.ResponseTime)
		sumW += float64(r.WaitingTime)
	}
	return sumT / n, sumR / n, sumW / n
}

func copyJobs(jobs []Job) []Job {
	out := make([]Job, len(jobs))
	copy(out, jobs)
	return out
}
