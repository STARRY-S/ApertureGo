#include <stdio.h>

#include <assimp/cimport.h>
#include <assimp/scene.h>
#include <assimp/postprocess.h>
#include <assimp/cfileio.h>

#include "ap_utils.h"
#include "ap_model.h"
#include "ap_custom_io.h"
#include "ap_memory.h"
#include "ap_cvector.h"
#include "ap_vertex.h"
#include "ap_texture.h"
#include "ap_mesh.h"

int ap_model_load_ptr(struct AP_Model *model, const char *path);

struct AP_Mesh *ap_model_process_mesh(struct AP_Model *model,
                        struct aiMesh *mesh,
                        const struct aiScene *scene);

int ap_model_mesh_push_back(struct AP_Model *model, struct AP_Mesh *mesh);

struct AP_Vector *ap_model_load_material_textures(
        struct AP_Model *model,
        struct aiMaterial *mat,
        enum aiTextureType type,
        int ap_type);

int ap_model_process_node(
        struct AP_Model *model,
        struct aiNode *node,
        const struct aiScene *scene);

int ap_model_texture_loaded_push_back(
        struct AP_Model *model,
        struct AP_Texture *texture);

// int ap_model_free()
// {
//         struct AP_Model *model_array = (struct AP_Model*) model_vector.data;
//         for (int i = 0; i < model_vector.length; ++i) {
//                 AP_FREE(model_array[i].directory);
//                 model_array[i].directory = NULL;
//                 AP_FREE(model_array[i].mesh);
//                 model_array[i].mesh = NULL;
//                 AP_FREE(model_array[i].texture);
//                 model_array[i].texture = NULL;
//         }
//         ap_vector_free(&model_vector);

//         return 0;
// }

int ap_model_init_ptr(struct AP_Model *model, const char *path)
{
        if (model == NULL) {
                return AP_ERROR_INVALID_POINTER;
        }

        memset(model, 0, sizeof(struct AP_Model));

        int dir_char_location = 0;
        for (int i = 0; i < strlen(path); ++i) {
                if (path[i] == '/') {
                        dir_char_location = i;
                }
        }
        if (dir_char_location >= 0) {
                char *dir_path =
                        AP_MALLOC(sizeof(char) * (dir_char_location + 2));
                memcpy(dir_path, path, (dir_char_location + 1) * sizeof(char));
                dir_path[dir_char_location + 1] = '\0';
                model->directory = dir_path;
        } else {
                char *dir_path = AP_MALLOC(sizeof(char));
                dir_path[0] = '\0';
                model->directory = dir_path;
        }

        return ap_model_load_ptr(model, path);
}

int ap_model_load_ptr(struct AP_Model *model, const char *path)
{
        if (model == NULL) {
                return AP_ERROR_INVALID_POINTER;
        }

        // custom file io for assimp
        struct aiFileIO fileIo;
        fileIo.CloseProc = ap_custom_file_close_proc;
        fileIo.OpenProc = ap_custom_file_open_proc;
        fileIo.UserData = NULL;

        const struct aiScene* scene = aiImportFileEx(
                path,
                aiProcess_Triangulate | aiProcess_GenSmoothNormals |
                aiProcess_FlipUVs | aiProcess_CalcTangentSpace,
                &fileIo
        );

        if(!scene || scene->mFlags & AI_SCENE_FLAGS_INCOMPLETE
                || !scene->mRootNode)
        {
                LOGE("aiImportFileEx failed: \n%s", aiGetErrorString());
                return AP_ERROR_ASSIMP_IMPORT_FAILED;
        }

        ap_model_process_node(model, scene->mRootNode, scene);

        aiReleaseImport(scene);

        return AP_ERROR_SUCCESS;
}

