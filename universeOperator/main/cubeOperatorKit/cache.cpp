#include <iostream>
#include <map>
using namespace std;

struct block {
    string key;
    string value;
} cache[1024];

map<string, int> mp;

extern "C" {
void insr(const char* key1, const char* value1);
const char* ask(const char* key1);
void del(const char* key1);
}

void insr(const char* key1, const char* value1) {

    string key(key1);
    string value(value1);

    auto it = mp.find(key);
    int pos;

    if (it != mp.end() && cache[it->second].key == key) {
        pos = it->second;
        cache[pos].value = value;
    } else {
        pos = 1023;

        auto tl = mp.find(cache[pos].key);

        if (tl != mp.end()) {
            mp.erase(tl);
        }

        cache[pos].key = key;
        cache[pos].value = value;
    }

    auto hd = mp.find(cache[0].key);

    if (hd != mp.end()) {
        hd->second = pos;
    }

    swap(cache[0], cache[pos]);
    mp[key] = 0;
}

const char * ask(const char* key1) {

    string key(key1);
    auto it = mp.find(key);
    int pos;

    if (it != mp.end() && cache[it->second].key == key) {
        pos = it->second;

        auto hd = mp.find(cache[0].key);

        if (hd != mp.end()) {
            hd->second = pos;
        }

        swap(cache[0], cache[pos]);
        mp[key] = 0;

        return cache[it->second].value.c_str();
    }

    return nullptr;
}

void del(const char* key1) {
    string key(key1);
    auto it = mp.find(key);

    if (it != mp.end() && cache[it->second].key == key) {
        cache[it->second].value = nullptr;
        mp.erase(it);
    }
}
