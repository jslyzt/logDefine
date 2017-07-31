#pragma once

#include <sstream>

namespace commom {

template<typename T>
void log(std::stringstream& stream, const T& node) {
    stream << node << "|";
}

template<typename T, typename ...Args>
void log(std::stringstream& stream, const T& node, const Args& ... args) {
    log(stream, node);
    log(stream, args ...);
}

}
