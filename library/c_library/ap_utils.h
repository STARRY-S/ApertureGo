#ifndef AP_UTILS_H
#define AP_UTILS_H

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#if defined(__ANDROID__)
#define AP_PLATFORM_ANDROID  1
#elif defined(_WIN32) || defined(__CYGWIN__)
#define AP_PLATFORM_WINDOWS  1
#elif defined(__linux__)
#define AP_PLATFORM_LINUX    1
#elif defined(__unix__) || defined(unix)
#define AP_PLATFORM_UNIX     1
#elif defined(__APPLE__)
#define AP_PLATFORM_MACOS    1
#endif

#define VEC2_SIZE sizeof(float[2])
#define VEC3_SIZE sizeof(float[3])
#define VEC4_SIZE sizeof(float[4])
#define MAT4_SIZE sizeof(float[16])
#define MAT3_SIZE sizeof(float[9])
#define UINT_SIZE sizeof(unsigned int)
#define  INT_SIZE sizeof(int)
#define UCHAR_SIZE sizeof(unsigned char)
#define CHAR_SIZE sizeof(char)

#define AP_COLOR_RED     "\x1b[31m"
#define AP_COLOR_GREEN   "\x1b[32m"
#define AP_COLOR_YELLOW  "\x1b[33m"
#define AP_COLOR_BLUE    "\x1b[34m"
#define AP_COLOR_MAGENTA "\x1b[35m"
#define AP_COLOR_CYAN    "\x1b[36m"
#define AP_COLOR_RESET   "\x1b[0m"

// Output log messages to stdout/stderr
#define LOGI(...) \
        fprintf(stdout, "%s[AP_MESSG] %s", AP_COLOR_GREEN, AP_COLOR_RESET); \
        fprintf(stdout, __VA_ARGS__); \
        fprintf(stdout, "\n");
#define LOGE(...) \
        fprintf(stderr, "%s[AP_ERROR] %s", AP_COLOR_RED, AP_COLOR_RESET);   \
        fprintf(stderr, __VA_ARGS__); \
        fprintf(stderr, "\n");
#define LOGW(...) \
        fprintf(stdout, "%s[AP_WARNG] %s", AP_COLOR_MAGENTA, AP_COLOR_RESET); \
        fprintf(stdout, __VA_ARGS__); \
        fprintf(stdout, "\n");

// common used error types
typedef enum {
        AP_ERROR_SUCCESS = 0,
        AP_ERROR_INVALID_POINTER,       // pointer is NULL
        AP_ERROR_INVALID_PARAMETER,     // invalid param
        AP_ERROR_MALLOC_FAILED,         // malloc failed
        AP_ERROR_MESH_UNINITIALIZED,    // mesh uninitialized
        AP_ERROR_ASSIMP_IMPORT_FAILED,  // assimp import failed
        AP_ERROR_ASSET_OPEN_FAILED,     // asset manager open file failed
        AP_ERROR_INIT_FAILED,           // initialize failed
        AP_ERROR_RENDER_FAILED,         // render failed with unknown error
        AP_ERROR_TEXTURE_FAILED,        // texture load failed unknown error
        AP_ERROR_SHADER_LOAD_FAILED,    // shader load failed
        AP_ERROR_CAMERA_NOT_SET,        // camera not set
        AP_ERROR_SHADER_NOT_SET,        // shader not set
        AP_ERROR_MODEL_NOT_SET,         // model not set
        AP_AUDIO_BUFFER_GEN_FAILED,     // audio buffer generate failed
        AP_ERROR_DECODE_FAILED,         // decode failed
        AP_ERROR_DECODE_NOT_INIT,       // decode not initialized
        AP_ERROR_DECODE_FMT_NSUPPORT,   // decode format doest not support
        AP_ERROR_UNKNOWN,               // unknown error
        AP_ERROR_LENGTH
} AP_Types;

extern const char *AP_ERROR_NAME[];

inline static void AP_CHECK(int check_i)
{
        if (check_i > 0 && check_i < AP_ERROR_LENGTH) {
                LOGE("%s", AP_ERROR_NAME[check_i]);
        }
}

#endif // AP_UTILS_H
