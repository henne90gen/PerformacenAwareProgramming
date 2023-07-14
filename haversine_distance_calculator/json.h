#pragma once

#include <string>
#include <unordered_map>
#include <vector>

#include "shared.h"

namespace JSON {
    enum class NodeType {
        DICTIONARY = 0,
        ARRAY = 1,
        STRING = 2,
        INTEGER = 3,
        FLOAT = 4,
        BOOLEAN = 5,
    };

    struct Node {
        NodeType type;

        explicit Node(NodeType t) : type(t) {}
    };

    struct Dictionary : public Node {
        std::unordered_map<std::string, Node *> dictionary;

        explicit Dictionary(std::unordered_map<std::string, Node *> d) : Node(NodeType::DICTIONARY), dictionary(d) {}
        ~Dictionary();
    };

    struct Array : public Node {
        std::vector<Node *> array;

        explicit Array(std::vector<Node *> a) : Node(NodeType::ARRAY), array(a) {}
        ~Array();
    };

    struct String : public Node {
        std::string s;

        explicit String(std::string s) : Node(NodeType::STRING), s(s) {}
    };

    struct Integer : public Node {
        int64_t i;

        explicit Integer(int64_t i) : Node(NodeType::INTEGER), i(i) {}
    };

    struct Float : public Node {
        f64 f;

        explicit Float(f64 f) : Node(NodeType::FLOAT), f(f) {}
    };

    struct Bool : public Node {
        bool b;

        explicit Bool(bool b) : Node(NodeType::BOOLEAN), b(b) {}
    };

    Node *Parse(char *buf, int length);
}
