#include "profiler.h"

#include <iomanip>
#include <iostream>

Profiler GlobalProfiler = {};

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

Timer::Timer(const std::string &name) : name(name) {
    start = ReadCPUTimer();
}

Timer::~Timer() {
    auto end = ReadCPUTimer();
    GlobalProfiler.measurements.emplace_back(name, start, end);
}

void
BeginProfiling() {
    GlobalProfiler = {};
}

Profiler::Profiler() {
    start = ReadCPUTimer();
}

void
PrintTiming(const std::string &name, u64 cpuTimerFreq, u64 totalElapsedF64, u64 start, u64 end) {
    auto diff = end - start;
    auto parseTimeNs = CPUTimerDiffToNanoseconds(diff, cpuTimerFreq);
    auto parseTimeMs = static_cast<f64>(parseTimeNs) / 1000000.0;
    auto parsePercentage = static_cast<f64>(diff) / totalElapsedF64 * 100.0;
    std::cout << std::left << std::setw(35) << name;
    std::cout << std::fixed << std::setprecision(2) << std::right << std::setw(5) << parsePercentage << "% ";
    std::cout << std::fixed << std::setprecision(3) << std::right << std::setw(10) << parseTimeMs << "ms" << std::endl;
}

void
EndProfiling() {
    auto end = ReadCPUTimer();
    auto totalElapsed = end - GlobalProfiler.start;
    auto totalElapsedF64 = static_cast<f64>(totalElapsed);
    auto cpuTimerFreq = EstimateCPUTimerFrequency();

    for (const auto &timePair: GlobalProfiler.measurements) {
        PrintTiming(timePair.name, cpuTimerFreq, totalElapsedF64, timePair.start, timePair.end);
    }
}
