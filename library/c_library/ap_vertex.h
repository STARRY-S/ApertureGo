#ifndef AP_VERTEX_H
#define AP_VERTEX_H

#define MAX_BONE_INFLUENCE 4

struct AP_Vertex {
        float position[3];
        float normal[3];
        float tex_coords[2];

        float tangent[3];
        float big_tangent[3];

        //bone indexes which will influence this vertex
        int bonel_ids[MAX_BONE_INFLUENCE];

        //weights from each bone
        float weights[MAX_BONE_INFLUENCE];
};

#endif // AP_VERTEX_H