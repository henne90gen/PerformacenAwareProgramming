#include "json.h"

#include <iostream>
#include <sstream>
#include <string_view>

enum class JSON_TokenType {
    NONE = 0,
    DICT_OPEN = 1,
    DICT_CLOSE = 2,
    DICT_COLON = 3,
    ARR_OPEN = 4,
    ARR_CLOSE = 5,
    COMMA = 6,
    STRING = 7,
    INTEGER = 8,
    FLOAT = 9,
    TRUE = 10,
    FALSE = 11,
};

struct JSON_Token {
    JSON_TokenType type;
    char *start;
    int length;
};

struct JSON_Context {
    char *buf;
    int length;
    int cursor;

    void skipWhitespace();
    JSON_Token nextToken();
};

JSON_Node *
parse(JSON_Context &ctx);

namespace std {
    std::string to_string(JSON_TokenType type) {
        switch (type) {
            case JSON_TokenType::NONE:
                return "NONE";
            case JSON_TokenType::DICT_OPEN:
                return "DICT_OPEN";
            case JSON_TokenType::DICT_CLOSE:
                return "DICT_CLOSE";
            case JSON_TokenType::DICT_COLON:
                return "DICT_COLON";
            case JSON_TokenType::ARR_OPEN:
                return "ARR_OPEN";
            case JSON_TokenType::ARR_CLOSE:
                return "ARR_CLOSE";
            case JSON_TokenType::COMMA:
                return "COMMA";
            case JSON_TokenType::STRING:
                return "STRING";
            case JSON_TokenType::INTEGER:
                return "INTEGER";
            case JSON_TokenType::FLOAT:
                return "FLOAT";
            case JSON_TokenType::TRUE:
                return "TRUE";
            case JSON_TokenType::FALSE:
                return "FALSE";
            default:
                return "UNKNOWN";
        }
    }

    std::string to_string(const JSON_Token &token) {
        std::stringstream ss;
        ss << to_string(token.type);
        ss << ": '";
        ss << std::string(token.start, token.length);
        ss << "'";
        return ss.str();
    }
}

JSON_Array::~JSON_Array() {
    for (auto n: array) {
        delete n;
    }
}

JSON_Dictionary::~JSON_Dictionary() {
    for (auto entry: dictionary) {
        delete entry.second;
    }
}

void
JSON_Context::skipWhitespace() {
    while (buf[cursor] == '\t' || buf[cursor] == '\n' || buf[cursor] == '\r' || buf[cursor] == ' ') {
        cursor++;
    }
}

JSON_Token
JSON_Context::nextToken() {
    skipWhitespace();

    if (buf[cursor] == '{') {
        cursor++;
        return { JSON_TokenType::DICT_OPEN, buf + cursor - 1, 1 };
    }
    if (buf[cursor] == '}') {
        cursor++;
        return { JSON_TokenType::DICT_CLOSE, buf + cursor - 1, 1 };
    }
    if (buf[cursor] == ':') {
        cursor++;
        return { JSON_TokenType::DICT_COLON, buf + cursor - 1, 1 };
    }
    if (buf[cursor] == '[') {
        cursor++;
        return { JSON_TokenType::ARR_OPEN, buf + cursor - 1, 1 };
    }
    if (buf[cursor] == ']') {
        cursor++;
        return { JSON_TokenType::ARR_CLOSE, buf + cursor - 1, 1 };
    }
    if (buf[cursor] == ',') {
        cursor++;
        return { JSON_TokenType::COMMA, buf + cursor - 1, 1 };
    }

    if (buf[cursor] == '"') {
        auto startCursor = cursor;
        cursor++;
        while (buf[cursor] != '"') {
            cursor++;
        }
        cursor++;
        return { JSON_TokenType::STRING, buf + startCursor, cursor - startCursor };
    }

    if ((buf[cursor] >= '0' && buf[cursor] <= '9') || buf[cursor] == '-') {
        auto startCursor = cursor;
        if (buf[cursor] == '-') {
            cursor++;
        }

        while (buf[cursor] >= '0' && buf[cursor] <= '9') {
            cursor++;
        }

        if (buf[cursor] != '.') {
            cursor++;
            return { JSON_TokenType::INTEGER, buf + startCursor, cursor - startCursor };
        }

        cursor++;

        while (buf[cursor] >= '0' && buf[cursor] <= '9') {
            cursor++;
        }

        return { JSON_TokenType::FLOAT, buf + startCursor, cursor - startCursor };
    }

    if (std::string_view(buf + cursor, 4) == "true") {
        cursor += 4;
        return { JSON_TokenType::TRUE, buf + cursor - 4, 4 };
    }

    if (std::string_view(buf + cursor, 5) == "false") {
        cursor += 5;
        return { JSON_TokenType::FALSE, buf + cursor - 5, 5 };
    }

    return { JSON_TokenType::NONE, buf, 0 };
}

