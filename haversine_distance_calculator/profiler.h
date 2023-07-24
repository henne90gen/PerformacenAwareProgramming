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

struct Timer {
    u64 parentIndex = 0;
    std::string label = {};
    u64 start = 0;

    Timer(const std::string &name);
    ~Timer();
};

struct TimeAggregate {
    std::string label = {};
    u64 elapsed = 0;
    u64 elapsedInChildren = 0;
};

struct Profiler {
    u64 start = 0;
    u32 nextAggregateIndex = 1;
    std::unordered_map<std::string, u32> timeAggregateIndices = {};
    std::array<TimeAggregate, 4096> timeAggregates = {};

    Profiler();
};

extern Profiler GlobalProfiler;

void BeginProfiling();
void EndProfiling();

#define CAT_(a, b) a##b
#define CAT(a, b) CAT_(a, b)

#define TimeFunction() Timer CAT(t, __LINE__)(__func__)
#define TimeBlock(name) Timer CAT(t, __LINE__)(name)
