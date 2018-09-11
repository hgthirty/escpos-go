typedef struct { const char *p; ptrdiff_t n; } _GoString;

typedef struct {
    _GoString *data;
    int len;
} GoStringArr;
