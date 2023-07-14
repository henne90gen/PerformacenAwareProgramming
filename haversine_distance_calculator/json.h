#pragma once

#include <string>
#include <unordered_map>
#include <vector>

#include "shared.h"

enum class JSON_NodeType {
    DICTIONARY = 0,
    ARRAY = 1,
    STRING = 2,
    INTEGER = 3,
    FLOAT = 4,
    BOOLEAN = 5,
};

struct JSON_Node {
    JSON_NodeType type;

    explicit JSON_Node(JSON_NodeType t) : type(t) {}
};

struct JSON_Dictionary : public JSON_Node {
    std::unordered_map<std::string, JSON_Node *> dictionary;

    explicit JSON_Dictionary(std::unordered_map<std::string, JSON_Node *> d) : JSON_Node(JSON_NodeType::DICTIONARY), dictionary(d) {}
    ~JSON_Dictionary();
};

struct JSON_Array : public JSON_Node {
    std::vector<JSON_Node *> array;

    explicit JSON_Array(std::vector<JSON_Node *> a) : JSON_Node(JSON_NodeType::ARRAY), array(a) {}
    ~JSON_Array();
};

struct JSON_String : public JSON_Node {
    std::string s;

    explicit JSON_String(std::string s) : JSON_Node(JSON_NodeType::STRING), s(s) {}
};

struct JSON_Integer : public JSON_Node {
    int64_t i;

    explicit JSON_Integer(int64_t i) : JSON_Node(JSON_NodeType::INTEGER), i(i) {}
};

struct JSON_Float : public JSON_Node {
    f64 f;

    explicit JSON_Float(f64 f) : JSON_Node(JSON_NodeType::FLOAT), f(f) {}
};

struct JSON_Bool : public JSON_Node {
    bool b;

    explicit JSON_Bool(bool b) : JSON_Node(JSON_NodeType::BOOLEAN), b(b) {}
};

JSON_Node *parseJSON(char *buf, int length);