int ap_model_process_node(
        struct AP_Model *model,
        struct aiNode *node,
        const struct aiScene *scene)
{
        if (model == NULL || node == NULL || scene == NULL) {
                return AP_ERROR_INVALID_POINTER;
        }
        // process each mesh located in the current node
        for (unsigned int i = 0; i < node->mNumMeshes; i++)
        {
                // the node object only contains indices to index the actual
                // objects in the scene. the scene contains all the data,
                // node is just to keep stuff organized
                // (like relations between nodes).
                struct aiMesh* mesh = scene->mMeshes[node->mMeshes[i]];
                struct AP_Mesh *mesh_new =
                        ap_model_process_mesh(model, mesh, scene);
                ap_model_mesh_push_back(model, mesh_new);
                ap_mesh_free(mesh_new);
        }

        // after we've processed all of the meshes (if any)
        // we then recursively process each of the children nodes
        for (unsigned int i = 0; i < node->mNumChildren; i++)
        {
                ap_model_process_node(model, node->mChildren[i], scene);
        }

        return 0;
}

struct AP_Mesh *ap_model_process_mesh(struct AP_Model *model,
                        struct aiMesh *mesh,
                        const struct aiScene *scene)
{
        if (model == NULL || mesh == NULL || scene == NULL) {
                return NULL;
        }
        // create buffers
        static struct AP_Vector vec_vertices;
        static struct AP_Vector vec_indices;
        static struct AP_Vector vec_textures;

        ap_vector_init(&vec_vertices, AP_VECTOR_VERTEX);
        ap_vector_init(&vec_indices, AP_VECTOR_UINT);
        ap_vector_init(&vec_textures, AP_VECTOR_TEXTURE);

        // walk through each of the mesh's vertices
        for (unsigned int i = 0; i < mesh->mNumVertices; i++)
        {
                struct AP_Vertex vertex;
                // positions
                vertex.position[0] = mesh->mVertices[i].x;
                vertex.position[1] = mesh->mVertices[i].y;
                vertex.position[2] = mesh->mVertices[i].z;

                // normals
                if (mesh->mNormals != NULL) {
                        vertex.normal[0] = mesh->mNormals[i].x;
                        vertex.normal[1] = mesh->mNormals[i].y;
                        vertex.normal[2] = mesh->mNormals[i].z;
                }
                // texture coordinates
                // does the mesh contain texture coordinates?
                if(mesh->mTextureCoords[0]) {
                        /**
                         * a vertex can contain up to 8 different texture
                         * coordinates. We thus make the assumption that we
                         * won't use models where a vertex can have multiple
                         * texture coordinates so we always take the
                         * first set (0).
                         */
                        vertex.tex_coords[0] = mesh->mTextureCoords[0][i].x;
                        vertex.tex_coords[1] = mesh->mTextureCoords[0][i].y;

                        // tangent
                        vertex.tangent[0] = mesh->mTangents[i].x;
                        vertex.tangent[1] = mesh->mTangents[i].y;
                        vertex.tangent[2] = mesh->mTangents[i].z;

                        // big tangent
                        vertex.big_tangent[0] = mesh->mBitangents[i].x;
                        vertex.big_tangent[1] = mesh->mBitangents[i].y;
                        vertex.big_tangent[2] = mesh->mBitangents[i].z;
                } else {
                        vertex.tex_coords[0] = 0.0f;
                        vertex.tex_coords[1] = 0.0f;
                }

                // vertices.push_back(vertex);
                ap_vector_push_back(&vec_vertices, (const char *) &vertex);
        }
        // now wak through each of the mesh's faces
        // (a face is a mesh its triangle) and retrieve the corresponding
        // vertex indices.
        for (unsigned int i = 0; i < mesh->mNumFaces; i++)
        {
                struct aiFace face = mesh->mFaces[i];
                // retrieve all indices of the face and store them in the
                // indices vector
                for(unsigned int j = 0; j < face.mNumIndices; j++) {
                        ap_vector_push_back(
                                &vec_indices,
                                (const char *) &face.mIndices[j]
                        );
                }
        }
        // process materials
        struct aiMaterial* material = scene->mMaterials[mesh->mMaterialIndex];
        /**
         * we assume a convention for sampler names in the shaders.
         * Each diffuse texture should be named as 'texture_diffuseN' where
         * N is a sequential number ranging from 1 to MAX_SAMPLER_NUMBER.
         *
         * Same applies to other texture as the following list summarizes:
         * diffuse: texture_diffuseN
         * specular: texture_specularN
         * normal: texture_normalN
         */

        // 1. diffuse maps
        size_t size = sizeof(struct AP_Texture);
        struct AP_Vector *vec_diffuse_ptr = NULL;
        vec_diffuse_ptr = ap_model_load_material_textures(
                model, material,
                aiTextureType_DIFFUSE, AP_TEXTURE_TYPE_DIFFUSE
        );
        if (vec_diffuse_ptr->length > 0) {
                ap_vector_insert_back(
                        &vec_textures,
                        vec_diffuse_ptr->data,
                        size * vec_diffuse_ptr->length
                );
        }
        ap_vector_free(vec_diffuse_ptr);

        // 2. specular maps
        struct AP_Vector *vec_specular_ptr = NULL;
        vec_specular_ptr = ap_model_load_material_textures(
                model, material,
                aiTextureType_SPECULAR, AP_TEXTURE_TYPE_SPECULAR
        );
        if (vec_specular_ptr->length > 0) {
                ap_vector_insert_back(
                        &vec_textures,
                        vec_specular_ptr->data,
                        size * vec_specular_ptr->length
                );
        }
        ap_vector_free(vec_specular_ptr);

        // 3. normal maps
        struct AP_Vector *vec_normal_ptr = NULL;
        vec_normal_ptr = ap_model_load_material_textures(
                model, material,
                aiTextureType_NORMALS, AP_TEXTURE_TYPE_NORMAL
        );
        if (vec_normal_ptr->length > 0) {
                ap_vector_insert_back(
                        &vec_textures,
                        vec_normal_ptr->data,
                        size * vec_normal_ptr->length
                );
        }
        ap_vector_free(vec_normal_ptr);

        // 4. height maps
        struct AP_Vector *vec_height_ptr;
        vec_height_ptr = ap_model_load_material_textures(
                model, material,
                aiTextureType_HEIGHT, AP_TEXTURE_TYPE_HEIGHT
        );
        if (vec_height_ptr->length > 0) {
                ap_vector_insert_back(
                        &vec_textures,
                        vec_height_ptr->data,
                        size * vec_height_ptr->length
                );
        }
        ap_vector_free(vec_height_ptr);

        // return a mesh object created from the extracted mesh data
        static struct AP_Mesh mesh_buffer;
        ap_mesh_init(
                &mesh_buffer,
                (struct AP_Vertex *) vec_vertices.data,
                vec_vertices.length,
                (unsigned int *) vec_indices.data,
                vec_indices.length,
                (struct AP_Texture *) vec_textures.data,
                vec_textures.length
        );

        // release buffers
        ap_vector_free(&vec_vertices);
        ap_vector_free(&vec_indices);
        ap_vector_free(&vec_textures);

        return &mesh_buffer;
}

