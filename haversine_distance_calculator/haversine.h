#pragma once

#include <string>
#include <vector>

#include "shared.h"


struct Buffer {
    char *data;
    u64 size;
};

Buffer ReadFile(const std::string &path);

struct PointPair {
    f64 x0, y0, x1, y1;
    PointPair(f64 _x0, f64 _y0, f64 _x1, f64 _y1)
      : x0(_x0), y0(_y0), x1(_x1), y1(_y1) {}
};

std::vector<PointPair>
ParsePointPairsCustom(const Buffer &buffer);
std::vector<PointPair>
ParsePointPairsGeneric(const Buffer &buffer);

std::vector<f64> CalculateHaversineDistances(const std::vector<PointPair> &pointPairs);
