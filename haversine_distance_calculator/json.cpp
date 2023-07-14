#include "json.h"

#include <iostream>
#include <sstream>
#include <string_view>

namespace JSON {
    enum class TokenType {
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

    struct Token {
        TokenType type;
        char *start;
        int length;
    };

    struct Context {
        char *buf;
        int length;
        int cursor;

        void skipWhitespace();
        Token nextToken();
    };

    Node *
    parseInternal(Context &ctx);
}

namespace std {
    std::string to_string(JSON::TokenType type) {
        switch (type) {
            case JSON::TokenType::NONE:
                return "NONE";
            case JSON::TokenType::DICT_OPEN:
                return "DICT_OPEN";
            case JSON::TokenType::DICT_CLOSE:
                return "DICT_CLOSE";
            case JSON::TokenType::DICT_COLON:
                return "DICT_COLON";
            case JSON::TokenType::ARR_OPEN:
                return "ARR_OPEN";
            case JSON::TokenType::ARR_CLOSE:
                return "ARR_CLOSE";
            case JSON::TokenType::COMMA:
                return "COMMA";
            case JSON::TokenType::STRING:
                return "STRING";
            case JSON::TokenType::INTEGER:
                return "INTEGER";
            case JSON::TokenType::FLOAT:
                return "FLOAT";
            case JSON::TokenType::TRUE:
                return "TRUE";
            case JSON::TokenType::FALSE:
                return "FALSE";
            default:
                return "UNKNOWN";
        }
    }

    std::string to_string(const JSON::Token &token) {
        std::stringstream ss;
        ss << to_string(token.type);
        ss << ": '";
        ss << std::string(token.start, token.length);
        ss << "'";
        return ss.str();
    }
}

namespace JSON {
    Array::~Array() {
        for (auto n: array) {
            delete n;
        }
    }

    Dictionary::~Dictionary() {
        for (auto entry: dictionary) {
            delete entry.second;
        }
    }

    void
    Context::skipWhitespace() {
        while (buf[cursor] == '\t' || buf[cursor] == '\n' || buf[cursor] == '\r' || buf[cursor] == ' ') {
            cursor++;
        }
    }

    Token
    Context::nextToken() {
        skipWhitespace();

        if (buf[cursor] == '{') {
            cursor++;
            return { TokenType::DICT_OPEN, buf + cursor - 1, 1 };
        }
        if (buf[cursor] == '}') {
            cursor++;
            return { TokenType::DICT_CLOSE, buf + cursor - 1, 1 };
        }
        if (buf[cursor] == ':') {
            cursor++;
            return { TokenType::DICT_COLON, buf + cursor - 1, 1 };
        }
        if (buf[cursor] == '[') {
            cursor++;
            return { TokenType::ARR_OPEN, buf + cursor - 1, 1 };
        }
        if (buf[cursor] == ']') {
            cursor++;
            return { TokenType::ARR_CLOSE, buf + cursor - 1, 1 };
        }
        if (buf[cursor] == ',') {
            cursor++;
            return { TokenType::COMMA, buf + cursor - 1, 1 };
        }

        if (buf[cursor] == '"') {
            auto startCursor = cursor;
            cursor++;
            while (buf[cursor] != '"') {
                cursor++;
            }
            cursor++;
            return { TokenType::STRING, buf + startCursor, cursor - startCursor };
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
                return { TokenType::INTEGER, buf + startCursor, cursor - startCursor };
            }

            cursor++;

            while (buf[cursor] >= '0' && buf[cursor] <= '9') {
                cursor++;
            }

            return { TokenType::FLOAT, buf + startCursor, cursor - startCursor };
        }

        if (std::string_view(buf + cursor, 4) == "true") {
            cursor += 4;
            return { TokenType::TRUE, buf + cursor - 4, 4 };
        }

        if (std::string_view(buf + cursor, 5) == "false") {
            cursor += 5;
            return { TokenType::FALSE, buf + cursor - 5, 5 };
        }

        return { TokenType::NONE, buf, 0 };
    }