struct AP_Vector *ap_model_load_material_textures(
        struct AP_Model *model,
        struct aiMaterial *mat,
        enum aiTextureType type,
        int ap_type)
{
        // create a new buffer for store new textures
        static struct AP_Vector vec_texture;
        // this vector stores the textures already loaded in model
        struct AP_Vector vec_model_tex = {
                .length = model->texture_length,
                .capacity = model->texture_length,
                .data = (char*) model->texture,
                .type = AP_VECTOR_TEXTURE
        };
        ap_vector_init(&vec_texture, AP_VECTOR_TEXTURE);

        uint32_t mat_texture_count = aiGetMaterialTextureCount(mat, type);
        for (uint32_t i = 0; i < mat_texture_count; i++)
        {
                struct aiString str;
                aiGetMaterialTexture(
                        mat, type, i, &str,
                        NULL, NULL, NULL, NULL, NULL, NULL
                );
                struct AP_Texture *ptr = NULL;
                ptr = ap_texture_get_ptr_by_file(
                        &vec_model_tex, str.data, model->directory);
                if (ptr == NULL) {
                        ptr = ap_texture_generate(
                                ap_type, str.data, model->directory);
                        if (ptr != NULL) {
                                ap_vector_push_back(&vec_texture, (char*) ptr);
                                ap_model_texture_loaded_push_back(model, ptr);
                                AP_FREE(ptr);
                                ptr = NULL;
                        }
                }
        }

        struct aiColor4D ai_color = { 0.f, 0.f, 0.f, 0.f };
        if (ap_type == AP_TEXTURE_TYPE_DIFFUSE) {
                aiGetMaterialColor(mat, AI_MATKEY_COLOR_DIFFUSE, &ai_color);
        } else if (ap_type == AP_TEXTURE_TYPE_SPECULAR) {
                aiGetMaterialColor(mat, AI_MATKEY_COLOR_SPECULAR, &ai_color);
        }

        if (!(ai_color.r || ai_color.g || ai_color.b)) {
                return &vec_texture;
        }

        float color[4] = { 0.0f };
        memcpy(color, &ai_color, sizeof(float) * 4);
        struct AP_Texture *ptr = NULL;
        if (vec_model_tex.data != NULL) {
                ptr = ap_texture_get_ptr_by_RGBA(&vec_model_tex, color);
        }
        if (ptr == NULL) {
                ptr = ap_texture_generate_RGBA(color, AP_TEXTURE_TYPE_DIFFUSE);
                if (ptr != NULL) {
                        ap_vector_push_back(&vec_texture, (char*) ptr);
                        ap_model_texture_loaded_push_back(model, ptr);
                        AP_FREE(ptr);
                        ptr = NULL;
                } else {
                        LOGE("ap_model_load_material_textures: ap_texture_generate_RGBA failed")
                }
        }

        return &vec_texture;
}

