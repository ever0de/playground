package scheduler

import (
	"math"
	"testing"
)

const eps = 1e-9

func approx(a, b float64) bool { return math.Abs(a-b) < eps }

func findResult(results []Result, id string) (Result, bool) {
	for _, r := range results {
		if r.JobID == id {
			return r, true
		}
	}
	return Result{}, false
}

// TestFIFO_BasicEqual: 교재 10.3절 — 3개 작업 모두 길이 10, 동시 도착
// 기대 평균 반환 시간 = 20
func TestFIFO_BasicEqual(t *testing.T) {
	jobs := []Job{
		{ID: "A", ArrivalTime: 0, BurstTime: 10},
		{ID: "B", ArrivalTime: 0, BurstTime: 10},
		{ID: "C", ArrivalTime: 0, BurstTime: 10},
	}
	sr := (&FIFO{}).Schedule(jobs)

	rA, _ := findResult(sr.Results, "A")
	rB, _ := findResult(sr.Results, "B")
	rC, _ := findResult(sr.Results, "C")

	if rA.TurnaroundTime != 10 {
		t.Errorf("A turnaround: got %d, want 10", rA.TurnaroundTime)
	}
	if rB.TurnaroundTime != 20 {
		t.Errorf("B turnaround: got %d, want 20", rB.TurnaroundTime)
	}
	if rC.TurnaroundTime != 30 {
		t.Errorf("C turnaround: got %d, want 30", rC.TurnaroundTime)
	}
	if !approx(sr.AvgTurnaround, 20.0) {
		t.Errorf("AvgTurnaround: got %.2f, want 20.00", sr.AvgTurnaround)
	}
}

// TestFIFO_ConvoyEffect: 교재 그림 10.2 — convoy effect
// A=100, B=C=10, 동시 도착 → 평균 반환 110
func TestFIFO_ConvoyEffect(t *testing.T) {
	jobs := []Job{
		{ID: "A", ArrivalTime: 0, BurstTime: 100},
		{ID: "B", ArrivalTime: 0, BurstTime: 10},
		{ID: "C", ArrivalTime: 0, BurstTime: 10},
	}
	sr := (&FIFO{}).Schedule(jobs)

	if !approx(sr.AvgTurnaround, 110.0) {
		t.Errorf("AvgTurnaround: got %.2f, want 110.00", sr.AvgTurnaround)
	}
}

// TestFIFO_IdleGap: CPU가 유휴 구간을 올바르게 처리하는지 확인
func TestFIFO_IdleGap(t *testing.T) {
	jobs := []Job{
		{ID: "A", ArrivalTime: 5, BurstTime: 10},
	}
	sr := (&FIFO{}).Schedule(jobs)
	rA, _ := findResult(sr.Results, "A")
	if rA.CompletionTime != 15 {
		t.Errorf("Completion: got %d, want 15", rA.CompletionTime)
	}
	if rA.ResponseTime != 0 {
		t.Errorf("ResponseTime: got %d, want 0", rA.ResponseTime)
	}
}
