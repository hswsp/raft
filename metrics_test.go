// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raft

import (
	"testing"
	"time"
)

func TestElectionMetrics(t *testing.T) {
	em := &ElectionMetrics{}
	
	// Test initial state
	snapshot := em.GetElectionMetrics()
	if snapshot.ElectionsInitiated != 0 {
		t.Errorf("Expected 0 elections initiated, got %d", snapshot.ElectionsInitiated)
	}
	
	// Test incrementing metrics
	em.ElectionsInitiated.Add(1)
	em.ElectionsWon.Add(1)
	em.VotesReceived.Add(3)
	em.VotesGranted.Add(1)
	em.LastElectionDuration.Store(int64(100 * time.Millisecond))
	em.TotalCandidateTime.Store(int64(150 * time.Millisecond))
	
	snapshot = em.GetElectionMetrics()
	if snapshot.ElectionsInitiated != 1 {
		t.Errorf("Expected 1 election initiated, got %d", snapshot.ElectionsInitiated)
	}
	if snapshot.ElectionsWon != 1 {
		t.Errorf("Expected 1 election won, got %d", snapshot.ElectionsWon)
	}
	if snapshot.VotesReceived != 3 {
		t.Errorf("Expected 3 votes received, got %d", snapshot.VotesReceived)
	}
	if snapshot.VotesGranted != 1 {
		t.Errorf("Expected 1 vote granted, got %d", snapshot.VotesGranted)
	}
	if snapshot.LastElectionDuration != 100*time.Millisecond {
		t.Errorf("Expected 100ms last election duration, got %v", snapshot.LastElectionDuration)
	}
}

func TestElectionSnapshotWinRate(t *testing.T) {
	tests := []struct {
		name       string
		initiated  uint64
		won        uint64
		expectedWR float64
	}{
		{"no elections", 0, 0, 0.0},
		{"all won", 4, 4, 100.0},
		{"half won", 4, 2, 50.0},
		{"none won", 4, 0, 0.0},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := ElectionSnapshot{
				ElectionsInitiated: tt.initiated,
				ElectionsWon:       tt.won,
			}
			winRate := snapshot.WinRate()
			if winRate != tt.expectedWR {
				t.Errorf("Expected win rate %f, got %f", tt.expectedWR, winRate)
			}
		})
	}
}

func TestElectionSnapshotAverageElectionDuration(t *testing.T) {
	tests := []struct {
		name             string
		initiated        uint64
		totalCandidateTime time.Duration
		expectedAvg      time.Duration
	}{
		{"no elections", 0, 0, 0},
		{"single election", 1, 100 * time.Millisecond, 100 * time.Millisecond},
		{"multiple elections", 4, 400 * time.Millisecond, 100 * time.Millisecond},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot := ElectionSnapshot{
				ElectionsInitiated: tt.initiated,
				TotalCandidateTime: tt.totalCandidateTime,
			}
			avgDuration := snapshot.AverageElectionDuration()
			if avgDuration != tt.expectedAvg {
				t.Errorf("Expected average duration %v, got %v", tt.expectedAvg, avgDuration)
			}
		})
	}
}

func TestElectionMetricsConcurrency(t *testing.T) {
	em := &ElectionMetrics{}
	
	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				em.ElectionsInitiated.Add(1)
				em.VotesReceived.Add(1)
				_ = em.GetElectionMetrics()
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	snapshot := em.GetElectionMetrics()
	if snapshot.ElectionsInitiated != 1000 {
		t.Errorf("Expected 1000 elections initiated, got %d", snapshot.ElectionsInitiated)
	}
	if snapshot.VotesReceived != 1000 {
		t.Errorf("Expected 1000 votes received, got %d", snapshot.VotesReceived)
	}
}