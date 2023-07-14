cmake_minimum_required(VERSION 3.20)

project(haversine)

include(FetchContent)

FetchContent_Declare(
        googletest
        GIT_REPOSITORY https://github.com/google/googletest.git
        GIT_TAG b796f7d44681514f58a683a3a71ff17c94edb0c1 # v1.13.0
)
FetchContent_Declare(
        benchmark
        GIT_REPOSITORY https://github.com/google/benchmark.git
        GIT_TAG d572f4777349d43653b21d6c2fc63020ab326db2 # v1.7.1
)

set(BENCHMARK_ENABLE_TESTING OFF CACHE BOOL "" FORCE) # benchmark
set(BENCHMARK_ENABLE_INSTALL OFF CACHE BOOL "" FORCE) # benchmark
set(BENCHMARK_INSTALL_DOCS OFF CACHE BOOL "" FORCE) # benchmark
set(BENCHMARK_ENABLE_GTEST_TESTS OFF CACHE BOOL "" FORCE) # benchmark

FetchContent_MakeAvailable(
    googletest
    benchmark
)

set(CMAKE_CXX_STANDARD 20)

add_library(haversine haversine.cpp json.cpp reference.cpp)

add_executable(haversine_exe main.cpp)
target_link_libraries(haversine_exe PRIVATE haversine)

add_executable(haversine_bench haversine_bench.cpp)
target_link_libraries(haversine_bench PRIVATE haversine benchmark::benchmark)

enable_testing()
include(GoogleTest)
add_executable(haversine_test haversine_test.cpp)
target_link_libraries(haversine_test PRIVATE haversine GTest::gtest_main)
gtest_discover_tests(haversine_test)
