#include "haversine.h"

#include <fstream>
#include <iomanip>
#include <iostream>

std::vector<f64>
parseAnswers() {
    std::vector<f64> answers = {};
    auto f = std::ifstream("../../answers.txt");
    double a;
    while (f >> a) {
        answers.push_back(a);
    }
    return answers;
}


int
main() {
    std::cout << std::fixed << std::setw(18) << std::setprecision(18);

    auto pointPairs = ParsePointPairs("../../point_pairs.json");

    std::cout << "parsed " << pointPairs.size() << " pairs" << std::endl;

    auto distances = CalculateHaversineDistances(pointPairs);

    auto answers = parseAnswers();
    auto itr = answers.begin() + (answers.size() - 1);
    auto expectedDistanceAverage = *itr;

    std::cout << "// GIVEN" << std::endl
              << "std::vector<PointPair> pointPairs = {" << std::endl;
    for (const auto &pair: pointPairs) {
        std::cout << "{" << pair.x0 << ", " << pair.y0 << ", " << pair.x1 << ", " << pair.y1 << "}," << std::endl;
    }
    std::cout << "};" << std::endl
              << std::endl
              << "// WHEN" << std::endl
              << "const auto result = CalculateHaversineDistances(pointPairs);" << std::endl
              << std::endl
              << "// THEN" << std::endl
              << "EXPECT_EQ(result.size(), " << distances.size() << ");" << std::endl;

    for (int i = 0; i < distances.size(); i++) {
        const auto &distance = distances[i];
        const auto &expectedDistance = answers[i];
        std::cout << "EXPECT_FLOAT_EQ(result[" << i << "], " << expectedDistance << ");" << std::endl;
    }

    return 0;
}
