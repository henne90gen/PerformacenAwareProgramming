#include <fstream>
#include <vector>

#include "shared.h"

struct PointPair
{
    f64 x0, y0, x1, y1;
};

std::vector<PointPair> parsePointPairs(std::ifstream &f)
{
    return {};
}

int main()
{
    auto f = std::ifstream("../point_pairs.json");
    auto pointPairs = parsePointPairs(f);
    return 0;
}
