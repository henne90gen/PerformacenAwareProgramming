#include "shared.h"

#include <string>
#include <vector>

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
    std::string name;
    u64 start = 0;

    Timer(const std::string &name);
    ~Timer();
};

struct TimePairs {
    std::string name = {};
    u64 start = 0;
    u64 end = 0;
};

struct Profiler {
    u64 start = 0;
    std::vector<TimePairs> measurements = {};

    Profiler();
};

extern Profiler GlobalProfiler;

void BeginProfiling();
void EndProfiling();

#define CAT_(a, b) a ## b
#define CAT(a, b) CAT_(a, b)

#define TimeFunction() Timer CAT(t, __LINE__)(__func__)
#define TimeBlock(name) Timer CAT(t, __LINE__)(name)
