#include "profiler.h"

#include <iomanip>
#include <iostream>

Profiler GlobalProfiler = {};
u32 GlobalProfilerParentIndex = 0;

#if _WIN32

#include <windows.h>

static u64
GetOSTimerFreq(void) {
    LARGE_INTEGER Freq;
    QueryPerformanceFrequency(&Freq);
    return Freq.QuadPart;
}

static u64
ReadOSTimer(void) {
    LARGE_INTEGER Value;
    QueryPerformanceCounter(&Value);
    return Value.QuadPart;
}

#else

#include <sys/time.h>

static u64
GetOSTimerFreq(void) {
    return 1000000;
}

static u64
ReadOSTimer(void) {
    // NOTE(casey): The "struct" keyword is not necessary here when compiling in C++,
    // but just in case anyone is using this file from C, I include it.
    struct timeval Value;
    gettimeofday(&Value, 0);

    u64 Result = GetOSTimerFreq() * (u64) Value.tv_sec + (u64) Value.tv_usec;
    return Result;
}

#endif

u64
EstimateCPUTimerFrequency() {
    u64 MillisecondsToWait = 100;

    u64 OSEnd = 0;
    u64 OSElapsed = 0;
    u64 OSFreq = GetOSTimerFreq();
    u64 OSWaitTime = OSFreq * MillisecondsToWait / 1000;
    u64 CPUStart = ReadCPUTimer();
    u64 OSStart = ReadOSTimer();

    while (OSElapsed < OSWaitTime) {
        OSEnd = ReadOSTimer();
        OSElapsed = OSEnd - OSStart;
    }

    u64 CPUEnd = ReadCPUTimer();
    u64 CPUElapsed = CPUEnd - CPUStart;
    return OSFreq * CPUElapsed / OSElapsed;
}

u64
CPUTimerDiffToNanoseconds(u64 cpuTimer, u64 cpuTimerFreq) {
    return cpuTimer / (static_cast<f64>(cpuTimerFreq) / 1000000000.0);
}

void
BeginProfiling() {
    GlobalProfiler = {};
}

Profiler::Profiler() {
    start = ReadCPUTimer();
}

void
PrintTiming(u64 cpuTimerFreq, u64 totalElapsedF64, const TimeAggregate &aggregate) {
    auto elapsed = aggregate.elapsedWithoutChildren;
    auto timeNs = CPUTimerDiffToNanoseconds(elapsed, cpuTimerFreq);
    auto timeMs = static_cast<f64>(timeNs) / 1000000.0;
    auto percentage = static_cast<f64>(elapsed) / totalElapsedF64 * 100.0;
    std::cout << std::left << std::setw(35) << aggregate.label
              << std::right << std::setw(10) << aggregate.hitCount
              << std::fixed << std::setprecision(3) << std::right << std::setw(12) << timeMs << "ms "
              << std::fixed << std::setprecision(2) << std::right << std::setw(7) << percentage << "% ";
    if (aggregate.elapsedWithChildren != elapsed) {
        auto percentageWithChildren = static_cast<f64>(aggregate.elapsedWithChildren) / totalElapsedF64 * 100.0;
        std::cout << std::fixed << std::setprecision(2) << std::right << std::setw(6) << percentageWithChildren << "% ";
    }
    std::cout << std::endl;
}

void
EndProfiling() {
    auto end = ReadCPUTimer();
    auto totalElapsed = end - GlobalProfiler.start;
    auto totalElapsedF64 = static_cast<f64>(totalElapsed);
    auto cpuTimerFreq = EstimateCPUTimerFrequency();

    std::cout << std::left << std::setw(35) << "Name"
              << std::right << std::setw(10) << "Hits"
              << std::right << std::setw(14) << "Time"
              << std::right << std::setw(9) << "Percent"
              << " Percent with Children" << std::endl;
    for (int i = 0; i < GlobalProfiler.timeAggregates.size(); i++) {
        const auto &aggregate = GlobalProfiler.timeAggregates[i];
        if (aggregate.hitCount == 0) {
            continue;
        }

        PrintTiming(cpuTimerFreq, totalElapsedF64, aggregate);
    }

    auto timeMs = totalElapsedF64 / 1000000.0;
    std::cout << std::left << std::setw(45) << "Total" << std::right << std::fixed << std::setprecision(3) << std::setw(12) << timeMs << "ms" << std::endl;
}