    Node *
    parseDict(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type != TokenType::DICT_OPEN) {
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        std::unordered_map<std::string, Node *> dictionary = {};

        while (true) {
            token = ctx.nextToken();
            if (token.type != TokenType::STRING) {
                for (auto entry: dictionary) {
                    delete entry.second;
                }
                ctx.cursor = cursorBefore;
                return nullptr;
            }

            auto key = std::string(token.start + 1, token.length - 2);

            token = ctx.nextToken();
            if (token.type != TokenType::DICT_COLON) {
                for (auto entry: dictionary) {
                    delete entry.second;
                }
                ctx.cursor = cursorBefore;
                return nullptr;
            }

            Node *value = parseInternal(ctx);
            if (value == nullptr) {
                for (auto entry: dictionary) {
                    delete entry.second;
                }
                ctx.cursor = cursorBefore;
                return nullptr;
            }

            dictionary[key] = value;

            token = ctx.nextToken();
            if (token.type == TokenType::COMMA) {
                continue;
            }

            if (token.type == TokenType::DICT_CLOSE) {
                break;
            }

            for (auto entry: dictionary) {
                delete entry.second;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        return new Dictionary(dictionary);
    }

    Node *
    parseArray(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type != TokenType::ARR_OPEN) {
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        std::vector<Node *> arr = {};
        while (true) {
            auto node = parseInternal(ctx);
            if (node == nullptr) {
                for (auto n: arr) {
                    delete n;
                }
                ctx.cursor = cursorBefore;
                return nullptr;
            }

            arr.emplace_back(node);

            token = ctx.nextToken();
            if (token.type == TokenType::COMMA) {
                continue;
            }

            if (token.type == TokenType::ARR_CLOSE) {
                break;
            }

            for (auto n: arr) {
                delete n;
            }
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        return new Array(arr);
    }

    Node *
    parseString(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type != TokenType::STRING) {
            ctx.cursor = cursorBefore;
            return nullptr;
        }

        auto s = std::string(token.start + 1, token.length - 2);
        return new String(s);
    }

    Node *
    parseInteger(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type != TokenType::INTEGER) {
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
        return new Integer(i);
    }

    Node *
    parseFloat(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type != TokenType::FLOAT) {
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
        return new Float(f);
    }

    Node *
    parseBool(Context &ctx) {
        auto cursorBefore = ctx.cursor;

        auto token = ctx.nextToken();
        if (token.type == TokenType::TRUE) {
            return new Bool(true);
        }

        if (token.type == TokenType::FALSE) {
            return new Bool(false);
        }

        ctx.cursor = cursorBefore;
        return nullptr;
    }

    Node *
    parseInternal(Context &ctx) {
        Node *dict = parseDict(ctx);
        if (dict != nullptr) {
            return dict;
        }

        Node *arr = parseArray(ctx);
        if (arr != nullptr) {
            return arr;
        }

        Node *str = parseString(ctx);
        if (str != nullptr) {
            return str;
        }

        Node *i = parseInteger(ctx);
        if (i != nullptr) {
            return i;
        }

        Node *f = parseFloat(ctx);
        if (f != nullptr) {
            return f;
        }

        Node *b = parseBool(ctx);
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
    traverse(Node *node, int level = 0) {
        indent(level);

        if (node->type == NodeType::DICTIONARY) {
            std::cout << "DICTIONARY: {" << std::endl;

            bool isFirst = true;
            for (auto n: ((Dictionary *) node)->dictionary) {
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

        if (node->type == NodeType::ARRAY) {
            std::cout << "ARRAY: [" << std::endl;

            for (auto n: ((Array *) node)->array) {
                traverse(n, level + 1);
            }

            indent(level);
            std::cout << "]" << std::endl;
            return;
        }

        if (node->type == NodeType::STRING) {
            std::cout << "STRING: " << ((String *) node)->s << std::endl;
            return;
        }

        if (node->type == NodeType::INTEGER) {
            std::cout << "INTEGER: " << ((Integer *) node)->i << std::endl;
            return;
        }

        if (node->type == NodeType::FLOAT) {
            std::cout << "FLOAT: " << ((Float *) node)->f << std::endl;
            return;
        }

        if (node->type == NodeType::BOOLEAN) {
            std::cout << "BOOL: " << ((Bool *) node)->b << std::endl;
            return;
        }
    }

    Node *
    Parse(char *buf, int length) {
        Context ctx = { buf, length, 0 };
        auto root = parseInternal(ctx);
        if (root == nullptr) {
            std::cerr << "failed to parse JSON" << std::endl;
            return nullptr;
        }

        return root;
    }
}
