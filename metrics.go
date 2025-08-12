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
	"sync/atomic"
	"time"
)

// ElectionMetrics tracks metrics related to leader elections
type ElectionMetrics struct {
	// Number of elections initiated by this node
	ElectionsInitiated atomic.Uint64
	
	// Number of elections won by this node
	ElectionsWon atomic.Uint64
	
	// Number of votes received from other nodes
	VotesReceived atomic.Uint64
	
	// Number of votes granted to other nodes
	VotesGranted atomic.Uint64
	
	// Last election duration in nanoseconds
	LastElectionDuration atomic.Int64
	
	// Total time spent in candidate state (nanoseconds)
	TotalCandidateTime atomic.Int64
}

// GetElectionMetrics returns a snapshot of current election metrics
func (em *ElectionMetrics) GetElectionMetrics() ElectionSnapshot {
	return ElectionSnapshot{
		ElectionsInitiated:   em.ElectionsInitiated.Load(),
		ElectionsWon:         em.ElectionsWon.Load(),
		VotesReceived:        em.VotesReceived.Load(),
		VotesGranted:         em.VotesGranted.Load(),
		LastElectionDuration: time.Duration(em.LastElectionDuration.Load()),
		TotalCandidateTime:   time.Duration(em.TotalCandidateTime.Load()),
	}
}

// ElectionSnapshot provides a point-in-time view of election metrics
type ElectionSnapshot struct {
	ElectionsInitiated   uint64
	ElectionsWon         uint64
	VotesReceived        uint64
	VotesGranted         uint64
	LastElectionDuration time.Duration
	TotalCandidateTime   time.Duration
}

// WinRate returns the election win rate as a percentage
func (es ElectionSnapshot) WinRate() float64 {
	if es.ElectionsInitiated == 0 {
		return 0.0
	}
	return float64(es.ElectionsWon) / float64(es.ElectionsInitiated) * 100.0
}

// AverageElectionDuration returns the average time spent per election
func (es ElectionSnapshot) AverageElectionDuration() time.Duration {
	if es.ElectionsInitiated == 0 {
		return 0
	}
	return es.TotalCandidateTime / time.Duration(es.ElectionsInitiated)
}