JSON_Node *
parseDict(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type != JSON_TokenType::DICT_OPEN) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    std::unordered_map<std::string, JSON_Node *> dictionary = {};

    while (true) {
        token = ctx.nextToken();
        if (token.type != JSON_TokenType::STRING) {
            for (auto entry: dictionary) {
                delete entry.second;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        auto key = std::string(token.start + 1, token.length - 2);

        token = ctx.nextToken();
        if (token.type != JSON_TokenType::DICT_COLON) {
            for (auto entry: dictionary) {
                delete entry.second;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        JSON_Node *value = parse(ctx);
        if (value == nullptr) {
            for (auto entry: dictionary) {
                delete entry.second;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        dictionary[key] = value;

        token = ctx.nextToken();
        if (token.type == JSON_TokenType::COMMA) {
            continue;
        }

        if (token.type == JSON_TokenType::DICT_CLOSE) {
            break;
        }

        for (auto entry: dictionary) {
            delete entry.second;
        }
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    return new JSON_Dictionary(dictionary);
}

JSON_Node *
parseArray(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type != JSON_TokenType::ARR_OPEN) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    std::vector<JSON_Node *> arr = {};
    while (true) {
        auto node = parse(ctx);
        if (node == nullptr) {
            for (auto n: arr) {
                delete n;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        arr.emplace_back(node);

        token = ctx.nextToken();
        if (token.type == JSON_TokenType::COMMA) {
            continue;
        }

        if (token.type == JSON_TokenType::ARR_CLOSE) {
            break;
        }

        for (auto n: arr) {
            delete n;
        }
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    return new JSON_Array(arr);
}

JSON_Node *
parseString(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type != JSON_TokenType::STRING) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    auto s = std::string(token.start + 1, token.length - 2);
    return new JSON_String(s);
}

JSON_Node *
parseInteger(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type != JSON_TokenType::INTEGER) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    i64 i = 0;
    try {
        i = std::stol(std::string(token.start, token.length));
    } catch (...) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }
    return new JSON_Integer(i);
}

JSON_Node *
parseFloat(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type != JSON_TokenType::FLOAT) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }

    f64 f = 0.0;
    try {
        f = std::stod(std::string(token.start, token.length));
    } catch (...) {
        ctx.cursor = cursorBefore;
        return nullptr;
    }
    return new JSON_Float(f);
}

JSON_Node *
parseBool(JSON_Context &ctx) {
    auto cursorBefore = ctx.cursor;

    auto token = ctx.nextToken();
    if (token.type == JSON_TokenType::TRUE) {
        return new JSON_Bool(true);
    }

    if (token.type == JSON_TokenType::FALSE) {
        return new JSON_Bool(false);
    }

    ctx.cursor = cursorBefore;
    return nullptr;
}

JSON_Node *
parse(JSON_Context &ctx) {
    JSON_Node *dict = parseDict(ctx);
    if (dict != nullptr) {
        return dict;
    }

    JSON_Node *arr = parseArray(ctx);
    if (arr != nullptr) {
        return arr;
    }

    JSON_Node *str = parseString(ctx);
    if (str != nullptr) {
        return str;
    }

    JSON_Node *i = parseInteger(ctx);
    if (i != nullptr) {
        return i;
    }

    JSON_Node *f = parseFloat(ctx);
    if (f != nullptr) {
        return f;
    }

    JSON_Node *b = parseBool(ctx);
    if (b != nullptr) {
        return b;
    }

    return nullptr;
}

static void
indent(int level) {
    for (int i = 0; i < level; i++) {
        std::cout << "  ";
    }
}

void
traverse(JSON_Node *node, int level = 0) {
    indent(level);

    if (node->type == JSON_NodeType::DICTIONARY) {
        std::cout << "DICTIONARY: {" << std::endl;

        bool isFirst = true;
        for (auto n: ((JSON_Dictionary *) node)->dictionary) {
            if (!isFirst) {
                std::cout << std::endl;
            } else {
                isFirst = false;
            }

            indent(level + 1);
            std::cout << n.first << std::endl;

            indent(level + 1);
            std::cout << ":" << std::endl;

            traverse(n.second, level + 1);
        }

        indent(level);
        std::cout << "}" << std::endl;
        return;
    }

    if (node->type == JSON_NodeType::ARRAY) {
        std::cout << "ARRAY: [" << std::endl;

        for (auto n: ((JSON_Array *) node)->array) {
            traverse(n, level + 1);
        }

        indent(level);
        std::cout << "]" << std::endl;
        return;
    }

    if (node->type == JSON_NodeType::STRING) {
        std::cout << "STRING: " << ((JSON_String *) node)->s << std::endl;
        return;
    }

    if (node->type == JSON_NodeType::INTEGER) {
        std::cout << "INTEGER: " << ((JSON_Integer *) node)->i << std::endl;
        return;
    }

    if (node->type == JSON_NodeType::FLOAT) {
        std::cout << "FLOAT: " << ((JSON_Float *) node)->f << std::endl;
        return;
    }

    if (node->type == JSON_NodeType::BOOLEAN) {
        std::cout << "BOOL: " << ((JSON_Bool *) node)->b << std::endl;
        return;
    }
}

JSON_Node *
parseJSON(char *buf, int length) {
    JSON_Context ctx = { buf, length, 0 };
    auto root = parse(ctx);
    if (root == nullptr) {
        std::cerr << "failed to parse JSON" << std::endl;
        return nullptr;
    }

    return root;
}
