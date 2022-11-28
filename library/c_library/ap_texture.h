#ifndef AP_TEXTURE_H
#define AP_TEXTURE_H

#include "ap_utils.h"

typedef enum {
        AP_TEXTURE_TYPE_UNKNOWN = 0,
        AP_TEXTURE_TYPE_DIFFUSE = 0x1001, // texture_diffuse
        AP_TEXTURE_TYPE_SPECULAR,         // texture_specular
        AP_TEXTURE_TYPE_NORMAL,           // texture_normal
        AP_TEXTURE_TYPE_HEIGHT            // texture_height
} AP_Texture_types;

struct AP_Texture {
        // type is the AP_Texture_types value
        int type;
        // name of the image file
        char *file_name;
        // file path of the image file
        char *file_path;
        // if file_name is null, then this texture will be treated as a pure
        // RGBA color value
        float RGBA[4];
};

int ap_texture_init(struct AP_Texture *texture);

/**
 * @param type [in] the type of the texture (AP_Texture_types)
 * @param filename [in] name of the image (PNG or JPG)
 * @return struct AP_Texture* points to the texture struct, need free manually
 */
struct AP_Texture* ap_texture_generate(
        int type,
        const char *filename,
        const char *filepath
);

/**
 * @param size texture image size (square)
 * @return struct AP_Texture* points to the texture struct, need free manually
 */
struct AP_Texture* ap_texture_generate_RGBA(
        float color[4],
        int type
);

int ap_texture_set_filename(struct AP_Texture *texture, const char *name);
int ap_texture_set_filepath(struct AP_Texture *texture, const char *path);

/**
 * @brief find texture ptr from vector by file name
 *
 * @param vec texture vector
 * @param name file name of the texture image
 * @return struct AP_Texture*
 */
struct AP_Texture *ap_texture_get_ptr_by_file(
        struct AP_Vector* vec, const char *fname, const char *fpath
);

/**
 * @brief find texture ptr from vector by RGBA color value
 *
 * @param vec texture vector
 * @param color file name of the texture image
 * @return struct AP_Texture*
 */
struct AP_Texture *ap_texture_get_ptr_by_RGBA(
        struct AP_Vector *vec, float color[4]
);

#endif // AP_TEXTURE_H
