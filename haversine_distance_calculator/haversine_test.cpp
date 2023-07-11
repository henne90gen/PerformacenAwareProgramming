#include <gtest/gtest.h>

#include <fstream>

#include "haversine.h"
#include "reference.h"

std::pair<std::vector<f64>, std::vector<PointPair>>
parseAnswers() {
    std::vector<f64> answers = {};
    std::vector<PointPair> pointPairs = {};
    auto f = std::ifstream("../../answers.txt");
    double a, b, c, d, e;
    while (f >> a >> b >> c >> d >> e) {
        answers.push_back(e);
        pointPairs.push_back({ a, b, c, d });
    }
    return { answers, pointPairs };
}

TEST(HaversineTest, testWithPointParsing) {
    auto pointPairs = ParsePointPairs("../../point_pairs.json");

    auto result = parseAnswers();
    auto answers = result.first;

    auto distances = CalculateHaversineDistances(pointPairs);

    for (int i = 0; i < distances.size(); i++) {
        const auto &distance = distances[i];
        const auto &expectedDistance = answers[i];
        ASSERT_FLOAT_EQ(expectedDistance, distance);
    }
}

TEST(HaversineTest, testJustPointParsing) {
    auto pointPairs = ParsePointPairs("../../point_pairs.json");

    auto result = parseAnswers();
    auto expectedPointPairs = result.second;

    for (int i = 0; i < pointPairs.size(); i++) {
        const auto &pointPair = pointPairs[i];
        const auto &expectedPointPair = expectedPointPairs[i];
        ASSERT_FLOAT_EQ(pointPair.x0, expectedPointPair.x0);
        ASSERT_FLOAT_EQ(pointPair.y0, expectedPointPair.y0);
        ASSERT_FLOAT_EQ(pointPair.x1, expectedPointPair.x1);
        ASSERT_FLOAT_EQ(pointPair.y1, expectedPointPair.y1);
    }
}

TEST(HaversineTest, testWithoutPointParsing) {
    auto result = parseAnswers();
    auto answers = result.first;
    auto pointPairs = result.second;

    auto distances = CalculateHaversineDistances(pointPairs);

    for (int i = 0; i < distances.size(); i++) {
        const auto &distance = distances[i];
        const auto &expectedDistance = answers[i];
        ASSERT_FLOAT_EQ(expectedDistance, distance);
    }
}
