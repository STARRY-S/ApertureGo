#include "ap_texture.h"
#include "ap_utils.h"
#include "ap_cvector.h"

static inline int EQUAL(float a, float b)
{
        return ((a - b) < 0.001 && (a - b) > -0.001);
}

int ap_texture_init(struct AP_Texture *texture)
{
        if (texture == NULL) {
                return AP_ERROR_INVALID_PARAMETER;
        }

        memset(texture, 0, sizeof(struct AP_Texture));
        return 0;
}

struct AP_Texture* ap_texture_generate(
        int type, const char *filename, const char* filepath)
{
        struct AP_Texture *texture = AP_MALLOC(sizeof(struct AP_Texture));
        if (texture == NULL) {
                LOGE("ap_texture_generate: malloc failed")
                return NULL;
        }
        memset(texture, 0, sizeof(struct AP_Texture));

        texture->type = type;
        ap_texture_set_filename(texture, filename);
        ap_texture_set_filepath(texture, filepath);

        return texture;
}

struct AP_Texture* ap_texture_generate_RGBA(float color[4], int type)
{
        if (color == NULL) {
                LOGE("ap_texture_generate_RGBA: invalid parameter")
                return NULL;
        }

        struct AP_Texture *texture = AP_MALLOC(sizeof(struct AP_Texture));
        if (texture == NULL) {
                LOGE("ap_texture_generate_RGBA: malloc failed")
                return NULL;
        }
        memset(texture, 0, sizeof(struct AP_Texture));

        texture->type = type;
        memcpy(texture->RGBA, color, VEC4_SIZE);

        return texture;
}

int ap_texture_set_filename(struct AP_Texture *texture, const char *filename)
{
        if (texture == NULL || filename == NULL) {
                return AP_ERROR_INVALID_POINTER;
        }

        if (texture->file_name != NULL) {
                LOGD("AP_FREE old texture file_name pointer: 0X%p", texture);
                AP_FREE(texture->file_name);
                texture->file_name = NULL;
        }
        texture->file_name = AP_MALLOC(CHAR_SIZE * (strlen(filename)+1));
        strcpy(texture->file_name, filename);
        return 0;
}

int ap_texture_set_filepath(struct AP_Texture *texture, const char *fpath)
{
        if (texture == NULL || fpath == NULL) {
                return AP_ERROR_INVALID_POINTER;
        }

        if (texture->file_path != NULL) {
                LOGD("AP_FREE old texture file_path pointer: 0X%p", texture);
                AP_FREE(texture->file_path);
                texture->file_path = NULL;
        }
        texture->file_path = AP_MALLOC(CHAR_SIZE * (strlen(fpath)+1));
        strcpy(texture->file_path, fpath);
        return 0;
}

struct AP_Texture *ap_texture_get_ptr_by_file(
        struct AP_Vector* vec, const char *fname, const char *fpath)
{
        if (fname == NULL || vec == NULL || vec->data == NULL) {
                LOGE("ap_texture_get_ptr_by_path: invalid parameter");
                return NULL;
        }

        struct AP_Texture *ptr = (struct AP_Texture*) vec->data;
        for (int i = 0; i < vec->length; ++i) {
                if (!ptr[i].file_name || !ptr[i].file_path) {
                        continue;
                }
                if (strcmp(fname, ptr[i].file_name) == 0
                        && strcmp(fpath, ptr[i].file_path) == 0)
                {
                        return ptr + i;
                }
        }
        return NULL;
}

struct AP_Texture *ap_texture_get_ptr_by_RGBA(
        struct AP_Vector *vec, float color[4])
{
        if (color == NULL || vec->data == NULL) {
                LOGE("ap_texture_get_ptr_by_RGBA: invalid parameter");
                return NULL;
        }

        struct AP_Texture *ptr = (struct AP_Texture*) vec->data;
        for (int i = 0; i < vec->length; ++i) {
                if (EQUAL(color[0], ptr[i].RGBA[0])
                   && EQUAL(color[1], ptr[i].RGBA[1])
                   && EQUAL(color[2], ptr[i].RGBA[2])
                   && EQUAL(color[3], ptr[i].RGBA[3]) )
                {
                        return ptr + i;
                }
        }
        return NULL;
}
