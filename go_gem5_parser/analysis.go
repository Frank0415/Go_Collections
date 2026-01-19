package main

import (
	"math"
)

// CacheStats holds the access statistics for a cache level.
type CacheStats struct {
	Access   uint64  // Total number of cache accesses (Hits + Misses)
	Hits     uint64  // Total number of cache hits
	Misses   uint64  // Total number of cache misses
	MissRate float64 // MissRate of Cache
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

	//

	// --- Pipeline Slot Metrics (Top-Down Base) ---
	SlotsIssued  uint64 // Total number of instructions issued to the backend (Dispatched)
	SlotsRetired uint64 // Total number of instructions committed (Retired)

	// --- Pipeline Stalls & Bubbles ---
	MispredRetired  uint64 // Count of retired branch mispredictions
	FetchCycles     uint64 // Cycles stalled due to I-Cache miss response latency
	RecoveryCycles  uint64 // Cycles stalled due to pipeline squashes (Branch Recovery)
	MachineClears   uint64 // Cycles stalled due to machine clears (e.g., memory ordering flushes)
	FetchStallSlots uint64 // Total Slots on the Pipeline that have been stalled

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
	MemLevelParallel float64    // The level of parallelism that the memory has, match ruby_system.m_outstandReqHistSeqr::mean
	L1D              CacheStats // L1 Data Cache statistics
	L1I              CacheStats // L1 Instruction Cache statistics
	L2               CacheStats // L2 Unified Cache statistics
	L3               CacheStats // L3 Cache statistics

	// --- Thread Context Stats ---
	Thread0 ThreadData // Statistics for hardware thread 0

	// Processed Stats Here:

}

type TMAOutStats struct {
	L1_retire           float64
	L1_badspec          float64
	L1_frontend         float64
	L1_backend          float64
	L0_fullfrontend     float64
	L0_frontendutil     float64
	L0_BranchPrediction float64
}

type L1TMAStats struct {
	L1_retire    float64
	L1_badspec   float64
	L1_frontend  float64
	L1_backend   float64
}

type L2TMAStats struct {
	// Frontend
	L2_fetch_latency   float64
	L2_fetch_bandwidth float64
	// Bad Speculation
	L2_branch_mispredict float64
	L2_machine_clear     float64
	// Backend
	L2_memory_bound float64
	L2_core_bound   float64
}

type TMAStats struct {
	pmu   *PMUStats
	mytma *TMAOutStats // TMA Stats collected from GEM5 Directly
	tmaL1 *L1TMAStats  // TMA Stats calculated
	tmaL2 *L2TMAStats  // Level 2 TMA stats Calculated
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
		"board.processor.cores.core.L0_BranchPrediction":  &mytma.L0_BranchPrediction,
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
		"board.processor.cores.core.fetch.status::icacheWaitResponse": &pmu.FetchCycles,
		"board.processor.cores.core.fetch.status::squashing":          &pmu.RecoveryCycles,
		"board.processor.cores.core.iew.dispatchStatus::squashing":    &pmu.MachineClears,
		"board.processor.cores.core.fetch.fetchStallSlots":            &pmu.FetchStallSlots,

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
		"board.processor.cores.core.lsq0.loadToUse::mean":               &pmu.MeanLoadAccessTime,
		"board.cache_hierarchy.ruby_system.m_outstandReqHistSeqr::mean": &pmu.MemLevelParallel,
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

	CacheList := []*CacheStats{&pmu.L1D, &pmu.L1I, &pmu.L2, &pmu.L3}

	for _, cache := range CacheList {
		if cache != nil && cache.Access != 0 {
			cache.MissRate = float64(cache.Misses) / float64(cache.Access)
		}
	}

	return pmu
}

func calcL1(pmu *PMUStats) *L1TMAStats {
	stats := new(L1TMAStats)
	var iWidth uint64 = 8
	var dWidth uint64 = 8
	var TotalSlots = pmu.Cycles * dWidth
	var SlotsIssued = pmu.SlotsIssued
	var SlotsRetired = pmu.SlotsRetired
	var RecoveryBubbles = pmu.RecoveryCycles * iWidth

	stats.L1_frontend = float64(pmu.FetchStallSlots) / float64(TotalSlots)
	stats.L1_badspec = float64(SlotsIssued-SlotsRetired+RecoveryBubbles) / float64(TotalSlots)
	stats.L1_retire = float64(SlotsRetired) / float64(TotalSlots)
	stats.L1_backend = 1 - (stats.L1_frontend + stats.L1_badspec + stats.L1_retire)

	// fmt.Printf("Backend: %.4f\nFrontend: %.4f\nRetire: %.4f\nBadSpeculation: %.4f\n", stats.L1_backend, stats.L1_frontend, stats.L1_retire, stats.L1_badspec)
	// L1TMAStats.L1_backend = PMUStats.num
	return stats
}

/*
 * We suppose L1 has 3 cycles of latency
 * L2 has 12 cycles
 * Memory has 40 cycles + 20000 ticks of latency, for 3GHZ -> 100 cycles
 */

func calcL2(pmu *PMUStats, l1 *L1TMAStats) *L2TMAStats {
	l2 := new(L2TMAStats)

	// --- 基于 JSON 配置校准的常量 ---
	const L2Lat = 8.0                      // 经计算约 7-9 cycles
	const MemLat = 90.0                    // 经计算约 80-95 cycles
	var MLP float64 = pmu.MemLevelParallel // 假设平均内存并行度为 2 (针对乱序核)

	// 1. Fetch Latency (建议检查 FetchCycles 是否仅包含 I-Cache 停顿)
	l2.L2_fetch_latency = float64(pmu.FetchCycles) / float64(pmu.Cycles)
	l2.L2_fetch_bandwidth = math.Max(0, l1.L1_frontend-l2.L2_fetch_latency)

	// 2. Bad Speculation
	totalBadEvents := float64(pmu.MispredRetired + pmu.MachineClears)
	if totalBadEvents > 0 {
		frac := float64(pmu.MispredRetired) / totalBadEvents
		l2.L2_branch_mispredict = l1.L1_badspec * frac
		l2.L2_machine_clear = l1.L1_badspec * (1.0 - frac)
	}

	// 3. Backend Bound - 更加精确的内存建模
	// 计算 L1 缺失但命中 L2 的次数
	l1MissL2Hit := float64(pmu.L1D.Misses) - float64(pmu.L2.Misses)
	if l1MissL2Hit < 0 {
		l1MissL2Hit = 0
	}

	// 估算受限周期，引入 MLP 因子防止过大
	rawMemStall := (l1MissL2Hit * L2Lat) + (float64(pmu.L2.Misses) * MemLat)
	adjMemStall := rawMemStall / MLP

	l2.L2_memory_bound = adjMemStall / float64(pmu.Cycles)

	// 归一化限制：Memory Bound 不能超过总 Backend Bound
	if l2.L2_memory_bound > l1.L1_backend {
		l2.L2_memory_bound = l1.L1_backend
	}
	l2.L2_core_bound = l1.L1_backend - l2.L2_memory_bound

	return l2
}

func GetStats(entries *map[string]Entry) *TMAStats {
	stats := new(TMAStats)
	stats.pmu = getPMU(entries)
	stats.mytma = getMyTMA(entries)
	stats.tmaL1 = calcL1(stats.pmu)
	stats.tmaL2 = calcL2(stats.pmu, stats.tmaL1)
	return stats
}
