project(
    'aperture',
    'c',
     default_options : ['c_std=c11'],
)

cc = meson.get_compiler('c')
m_dep = cc.find_library('m', required : false)

dependencies = [
    dependency('assimp'),
    m_dep,
]

library(
    'aperture',
    'ap_custom_io.c',
    'ap_cvector.c',
    'ap_memory.c',
    'ap_mesh.c',
    'ap_model.c',
    'ap_texture.c',
    'ap_utils.c',
    'ap_vertex.c',
    dependencies : dependencies,
    version : '1.0.0',
)
