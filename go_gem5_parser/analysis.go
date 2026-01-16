package main

// CacheStats holds the access statistics for a cache level.
type CacheStats struct {
	Access uint64 // Total number of cache accesses (Hits + Misses)
	Hits   uint64 // Total number of cache hits
	Misses uint64 // Total number of cache misses
}

// ThreadData holds statistics specific to a hardware thread context.
type ThreadData struct {
	numInsts uint64 // Total number of instructions committed by this thread
	numOps   uint64 // Total number of micro-ops (uOps) committed by this thread
}

// PMUStats aggregates hardware performance counters from the simulation.
type PMUStats struct {
	// --- Base Timing ---
	Cycles   uint64 // Total CPU execution cycles (Clocks)
	Simticks uint64 // Total simulation time in ticks (1 tick = 1ps usually)

	// --- Pipeline Slot Metrics (Top-Down Base) ---
	InstFetched  uint64 // (Not directly mapped, instsIssued is used instead)
	SlotsIssued  uint64 // Total number of instructions issued to the backend (Dispatched)
	SlotsRetired uint64 // Total number of instructions committed (Retired)

	// --- Pipeline Stalls & Bubbles ---
	MispredRetired  uint64 // Count of retired branch mispredictions
	FetchBubbles    uint64 // Cycles stalled due to I-Cache miss response latency
	RecoveryBubbles uint64 // Cycles stalled due to pipeline squashes (Branch Recovery)
	MachineClears   uint64 // Cycles stalled due to machine clears (e.g., memory ordering flushes)

	// --- Execution Unit Metrics ---
	OpsExecuted uint64 // Total number of micro-ops executed in execution units

	// --- Structural Stalls (Event Counts) ---
	LoadQueueFull  uint64 // Count of events where the Load Queue (LQ) was full
	StoreQueueFull uint64 // Count of events where the Store Queue (SQ) was full
	InstQueueFull  uint64 // Count of events where the Instruction Queue (IQ) was full

	// --- Load/Store Unit Specifics ---
	LSQBlockedByCache  uint64  // Count of times LSQ was blocked by cache ports or contention
	MeanLoadAccessTime float64 // Average latency in cycles from Load issue to data return

	// --- Memory Hierarchy (DRAM) ---
	MemReadReqs        uint64 // Total number of read requests sent to the memory controller
	MemQueueStallCount uint64 // Count of stalls in the Ruby mandatory queue (Protocol/Contention stalls)

	// --- Cache Hierarchy Stats ---
	L1D CacheStats // L1 Data Cache statistics
	L1I CacheStats // L1 Instruction Cache statistics
	L2  CacheStats // L2 Unified Cache statistics

	// --- Thread Context Stats ---
	Thread0 ThreadData // Statistics for hardware thread 0

	// Processed Stats Here:

}

type TMAOutStats struct {
	L1_retire       float64
	L1_badspec      float64
	L1_frontend     float64
	L1_backend      float64
	L0_fullfrontend float64
	L0_frontendutil float64
}

type L1TMAStats struct {
	L1_retire   float64
	L1_badspec  float64
	L1_frontend float64
	L1_backend  float64
}

type TMAStats struct {
	pmu   *PMUStats
	mytma *TMAOutStats // TMA Stats collected from GEM5 Directly
	tmaL1 *L1TMAStats  // TMA Stats calculated
}

func getMyTMA(entries *map[string]Entry) *TMAOutStats {
	mytma := new(TMAOutStats)
	mapping := map[string]*float64{
		"board.processor.cores.core.L1_Retiring":          &mytma.L1_retire,
		"board.processor.cores.core.L1_BadSpeculation":    &mytma.L1_badspec,
		"board.processor.cores.core.L1_FrontendBound":     &mytma.L1_frontend,
		"board.processor.cores.core.L1_BackendBound":      &mytma.L1_backend,
		"board.processor.cores.core.L0_FullFrontendBound": &mytma.L0_fullfrontend,
		"board.processor.cores.core.L0_FrontendUtil":      &mytma.L0_frontendutil,
	}

	for key, targetPtr := range mapping {
		if ent, ok := (*entries)[key]; ok {
			*targetPtr = ent.Value
		}
	}

	return mytma
}

