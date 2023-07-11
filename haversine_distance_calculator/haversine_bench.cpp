#include <benchmark/benchmark.h>
#include <fstream>

#include "haversine.h"
#include "reference.h"

static void
BM_ParsePointsCustom(benchmark::State &state) {
    auto fileName = tmpnam(nullptr);
    auto f = std::ofstream(fileName);
    f << "{\"pairs\":[{\"x0\":60.06493178495121,\"y0\":-31.102444435872737,\"x1\":35.67709437177292,\"y1\":-5.85944106452626},{\"x0\":-22.66247210298143,\"y0\":24.830077100813543,\"x1\":-30.196277871336306,\"y1\":28.656104157990143},{\"x0\":118.38421622135718,\"y0\":34.6507820744768,\"x1\":74.25540994301548,\"y1\":41.87667158056794},{\"x0\":92.4134749902354,\"y0\":-51.20604642909591,\"x1\":93.08391243761406,\"y1\":1.3557561086926952},{\"x0\":118.31794359672764,\"y0\":48.46106212094184,\"x1\":105.54768962904417,\"y1\":37.50793416786314},{\"x0\":62.68206228790001,\"y0\":13.672930981181153,\"x1\":53.27199043744885,\"y1\":-36.82027677213894},{\"x0\":-32.325591719735314,\"y0\":37.49429043316076,\"x1\":-54.80208016806659,\"y1\":32.29890339193499},{\"x0\":23.040599244220896,\"y0\":-42.36340322527565,\"x1\":-22.663910633595968,\"y1\":-24.124114256936224},{\"x0\":87.07687914735949,\"y0\":37.07801016798828,\"x1\":2.1940648570643546,\"y1\":22.04778006853196},{\"x0\":92.88401609477133,\"y0\":-45.33555007046114,\"x1\":85.28880683245166,\"y1\":-45.91752245107136}]}";
    f.close();

    for (auto _: state) {
        auto pointPairs = ParsePointPairs(fileName);
        benchmark::DoNotOptimize(pointPairs);
    }
}
BENCHMARK(BM_ParsePointsCustom);

static void
BM_ReferenceDistanceCalculator(benchmark::State &state) {
    for (auto _: state) {
        auto distance = ReferenceHaversine(0.0, 0.0, 0.0, 0.0, 0.0);
        benchmark::DoNotOptimize(distance);
    }
}
BENCHMARK(BM_ReferenceDistanceCalculator);

BENCHMARK_MAIN();
