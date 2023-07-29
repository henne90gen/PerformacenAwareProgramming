#pragma once

#include "shared.h"

#include <array>
#include <string>
#include <unordered_map>

#if WIN32
#include <intrin.h>
#else
#include <x86intrin.h>
#endif

u64 EstimateCPUTimerFrequency();
u64 CPUTimerDiffToNanoseconds(u64 cpuTimer, u64 cpuTimerFreq);

inline u64
ReadCPUTimer(void) {
    return __rdtsc();
}

struct TimeAggregate {
    std::string label = {};
    u64 elapsedWithoutChildren = 0;
    u64 elapsedWithChildren = 0;
    u64 hitCount = 0;
};

struct Profiler {
    u64 start = 0;
    u32 nextAggregateIndex = 1;
    std::unordered_map<std::string, u32> timeAggregateIndices = {};
    std::array<TimeAggregate, 4096> timeAggregates = {};

    Profiler();
};

extern Profiler GlobalProfiler;
extern u32 GlobalProfilerParentIndex;

struct Timer {
    std::string label = {};
    u64 parentIndex = 0;
    u64 start = 0;
    u64 oldElapsedWithChildren = 0;

    inline Timer(const std::string &name) : label(name) {
        u32 index = 0;
        auto itr = GlobalProfiler.timeAggregateIndices.find(name);
        if (itr == GlobalProfiler.timeAggregateIndices.end()) {
            index = GlobalProfiler.nextAggregateIndex++;
            GlobalProfiler.timeAggregateIndices[name] = index;
        } else {
            index = itr->second;
        }

        auto &aggregate = GlobalProfiler.timeAggregates[index];
        oldElapsedWithChildren = aggregate.elapsedWithChildren;
        parentIndex = GlobalProfilerParentIndex;
        GlobalProfilerParentIndex = index;
        start = ReadCPUTimer();
    }

    inline ~Timer() {
        GlobalProfilerParentIndex = parentIndex;

        auto end = ReadCPUTimer();
        auto elapsed = end - start;

        auto index = GlobalProfiler.timeAggregateIndices[label];
        auto &parentAggregate = GlobalProfiler.timeAggregates[parentIndex];
        auto &aggregate = GlobalProfiler.timeAggregates[index];

        parentAggregate.elapsedWithoutChildren -= elapsed;
        aggregate.elapsedWithoutChildren += elapsed;
        aggregate.elapsedWithChildren = oldElapsedWithChildren + elapsed;
        aggregate.hitCount++;
        aggregate.label = label;
    }
};

void BeginProfiling();
void EndProfiling();

#define CAT_(a, b) a##b
#define CAT(a, b) CAT_(a, b)

#define TimeFunction() Timer CAT(t, __LINE__)(__func__)
#define TimeBlock(name) Timer CAT(t, __LINE__)(name)