func getPMU(entries *map[string]Entry) *PMUStats {
	pmu := new(PMUStats)

	intmapping := map[string]*uint64{
		// Base
		"board.processor.cores.core.numCycles": &pmu.Cycles,
		"simTicks":                             &pmu.Simticks,

		// Pipeline Slots
		"board.processor.cores.core.instsIssued":                       &pmu.SlotsIssued,
		"board.processor.cores.core.commit.committedInstType_0::total": &pmu.SlotsRetired,

		// Pipeline Stalls
		"board.processor.cores.core.commit.branchMispredicts":         &pmu.MispredRetired,
		"board.processor.cores.core.fetch.status::icacheWaitResponse": &pmu.FetchBubbles,
		"board.processor.cores.core.fetch.status::squashing":          &pmu.RecoveryBubbles,
		"board.processor.cores.core.iew.dispatchStatus::squashing":    &pmu.MachineClears,

		// Execution
		"board.processor.cores.core.executeStats0.numInsts": &pmu.OpsExecuted,

		// Structural Stalls
		"board.processor.cores.core.rename.LQFullEvents": &pmu.LoadQueueFull,
		"board.processor.cores.core.rename.SQFullEvents": &pmu.StoreQueueFull,
		"board.processor.cores.core.rename.IQFullEvents": &pmu.InstQueueFull,

		// LSQ
		"board.processor.cores.core.lsq0.blockedByCache": &pmu.LSQBlockedByCache,

		// Memory Controller
		"board.memory.mem_ctrl.readReqs":                                                &pmu.MemReadReqs,
		"board.cache_hierarchy.ruby_system.l1_controllers.mandatoryQueue.m_stall_count": &pmu.MemQueueStallCount,

		// L1 Cache
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Dcache.m_demand_accesses": &pmu.L1D.Access,
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Dcache.m_demand_hits":     &pmu.L1D.Hits,
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Dcache.m_demand_misses":   &pmu.L1D.Misses,
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Icache.m_demand_accesses": &pmu.L1I.Access,
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Icache.m_demand_hits":     &pmu.L1I.Hits,
		"board.cache_hierarchy.ruby_system.l1_controllers.L1Icache.m_demand_misses":   &pmu.L1I.Misses,

		// L2 Cache
		"board.cache_hierarchy.ruby_system.l2_controllers.L2cache.m_demand_accesses": &pmu.L2.Access,
		"board.cache_hierarchy.ruby_system.l2_controllers.L2cache.m_demand_hits":     &pmu.L2.Hits,
		"board.cache_hierarchy.ruby_system.l2_controllers.L2cache.m_demand_misses":   &pmu.L2.Misses,

		// Thread 0
		"board.processor.cores.core.thread_0.numInsts": &pmu.Thread0.numInsts,
		"board.processor.cores.core.thread_0.numOps":   &pmu.Thread0.numOps,
	}

	floatmapping := map[string]*float64{
		"board.processor.cores.core.lsq0.loadToUse::mean": &pmu.MeanLoadAccessTime,
	}

	// Assign uint64 fields
	for key, targetPtr := range intmapping {
		if ent, ok := (*entries)[key]; ok {
			*targetPtr = uint64(ent.Value)
		}
	}

	for key, targetPtr := range floatmapping {
		if ent, ok := (*entries)[key]; ok {
			*targetPtr = ent.Value
		}
	}

	return pmu
}

func calcL1(*PMUStats) *L1TMAStats {
	stats := new(L1TMAStats)

	return stats
}

func GetStats(entries *map[string]Entry) *TMAStats {
	stats := new(TMAStats)
	stats.pmu = getPMU(entries)
	stats.mytma = getMyTMA(entries)
	stats.tmaL1 = calcL1(stats.pmu)

	return stats
}
