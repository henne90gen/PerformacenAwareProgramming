#include "shared.h"

#if _WIN32
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
