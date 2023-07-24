#include <cstring>
#include <fstream>
#include <iostream>
#include <string>
#include <string_view>
#include <unordered_map>
#include <vector>

#include "haversine.h"
#include "json.h"
#include "profiler.h"
#include "reference.h"

const f64 EARTHS_RADIUS = 6372.8;

std::vector<PointPair>
ParsePointPairsCustom(const Buffer &buffer) {
    auto result = std::vector<PointPair>();
    f64 currentPair[4] = {};
    int currentNum = 0;
    int numberStartIndex = 0;
    bool pairsHaveStarted = false;
    for (int i = 0; i < buffer.size; i++) {
        if (buffer.data[i] == ']') {
            break;
        }

        if (buffer.data[i] == ' ' || buffer.data[i] == '\t' || buffer.data[i] == '\n') {
            continue;
        }

        if (!pairsHaveStarted) {
            if (buffer.data[i] == '{') {
                continue;
            }
            if (std::string_view(buffer.data + i, 7) == "\"pairs\"") {
                i += 6;
                continue;
            }
            if (buffer.data[i] == ':') {
                continue;
            }
            if (buffer.data[i] == '[') {
                pairsHaveStarted = true;
                continue;
            }
        }

        if (buffer.data[i] == '{') {
            if (currentNum > 0) {
                result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);
            }
            std::memset(currentPair, 0.0, 4);
            continue;
        }

        if (buffer.data[i] == '"') {
            if (std::string_view(buffer.data + i + 1, 2) == "x0") {
                currentNum = 0;
            }
            if (std::string_view(buffer.data + i + 1, 2) == "y0") {
                currentNum = 1;
            }
            if (std::string_view(buffer.data + i + 1, 2) == "x1") {
                currentNum = 2;
            }
            if (std::string_view(buffer.data + i + 1, 2) == "y1") {
                currentNum = 3;
            }
            i += 4;
            numberStartIndex = i + 1;
            continue;
        }

        if (buffer.data[i] == ',' || buffer.data[i] == '}') {
            auto strLength = i - numberStartIndex;
            auto strToParse = std::string(buffer.data + numberStartIndex, strLength);
            auto d = std::stod(strToParse);
            currentPair[currentNum] = d;
            if (buffer.data[i] == '}') {
                // skip the comma after the '}' as well
                i++;
            }
            continue;
        }
    }

    result.emplace_back(currentPair[0], currentPair[1], currentPair[2], currentPair[3]);

    return result;
}

Buffer
ReadFile(const std::string &path) {
    TimeFunction();

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

    return { buf, (u64) fileSize };
}

std::vector<PointPair>
ParsePointPairsGeneric(const Buffer &buf) {
    TimeFunction();

    auto root = JSON::Parse(buf.data, buf.size);
    if (root == nullptr) {
        return {};
    }

    if (root->type != JSON::NodeType::DICTIONARY) {
        delete root;
        return {};
    }

    auto dict = (JSON::Dictionary *) root;
    auto itr = dict->dictionary.find("pairs");
    if (itr == dict->dictionary.end()) {
        delete root;
        return {};
    }

    auto node = itr->second;
    if (node->type != JSON::NodeType::ARRAY) {
        delete root;
        return {};
    }

    std::vector<PointPair> result = {};
    auto arr = (JSON::Array *) node;
    for (auto n: arr->array) {
        TimeBlock("Hello");
        if (n->type != JSON::NodeType::DICTIONARY) {
            delete root;
            return {};
        }

        auto d = (JSON::Dictionary *) n;
        auto itrX0 = d->dictionary.find("x0");
        if (itrX0 == d->dictionary.end() || itrX0->second->type != JSON::NodeType::FLOAT) {
            delete root;
            return {};
        }

        auto itrY0 = d->dictionary.find("y0");
        if (itrY0 == d->dictionary.end() || itrY0->second->type != JSON::NodeType::FLOAT) {
            delete root;
            return {};
        }

        auto itrX1 = d->dictionary.find("x1");
        if (itrX1 == d->dictionary.end() || itrX1->second->type != JSON::NodeType::FLOAT) {
            delete root;
            return {};
        }

        auto itrY1 = d->dictionary.find("y1");
        if (itrY1 == d->dictionary.end() || itrY1->second->type != JSON::NodeType::FLOAT) {
            delete root;
            return {};
        }

        result.emplace_back(((JSON::Float *) itrX0->second)->f, ((JSON::Float *) itrY0->second)->f, ((JSON::Float *) itrX1->second)->f, ((JSON::Float *) itrY1->second)->f);
    }

    delete root;
    return result;
}

std::vector<f64>
CalculateHaversineDistances(const std::vector<PointPair> &pointPairs) {
    TimeFunction();

    auto result = std::vector<f64>();
    result.reserve(pointPairs.size());

    for (const auto &pointPair: pointPairs) {
        auto distance = ReferenceHaversine(pointPair.x0, pointPair.y0, pointPair.x1, pointPair.y1, EARTHS_RADIUS);
        result.push_back(distance);
    }

    return result;
}
