#include <cstring>
#include <fstream>
#include <iostream>
#include <string>
#include <string_view>
#include <vector>

#include "haversine.h"
#include "reference.h"

const f64 EARTHS_RADIUS = 6372.8;

std::vector<PointPair>
ParsePointPairs(const std::string &path) {
    auto f = std::ifstream(path, std::ios::binary | std::ios::ate);
    if (!f.is_open()) {
        std::cerr << "failed to open '" << path << "'" << std::endl;
        return {};
    }

    const auto fileSize = f.tellg();
    f.seekg(0, std::ios::beg);

    auto buf = (char *) malloc(fileSize);
    if (!f.read(buf, fileSize)) {
        std::cerr << "failed to read data" << std::endl;
        return {};
    }

    auto result = std::vector<PointPair>();
    f64 currentPair[4] = {};
    int currentNum = 0;
    int numberStartIndex = 0;
    bool pairsHaveStarted = false;
    for (int i = 0; i < fileSize; i++) {
        if (buf[i] == ']') {
            break;
        }

        if (buf[i] == ' ' || buf[i] == '\t' || buf[i] == '\n') {
            continue;
        }

        if (!pairsHaveStarted) {
            if (buf[i] == '{') {
                continue;
            }
            if (std::string_view(buf + i, 7) == "\"pairs\"") {
                i += 6;
                continue;
            }
            if (buf[i] == ':') {
                continue;
            }
            if (buf[i] == '[') {
                pairsHaveStarted = true;
                continue;
            }
        }

        if (buf[i] == '{') {
            if (currentNum > 0) {
                result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);
            }
            std::memset(currentPair, 0.0, 4);
            continue;
        }

        if (buf[i] == '"') {
            if (std::string_view(buf + i + 1, 2) == "x0") {
                currentNum = 0;
            }
            if (std::string_view(buf + i + 1, 2) == "y0") {
                currentNum = 1;
            }
            if (std::string_view(buf + i + 1, 2) == "x1") {
                currentNum = 2;
            }
            if (std::string_view(buf + i + 1, 2) == "y1") {
                currentNum = 3;
            }
            i += 4;
            numberStartIndex = i + 1;
            continue;
        }

        if (buf[i] == ',' || buf[i] == '}') {
            auto strLength = i - numberStartIndex;
            auto strToParse = std::string(buf + numberStartIndex, strLength);
            auto d = std::stod(strToParse);
            currentPair[currentNum] = d;
            if (buf[i] == '}') {
                // skip the comma after the '}' as well
                i++;
            }
            continue;
        }
    }

    result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);

    return result;
}

std::vector<f64>
CalculateHaversineDistances(const std::vector<PointPair> &pointPairs) {
    auto result = std::vector<f64>();
    result.reserve(pointPairs.size());

    for (const auto &pointPair: pointPairs) {
        auto distance = ReferenceHaversine(pointPair.x0, pointPair.y0, pointPair.x1, pointPair.y1, EARTHS_RADIUS);
        result.push_back(distance);
    }

    return result;
}
