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

package raft_test

import (
	"fmt"
	"log"

	"go.etcd.io/raft/v3"
)

// This example demonstrates how to use the new election metrics feature
// to monitor leader election performance and behavior.
func ExampleElectionMetrics() {
	// Create a test storage
	storage := raft.NewMemoryStorage()
	
	// Create a configuration for a single node cluster
	c := &raft.Config{
		ID:              1,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	
	// Create a RawNode
	rn, err := raft.NewRawNode(c)
	if err != nil {
		log.Fatal(err)
	}
	
	// Bootstrap the cluster with this node as the only member
	err = rn.Bootstrap([]raft.Peer{{ID: 1}})
	if err != nil {
		log.Fatal(err)
	}
	
	// Get initial metrics (should be all zeros)
	metrics := rn.ElectionMetrics()
	fmt.Printf("Initial elections initiated: %d\n", metrics.ElectionsInitiated)
	fmt.Printf("Initial elections won: %d\n", metrics.ElectionsWon)
	fmt.Printf("Initial win rate: %.1f%%\n", metrics.WinRate())
	
	// Output:
	// Initial elections initiated: 0
	// Initial elections won: 0
	// Initial win rate: 0.0%
}

// This example shows how to use election metrics with a regular Node
func ExampleNode_ElectionMetrics() {
	// Create storage and config
	storage := raft.NewMemoryStorage()
	c := &raft.Config{
		ID:              1,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	
	// Start a Node
	node := raft.StartNode(c, []raft.Peer{{ID: 1}})
	defer node.Stop()
	
	// Get election metrics
	metrics := node.ElectionMetrics()
	fmt.Printf("Node metrics: Elections initiated: %d, votes granted: %d\n", 
		metrics.ElectionsInitiated, metrics.VotesGranted)
	
	// Output:
	// Node metrics: Elections initiated: 0, votes granted: 0
}