int ap_model_texture_loaded_push_back(
        struct AP_Model *model,
        struct AP_Texture *texture)
{
        if (model == NULL || texture == NULL) {
                return AP_ERROR_INVALID_PARAMETER;
        }

        // add a new texture struct object into model
        model->texture = AP_REALLOC(model->texture,
                sizeof(struct AP_Texture) * (model->texture_length + 1));
        if (model->texture == NULL) {
                LOGE("realloc error");
                return AP_ERROR_MALLOC_FAILED;
        }
        struct AP_Texture *texture_new =
                model->texture + (model->texture_length);
        model->texture_length++;
        ap_texture_init(texture_new);
        texture_new->type = texture->type;
        ap_texture_set_filename(texture_new, texture->file_name);
        ap_texture_set_filepath(texture_new, texture->file_path);
        // texture_new->id = texture->id;
        memcpy(texture_new->RGBA, texture->RGBA, VEC4_SIZE);

        return 0;
}

int ap_model_mesh_push_back(struct AP_Model *model, struct AP_Mesh *mesh)
{
        if (model == NULL || mesh == NULL) {
                LOGE("ap_model_mesh_push_back: invalid param");
                return AP_ERROR_INVALID_POINTER;
        }

        // add a new mesh struct object into model
        model->mesh = AP_REALLOC(
                model->mesh,
                sizeof(struct AP_Mesh) * (model->mesh_length + 1)
        );
        if (model->mesh == NULL) {
                LOGE("ap_model_mesh_push_back: realloc error");
                return AP_ERROR_MALLOC_FAILED;
        }
        struct AP_Mesh *mesh_ptr = &model->mesh[model->mesh_length];
        model->mesh_length++;
        // copy old mesh data memory to new mesh
        ap_mesh_copy(mesh_ptr, mesh);

        return 0;
}