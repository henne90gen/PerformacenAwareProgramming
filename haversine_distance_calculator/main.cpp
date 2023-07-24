#include <cmath>
#include <fstream>
#include <iomanip>
#include <iostream>

#include "haversine.h"
#include "metrics.h"

std::pair<std::vector<f64>, std::vector<PointPair>>
parseAnswers() {
    std::vector<f64> answers = {};
    std::vector<PointPair> pointPairs = {};
    auto f = std::ifstream("../answers.txt");
    double a, b, c, d, e;
    while (f >> a >> b >> c >> d >> e) {
        answers.push_back(e);
        pointPairs.push_back({ a, b, c, d });
    }
    return { answers, pointPairs };
}

bool
approximatelyEqual(f64 a, f64 b, f64 epsilon) {
    return std::fabs(a - b) <= ((std::fabs(a) < std::fabs(b) ? std::fabs(b) : std::fabs(a)) * epsilon);
}

void
verifyAnswers(const std::vector<f64> &distances) {
    auto [answers, answerPointPairs] = parseAnswers();
    if (answers.empty()) {
        std::cout << "failed to parse answers" << std::endl;
        return;
    }

    bool failure = false;
    auto epsilon = 0.0000001;
    auto expectedDistanceAverage = 0.0;
    auto distanceAverage = 0.0;
    for (int i = 0; i < distances.size(); i++) {
        const auto &distance = distances[i];
        const auto &expectedDistance = answers[i];
        if (!approximatelyEqual(distance, expectedDistance, epsilon)) {
            std::cout << "failed " << distance << " != " << expectedDistance << std::endl;
            failure = true;
        }
        distanceAverage += distance;
        expectedDistanceAverage += expectedDistance;
    }

    distanceAverage /= static_cast<f64>(distances.size());
    expectedDistanceAverage /= static_cast<f64>(distances.size());
    if (!approximatelyEqual(distanceAverage, expectedDistanceAverage, epsilon)) {
        std::cout << "average failed " << distanceAverage << " != " << expectedDistanceAverage << std::endl;
        failure = true;
    }

    if (failure) {
        std::cout << "failure" << std::endl;
    } else {
        std::cout << "success" << std::endl;
    }
}

void
increasePrecisionOfFloatPrintout() {
    std::cout << std::fixed << std::setw(18) << std::setprecision(18);
    std::cout << "" << std::endl;
}

void
PrintTiming(const std::string &name, u64 cpuTimerFreq, u64 totalElapsedF64, u64 start, u64 end) {
    auto diff = end - start;
    auto parseTime = CPUTimerDiffToNanoseconds(diff, cpuTimerFreq);
    auto parsePercentage = static_cast<f64>(diff) / totalElapsedF64 * 100.0;
    std::cout << name << ": " << parsePercentage << "% " << parseTime << "ns" << std::endl;
}

int
main() {
    // increasePrecisionOfFloatPrintout();

    auto start = ReadCPUTimer();

    auto buf = ReadFile("../point_pairs.json");
    if (buf.data == nullptr) {
        return 1;
    }

    auto afterRead = ReadCPUTimer();

    auto pointPairs = ParsePointPairsGeneric(buf);
    if (pointPairs.empty()) {
        return 1;
    }

    auto afterParsing = ReadCPUTimer();

    auto distances = CalculateHaversineDistances(pointPairs);

    auto afterDistanceCalculation = ReadCPUTimer();

    auto totalElapsed = afterDistanceCalculation - start;
    auto totalElapsedF64 = static_cast<f64>(totalElapsed);
    auto cpuTimerFreq = EstimateCPUTimerFrequency();

    PrintTiming("reading file", cpuTimerFreq, totalElapsedF64, start, afterRead);
    PrintTiming("parsing", cpuTimerFreq, totalElapsedF64, afterRead, afterParsing);
    PrintTiming("calculating", cpuTimerFreq, totalElapsedF64, afterParsing, afterDistanceCalculation);

    verifyAnswers(distances);

    return 0;
}
