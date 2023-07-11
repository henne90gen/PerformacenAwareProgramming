#include <cstring>
#include <fstream>
#include <iostream>
#include <string>
#include <string_view>
#include <vector>

#include "shared.h"

struct PointPair {
    f64 x0, y0, x1, y1;
    PointPair(f64 _x0, f64 _y0, f64 _x1, f64 _y1)
      : x0(_x0), y0(_y0), x1(_x1), y1(_y1) {}
};

std::vector<PointPair>
parsePointPairs(const std::string &path) {
    std::cout << "parsing point pairs" << std::endl;

    auto f = std::ifstream(path, std::ios::binary | std::ios::ate);
    if (!f.is_open()) {
        std::cerr << "failed to open file" << std::endl;
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
        if (std::string_view(buf + i, 2) == "}]") {
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
                continue;
            }
            if (buf[i] == ':') {
                pairsHaveStarted = true;
                continue;
            }
        }

        if (buf[i] == '{') {
            if (currentNum >0 ) {
                result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);
                std::cout << currentPair[0] << ", " << currentPair[1] << ", "
                          << currentPair[2] << ", " << currentPair[3] << std::endl;
            }
            std::cout << "start of object" << std::endl;
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
            std::cout << "reading number " << currentNum << std::endl;
            i += 4;
            numberStartIndex = i + 1;
            continue;
        }
        if (buf[i] == ',' || buf[i] == '}') {
            auto strLength = i - numberStartIndex;
            auto strToParse = std::string(buf + numberStartIndex, strLength);
            auto d = std::stod(strToParse);
            std::cout << "string to parse " << strToParse << " " << d << std::endl;
            currentPair[currentNum] = d;
            if (buf[i] == '}') {
                // skip the comma after the '}' as well
                i++;
            }
            continue;
        }
        std::cout << buf[i] << std::endl;
    }

    result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);

    std::cout << "parsed " << result.size() << " pairs" << std::endl;

    for (const auto &pair: result) {
        std::cout << "x0=" << pair.x0 << " y0=" << pair.y0 << " x1=" << pair.x1 << " y1=" << pair.y1 << std::endl;
    }

    return result;
}

int
main() {
    auto pointPairs = parsePointPairs("../point_pairs.json");
    return 0;
}
