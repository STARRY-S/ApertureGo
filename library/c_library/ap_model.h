#ifndef AP_LOAD_MODEL_H
#define AP_LOAD_MODEL_H

struct AP_Model {
        struct AP_Texture *texture;
        int texture_length;
        struct AP_Mesh *mesh;
        int mesh_length;
        char *directory;
};

int ap_model_init_ptr(struct AP_Model *model, const char *path);

#endif