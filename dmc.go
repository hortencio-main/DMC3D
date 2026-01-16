// This software is released into the public domain under the Unlicense.
//
// Anyone is free to copy, modify, publish, use, compile, sell, or
// distribute this software, either in source code form or as a compiled
// binary, for any purpose, commercial or non-commercial, and by any
// means.
//
// In jurisdictions that do not recognize the public domain, the author
// grants a perpetual, irrevocable, royalty-free license to exercise all
// rights associated with this software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package main

import (
    "compress/flate"
    "encoding/binary"
    "fmt"
    "sort"
    "strings"
    "image"
    "image/draw"
    "image/png"
    "os"
    "log"
    "time"
    "math"
    "math/rand"
	"runtime"
    "io"
    //~ "reflect"
	
    //~ graphics
	//~ "github.com/go-gl/gl/v4.6-compatibility/gl"
	"github.com/go-gl/gl/v3.3-compatibility/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
    //~ "github.com/go-gl/gl/v2.1/gl"

    //~ audio
    "github.com/gen2brain/malgo"
	"github.com/youpy/go-wav"
    
)

var width, height = 800, 600
//~ var widthInternal int = 639
//~ var heightInternal int = 479
var widthInternal int = 639
var heightInternal int = 479

const DEBUG = false
const (
    PI = 3.14159
    SUBCHUNK_H = 16
    RENDER_DISTANCE = 8
    //~ CHUNK_V = 120
    //~ CHUNK_V = 120
    CHUNK_V = 255
)

type BlockID = uint16
const (
    AIR = iota
    STONE
    GRASS
    DIRT
    COBBLE
    PLANK
    BEDROCK
    WATER
    LAVA
    SAND
    GRAVEL
    ORE_IRON
    ORE_COAL
    WOOD
    LEAF
    CRAFTING_TABLE
    CLOUD
    STICK_PILE
    GLOW
    BREAKING
    BRICK
    PAINTING
    SLAB
    FURNACE
    ICE
    SNOW
    DUCK
    BEE
    MOON
    CLAY
    FLESH
    OBSCURE
    CHEST
    
    ERROR
    RAYTRACED_SPHERE
)

const (
    ITEM_0 = 1024 + iota
    PICKAXE
    STICK
    SAPLING
    CANDLE
    FLOWER
    PIXIE
    GRASS_TUFF
    LAVENDER
    AXE
    SHOVEL
    SWORD
    INGOT
    LILY
    WAX
    SIGN
    N_BLOCK_IDS
)


var BlockIsLiquid = [N_BLOCK_IDS]bool{ 
    WATER: true, 
    LAVA: true,
}


const (
    FONT_0 = 6
    FONT_A = 6 + 10
)


type Recipe struct {
    from []BlockID
    fromc []int32
    to []BlockID
    toc []int32
}

var crafting_table_recipes = []Recipe{
    {
        from:  []BlockID{STICK, COBBLE},
        fromc: []int32{1,2},
        to:    []BlockID{SWORD},
        toc:   []int32{64},
    },
    {
        from:  []BlockID{STICK, COBBLE},
        fromc: []int32{2,1},
        to:    []BlockID{SHOVEL},
        toc:   []int32{64},
    },
    {
        from:  []BlockID{STICK, COBBLE},
        fromc: []int32{2,2},
        to:    []BlockID{AXE},
        toc:   []int32{64},
    },
    {
        from:  []BlockID{STICK, COBBLE},
        fromc: []int32{2,3},
        to:    []BlockID{PICKAXE},
        toc:   []int32{64},
    },
    {
        from:  []BlockID{STICK, PLANK},
        fromc: []int32{1,2},
        to:    []BlockID{SWORD},
        toc:   []int32{32},
    },
    {
        from:  []BlockID{STICK, PLANK},
        fromc: []int32{2,1},
        to:    []BlockID{SHOVEL},
        toc:   []int32{32},
    },
    {
        from:  []BlockID{STICK, PLANK},
        fromc: []int32{2,2},
        to:    []BlockID{AXE},
        toc:   []int32{32},
    },
    {
        from:  []BlockID{STICK, PLANK},
        fromc: []int32{2,3},
        to:    []BlockID{PICKAXE},
        toc:   []int32{32},
    },

    {
        from:  []BlockID{WOOD},
        fromc: []int32{1},
        to:    []BlockID{PLANK},
        toc:   []int32{4},
    },
    {
        from:  []BlockID{COBBLE},
        fromc: []int32{8},
        to:    []BlockID{FURNACE},
        toc:   []int32{1},
    },
    {
        from:  []BlockID{FLESH},
        fromc: []int32{1},
        to:    []BlockID{PLANK},
        toc:   []int32{4},
    },
    {
        from:  []BlockID{PLANK},
        fromc: []int32{4},
        to:    []BlockID{CRAFTING_TABLE},
        toc:   []int32{1},
    },
    {
        from:  []BlockID{PLANK},
        fromc: []int32{2},
        to:    []BlockID{STICK},
        toc:   []int32{4},
    },
    {
        from:  []BlockID{COBBLE},
        fromc: []int32{8},
        to:    []BlockID{FURNACE},
        toc:   []int32{1},
    },
    {
        from:  []BlockID{STICK, ORE_COAL},
        fromc: []int32{1,1},
        to:    []BlockID{CANDLE},
        toc:   []int32{1},
    },
    {
        from:  []BlockID{WAX},
        fromc: []int32{4},
        to:    []BlockID{CANDLE},
        toc:   []int32{1},
    },
}


var furnace_recipes = []Recipe{
    {
        from:  []BlockID{ORE_IRON, ORE_COAL},
        fromc: []int32{1, 1},
        to:    []BlockID{INGOT},
        toc:   []int32{4},
    },
    {
        from:  []BlockID{COBBLE, ORE_COAL},
        fromc: []int32{1, 1},
        to:    []BlockID{STONE},
        toc:   []int32{1},
    },
    { /* note: not make another blockid just for charcoal */
        from:  []BlockID{WOOD},
        fromc: []int32{2},
        to:    []BlockID{ORE_COAL},
        toc:   []int32{1},
    },
}
    
    
type Crafting struct {
    prog float32
    dir Vec3
    pos Vec3
    id BlockID
    c int32
}

var block_crafting []Crafting

const (
    LIQUID_FLOWING_XP = iota
    LIQUID_FLOWING_XN
    LIQUID_FLOWING_ZP
    LIQUID_FLOWING_ZN
    LIQUID_FLOWING_DOWN
)

type Attributes uint8


const lightLevels = 8
const DAY_LEN = 24000
    
const (
    WATER_XP = iota
    WATER_XN
    WATER_ZP
    WATER_ZN
)
type Chunk struct {
    block             [SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]BlockID
    light             [SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]uint8
    cloud             [SUBCHUNK_H][SUBCHUNK_H]byte
    
    display_list      uint32
    LOD1              uint32
    
    uniform_index int
    uniform_array     [256*3]float32

    scheduledForLightUpdate bool

    blockInventories map[nearVec]Inventory
    blockText map[nearVec][]byte
    
    blockAttributes map[nearVec][]Attributes

    water_flow map[nearVec]struct{h, d uint8}
}


type nearVec struct {
    x,y,z uint8
}

type Vec3 struct{
    x,y,z float32
}

type iVec2 struct{
    x,y int
}

type Inventory struct {
    ID []BlockID
    C  []int32 /* count */
    D  []int32 /* durability */
    T  int /* time since last insertion */
    S int  /* size in slots */
}

const (
    MOBID_DUCK = iota
    MOBID_BEE
    PIG
    SHEEP
    COW
    N_MOB_IDS
)

const (
    MOB_IDLE = iota
    MOB_TURNING
    MOB_JUMPING
    MOB_WALKING
    N_MOB_STATES
)

type EntityID uint8
type EntityState uint8
type Entity struct {
    hp float32
    pos Vec3
    dir Vec3
    id EntityID
    state EntityState
}
var Entities []Entity

type Drop struct {
    pos Vec3
    dir Vec3
    rot float32
    id BlockID
    c int32
}

var (
    player_pos Vec3
    player_vel Vec3
    player_inventory Inventory
	cam_yaw   float32
	cam_pitch float32
    atlas uint32 
    Background uint32
    
    
    TextureSun uint32
    
    item_atlas uint32
	first_mouse bool  = true
	last_mouse_x float64
	last_mouse_y float64
    last_mouse_z float64
    fov = 60.0
    
    // use this when adding a new dimension
    //~ Dimensions []map[iVec2]*Chunk
    World map[iVec2]*Chunk
    blockUpdates [][3]int
    itemDrops []Drop

    breakingBlocks []struct{x, y, z, t, id, elapsed int}

    //~ TICK uint32 = DAY_LEN/2
    TICK uint32 = DAY_LEN/4
    //~ TICK uint32 = 0
    sceneTex uint32
    blurTex uint32
    openInventory bool = false
    openInventoryX int
    openInventoryY int
    openInventoryZ int


    atlasFrames [4]uint32
)




func init() {
	runtime.LockOSThread()
    World = make(map[iVec2]*Chunk)
}

func bringInventory (x, y, z int) {
    openInventory = true
    openInventoryX = x
    openInventoryY = y
    openInventoryZ = z
}
func closeInventory () {
    openInventory = false
}

func getInventory(x,y,z int) Inventory {
    nv := nearVec{uint8(x%SUBCHUNK_H),uint8(y),uint8(z%SUBCHUNK_H)}
    return World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].blockInventories[nv]
}

func getWaterflow(x,y,z int) (h uint8) {
    nv := nearVec{uint8(x%SUBCHUNK_H),uint8(y),uint8(z%SUBCHUNK_H)}
    
    var a struct{h, d uint8} = World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].water_flow[nv]
    
    return a.h
}
func setWaterflow(x,y,z int, h, d uint8) {
    c := World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}]
    if c.water_flow == nil {
        c.water_flow = make(map[nearVec]struct{h, d uint8})
    }
    nv := nearVec{uint8(x%SUBCHUNK_H),uint8(y),uint8(z%SUBCHUNK_H)}
    c.water_flow[nv] = struct{h, d uint8} {h, d}
}

// ProjectWorldToScreen converts a 3D world position to screen coordinates
func ProjectWorldToScreen(x, y, z float32) (sx, sy, sz float32) {
    var modelview [16]float32
    var projection [16]float32
    var viewport [4]int32

    // Get current matrices and viewport
    gl.GetFloatv(gl.MODELVIEW_MATRIX, &modelview[0])
    gl.GetFloatv(gl.PROJECTION_MATRIX, &projection[0])
    gl.GetIntegerv(gl.VIEWPORT, &viewport[0])

    // Multiply point by ModelView
    mvx := modelview[0]*x + modelview[4]*y + modelview[8]*z + modelview[12]
    mvy := modelview[1]*x + modelview[5]*y + modelview[9]*z + modelview[13]
    mvz := modelview[2]*x + modelview[6]*y + modelview[10]*z + modelview[14]
    mvw := modelview[3]*x + modelview[7]*y + modelview[11]*z + modelview[15]

    // Multiply by Projection
    px := projection[0]*mvx + projection[4]*mvy + projection[8]*mvz + projection[12]*mvw
    py := projection[1]*mvx + projection[5]*mvy + projection[9]*mvz + projection[13]*mvw
    pz := projection[2]*mvx + projection[6]*mvy + projection[10]*mvz + projection[14]*mvw
    pw := projection[3]*mvx + projection[7]*mvy + projection[11]*mvz + projection[15]*mvw

    if pw == 0.0 {
        return -1, -1, -1 // point is invalid
    }

    // Perspective divide (Normalized Device Coordinates)
    ndcX := px / pw
    ndcY := py / pw
    ndcZ := pz / pw

    // Map to window coordinates
    sx = (ndcX*0.5 + 0.5) * float32(viewport[2]) + float32(viewport[0])
    sy = (ndcY*0.5 + 0.5) * float32(viewport[3]) + float32(viewport[1])
    sz = (ndcZ*0.5 + 0.5)

    return
}



func init() {

}


var RaytracingProgram uint32
var RaytracingUniformuSphereCenter uint32

// RaytracingSourceEnum 
const (
    RaytracingSource_sphere = float32(0.0)
    RaytracingSource_octa   = float32(1.0)
    RaytracingSource_donut  = float32(2.0)
)

var RaytracingSourceV = `
#version 120
varying vec2 vClip;
void main()
{
    gl_Position = gl_Vertex;
    vClip = gl_Vertex.xy; // clip-space quad
}

`//~ #version 120
//~ void main()
//~ {
    //~ gl_Position = gl_Vertex;
//~ }



var RaytracingSourceF = `
#version 120

#define MAX_SPHERES 256

uniform mat4 uInvProjection;
uniform mat4 uInvView;

uniform float Seed;

uniform vec4 uSphereCenter[MAX_SPHERES];
const float uSphereRadius = 0.5;

vec3 uLightDir = vec3(1,1,1);

varying vec2 vClip;


// Torus
#define MAX_STEPS 32
#define HIT_EPS 0.001
#define FAR_DIST 50.0
const float TORUS_R = 0.6;  // major radius
const float TORUS_r = 0.2;  // tube radius
mat3 rotY(float a)
{
    float c = cos(a);
    float s = sin(a);
    return mat3(
         c, 0.0, -s,
        0.0, 1.0, 0.0,
         s, 0.0,  c
    );
}
float sdTorus(vec3 p)
{
    p = rotY(0.8) * p;
    vec2 q = vec2(length(p.xz) - TORUS_R, p.y);
    return length(q) - TORUS_r;
}
vec3 calcTorusNormal(vec3 p)
{
    const vec2 e = vec2(0.001, 0.0);
    return normalize(vec3(
        sdTorus(p + e.xyy) - sdTorus(p - e.xyy),
        sdTorus(p + e.yxy) - sdTorus(p - e.yxy),
        sdTorus(p + e.yyx) - sdTorus(p - e.yyx)
    ));
}



// Octahedron
const float uOctRadius = 0.5;

bool intersectOctahedron(
    vec3 ro, vec3 rd,
    vec3 center,
    out float tHit,
    out vec3 hitNormal)
{
    vec3 o = ro - center;
    float tMin = -1e20;
    float tMax =  1e20;
    vec3 nMin = vec3(0.0);

    // Planes: ±x ±y ±z = R
    for (int sx = -1; sx <= 1; sx += 2)
    for (int sy = -1; sy <= 1; sy += 2)
    for (int sz = -1; sz <= 1; sz += 2)
    {
        vec3 n = vec3(float(sx), float(sy), float(sz));
        float denom = dot(n, rd);
        float dist  = uOctRadius - dot(n, o);

        if (abs(denom) < 1e-6)
        {
            if (dist < 0.0)
                return false;
        }
        else
        {
            float t = dist / denom;
            if (denom < 0.0)
            {
                if (t > tMin)
                {
                    tMin = t;
                    nMin = n;
                }
            }
            else
            {
                tMax = min(tMax, t);
            }
            if (tMin > tMax)
                return false;
        }
    }

    if (tMin < 0.0)
        return false;

    tHit = tMin;
    hitNormal = normalize(nMin);
    return true;
}


void main()
{
    // Noise subsampling
    //~ float n = fract(sin(dot(gl_FragCoord.xy, vec2((1+Seed)/12.9898,78.233))) * 43758.5453);
    //~ if (n > 0.1)
        //~ discard;

    // Reconstruct ray
    vec4 pNear = uInvProjection * vec4(vClip, -1.0, 1.0);
    pNear /= pNear.w;

    vec4 pFar = uInvProjection * vec4(vClip, 1.0, 1.0);
    pFar /= pFar.w;

    vec3 rayDirView = normalize(pFar.xyz - pNear.xyz);

    vec3 rayOrigin = (uInvView * vec4(pNear.xyz, 1.0)).xyz;
    vec3 rayDir    = normalize((uInvView * vec4(rayDirView, 0.0)).xyz);

    float tMin = 1e20;
    bool hit = false;
    vec3 hitCenter = vec3(0.0);

    for (int i = 0; i < MAX_SPHERES; i++)
    {
        vec3 center = uSphereCenter[i].xyz;
        if (center.x == 0.0)
            continue;
        if (uSphereCenter[i].w ==` + fmt.Sprintf("%f",RaytracingSource_sphere) + `) {
            vec3 oc = rayOrigin - center;

            float b = dot(oc, rayDir);
            float c = dot(oc, oc) - uSphereRadius * uSphereRadius;
            float h = b*b - c;

            if (h < 0.0)
                continue;

            float t = -b - sqrt(h);
            if (t > 0.0 && t < tMin)
            {
                tMin = t;
                hitCenter = center;
                hit = true;
            }
        } else if (uSphereCenter[i].w ==` + fmt.Sprintf("%f",RaytracingSource_donut) + `) {
            float t = 0.0;
            bool localHit = false;

            for (int step = 0; step < MAX_STEPS; step++)
            {
                vec3 p = rayOrigin + rayDir * t;
                float d = sdTorus(p - center);

                if (d < HIT_EPS)
                {
                    localHit = true;
                    break;
                }

                t += d;
                if (t > FAR_DIST)
                    break;
            }

            if (localHit && t > 0.0 && t < tMin)
            {
                tMin = t;
                hitCenter = center;
                hit = true;
            }
        } else if (uSphereCenter[i].w == ` + fmt.Sprintf("%f", RaytracingSource_octa) + `) {
            float t;
            vec3 n;
            if (intersectOctahedron(rayOrigin, rayDir, center, t, n))
            {
                if (t > 0.0 && t < tMin)
                {
                    tMin = t;
                    //~ hitNormal = n;
                    hitCenter = center;
                    hit = true;
                }
            }
        }
        

    }

    if (!hit)
        discard;

    vec3 hitPos = rayOrigin + rayDir * tMin;
    vec4 clipPos = gl_ProjectionMatrix * gl_ModelViewMatrix * vec4(hitPos, 1.0);

    float ndcDepth = clipPos.z / clipPos.w;
    gl_FragDepth = ndcDepth * 0.5 + 0.5;
    
    
    vec3 normal = normalize(hitPos - hitCenter);
    float light = max(dot(normal, normalize(uLightDir)), 0.0);

    // quantization
    //~ light = floor(light * 8.0) / 8.0;

    vec3 color = vec3(1.0, 1.0, 1.0) * light;

    gl_FragColor = vec4(color, 1.0);
}
`


var subpixelSnapProgram uint32
var subpixelSnapUniformTextured int32
var subpixelSnapUniformSeed int32


var minimumFragmentSource = `
#version 130

varying vec3 v_pos;

uniform sampler2D tex;
uniform bool uBillboard;

uniform int uSeed;

//~ noperspective in vec2 v_uv;

float hash(vec3 p)
{
    return fract(sin(dot(p, vec3(12.9898,78.233,45.164))) * 43758.5453);
}


float valueNoise(vec3 p)
{
    vec3 i = floor(p);
    vec3 f = fract(p);
    f = f * f * (3.0 - 2.0 * f); // smoothstep
    float n000 = hash(i + vec3(0,0,0));
    float n100 = hash(i + vec3(1,0,0));
    float n010 = hash(i + vec3(0,1,0));
    float n110 = hash(i + vec3(1,1,0));
    float n001 = hash(i + vec3(0,0,1));
    float n101 = hash(i + vec3(1,0,1));
    float n011 = hash(i + vec3(0,1,1));
    float n111 = hash(i + vec3(1,1,1));
    return mix(
        mix(mix(n000, n100, f.x), mix(n010, n110, f.x), f.y),
        mix(mix(n001, n101, f.x), mix(n011, n111, f.x), f.y),
        f.z
    );
}

void main()
{
    // Control scale of fog noise
    float noiseScale = 0.08;
    float vn = valueNoise(v_pos * noiseScale);
    float fogNoise = vn;
    
    fogNoise = mix(1.0, 0.75, fogNoise);
    
    float fogFactor = ((gl_Fog.end * fogNoise) - gl_FogFragCoord) * gl_Fog.scale;

    // still figuring out how to make water look more foggy
    fogFactor = clamp(fogFactor, 0.0, 1.0);

    vec4 color;
    
    if (uBillboard)
        color = texture2D(tex, gl_TexCoord[0].st) * gl_Color;
        //~ color = texture2D(tex, v_uv) * gl_Color;
    else {
        color = gl_Color;
    }
    
    vec3 v_pos_l = v_pos;
    
    v_pos_l.y += 5*sin(float(uSeed)/30.0);
    
    color.r = mix(  mix(0, color.r, valueNoise(v_pos_l * 0.2)), color.r, gl_Color.w);
    color.g = mix(  mix(0, color.g, valueNoise(v_pos_l * 0.2)), color.g, gl_Color.w);
    color.b = mix(  mix(0, color.b, valueNoise(v_pos_l * 0.2)), color.b, gl_Color.w);
    

    //~ color.r = mix(1.0, color.r, gl_Color.w);

    // temp fix, I want items to be transprent but not the sea bottom
    if (gl_Color.w < 0.99f) {
        color.w = 1.0;
    }
    gl_FragColor = mix(gl_Fog.color, color, fogFactor);
}
`


var subpixelSnapSource = `
#version 130

varying vec3 v_pos;
//~ noperspective out vec2 v_uv;

uniform vec2 uViewport;
uniform bool textured;
uniform int uSeed;

const float SNAP = 2.0;

float hash(vec3 p) {
    return fract(sin(dot(p, vec3(12.9898,78.233,45.164))) * 43758.5453);
}

float valueNoise(vec3 p) {
    vec3 i = floor(p);
    vec3 f = fract(p);
    f = f * f * (3.0 - 2.0 * f); // smoothstep
    float n000 = hash(i + vec3(0,0,0));
    float n100 = hash(i + vec3(1,0,0));
    float n010 = hash(i + vec3(0,1,0));
    float n110 = hash(i + vec3(1,1,0));
    float n001 = hash(i + vec3(0,0,1));
    float n101 = hash(i + vec3(1,0,1));
    float n011 = hash(i + vec3(0,1,1));
    float n111 = hash(i + vec3(1,1,1));
    return mix(
        mix(mix(n000, n100, f.x), mix(n010, n110, f.x), f.y),
        mix(mix(n001, n101, f.x), mix(n011, n111, f.x), f.y),
        f.z
    );
}


void main() {

    // eye-space position for volumetric fog
    vec4 eyePos = gl_Vertex;
    v_pos = eyePos.xyz;
    
    
    // Fog Z
    float dist = abs((gl_ModelViewMatrix * gl_Vertex).z);
    
    vec4 clip = gl_ModelViewProjectionMatrix * gl_Vertex;

    vec3 v_pos_l = v_pos;    
    v_pos_l.y += 3*sin(float(uSeed)/30.0);
    clip.y = mix(clip.y+valueNoise(v_pos_l * 0.2), clip.y, gl_Color.w);
    
    
    float invW = 1.0 / clip.w;
    vec2 ndc = clip.xy * invW;

    vec2 screen = (ndc * 0.5 + 0.5) * uViewport;

    vec2 center = uViewport * 0.5;
    vec2 d = screen - center;
    vec2 s = d / SNAP;

    // branchless snap away
    screen =
        center +
        SNAP *
        mix(floor(s), ceil(s), step(0.0, d));

    ndc = (screen / uViewport) * 2.0 - 1.0;

    // add some noise to the depth fog
    gl_FogFragCoord = dist; 
    
    gl_Position = vec4(ndc * clip.w, clip.z, clip.w);

    gl_FrontColor = gl_Color;
    gl_TexCoord[0] = gl_MultiTexCoord0;
    //~ v_uv = gl_MultiTexCoord0.st;

    // Projection Y scale (vertical FOV)
    float projScaleY = gl_ProjectionMatrix[1][1];

    // Base size of a 1x1 camera-facing quad in pixels
    float sizePixels = (projScaleY / clip.w) * uViewport.y * 0.5;

    // Worst-case expansion
    //~ sizePixels *= 1.41421356; // sqrt(2)
    sizePixels *= 1.5; // sqrt(2)

    // Snap UP to avoid gaps
    sizePixels = ceil(sizePixels / SNAP) * SNAP;

    gl_PointSize = clamp(sizePixels, 1.0, 128.0);
}
`


/*
  0   8   2  10
 12   4  14   6
  3  11   1   9
 15   7  13   5
*/
 
var orderedDitherProgram uint32
var orderedDitherSource = `
#version 120

uniform sampler2D tex;


float bayer4x4(vec2 p)
{
    int x = int(mod(p.x, 4.0));
    int y = int(mod(p.y, 4.0));
    int i = x + y * 4;

    if (i ==  0) return  0.0/16.0;
    if (i ==  1) return  8.0/16.0;
    if (i ==  2) return  2.0/16.0;
    if (i ==  3) return 10.0/16.0;
    if (i ==  4) return 12.0/16.0;
    if (i ==  5) return  4.0/16.0;
    if (i ==  6) return 14.0/16.0;
    if (i ==  7) return  6.0/16.0;
    if (i ==  8) return  3.0/16.0;
    if (i ==  9) return 11.0/16.0;
    if (i == 10) return  1.0/16.0;
    if (i == 11) return  9.0/16.0;
    if (i == 12) return 15.0/16.0;
    if (i == 13) return  7.0/16.0;
    if (i == 14) return 13.0/16.0;
    return 5.0/16.0;
}

float bayer8x8(vec2 p)
{
    int x = int(mod(p.x, 8.0));
    int y = int(mod(p.y, 8.0));
    int i = x + y * 8;

    if (i ==  0) return  0.0/64.0;
    if (i ==  1) return 48.0/64.0;
    if (i ==  2) return 12.0/64.0;
    if (i ==  3) return 60.0/64.0;
    if (i ==  4) return  3.0/64.0;
    if (i ==  5) return 51.0/64.0;
    if (i ==  6) return 15.0/64.0;
    if (i ==  7) return 63.0/64.0;

    if (i ==  8) return 32.0/64.0;
    if (i ==  9) return 16.0/64.0;
    if (i == 10) return 44.0/64.0;
    if (i == 11) return 28.0/64.0;
    if (i == 12) return 35.0/64.0;
    if (i == 13) return 19.0/64.0;
    if (i == 14) return 47.0/64.0;
    if (i == 15) return 31.0/64.0;

    if (i == 16) return  8.0/64.0;
    if (i == 17) return 56.0/64.0;
    if (i == 18) return  4.0/64.0;
    if (i == 19) return 52.0/64.0;
    if (i == 20) return 11.0/64.0;
    if (i == 21) return 59.0/64.0;
    if (i == 22) return  7.0/64.0;
    if (i == 23) return 55.0/64.0;

    if (i == 24) return 40.0/64.0;
    if (i == 25) return 24.0/64.0;
    if (i == 26) return 36.0/64.0;
    if (i == 27) return 20.0/64.0;
    if (i == 28) return 43.0/64.0;
    if (i == 29) return 27.0/64.0;
    if (i == 30) return 39.0/64.0;
    if (i == 31) return 23.0/64.0;

    if (i == 32) return  2.0/64.0;
    if (i == 33) return 50.0/64.0;
    if (i == 34) return 14.0/64.0;
    if (i == 35) return 62.0/64.0;
    if (i == 36) return  1.0/64.0;
    if (i == 37) return 49.0/64.0;
    if (i == 38) return 13.0/64.0;
    if (i == 39) return 61.0/64.0;

    if (i == 40) return 34.0/64.0;
    if (i == 41) return 18.0/64.0;
    if (i == 42) return 46.0/64.0;
    if (i == 43) return 30.0/64.0;
    if (i == 44) return 33.0/64.0;
    if (i == 45) return 17.0/64.0;
    if (i == 46) return 45.0/64.0;
    if (i == 47) return 29.0/64.0;

    if (i == 48) return 10.0/64.0;
    if (i == 49) return 58.0/64.0;
    if (i == 50) return  6.0/64.0;
    if (i == 51) return 54.0/64.0;
    if (i == 52) return  9.0/64.0;
    if (i == 53) return 57.0/64.0;
    if (i == 54) return  5.0/64.0;
    if (i == 55) return 53.0/64.0;

    if (i == 56) return 42.0/64.0;
    if (i == 57) return 26.0/64.0;
    if (i == 58) return 38.0/64.0;
    if (i == 59) return 22.0/64.0;
    if (i == 60) return 41.0/64.0;
    if (i == 61) return 25.0/64.0;
    if (i == 62) return 37.0/64.0;
    return 21.0/64.0;
}



float dither4bit(float c, float t)
{
    /* default 16 */
    /* biblically accurate 2 */
    float levels = 5.0;
    float scaled = c * levels;
    float base   = floor(scaled);
    float error  = scaled - base;
    if (error > t)
        base += 1.0;

    base = min(base, levels - 1.0);
    return base / (levels - 1.0);
}
void main()
{
    vec3 color = texture2D(tex, gl_TexCoord[0].st).rgb;

    // Optional gamma correctness:
    color = pow(color, vec3(2.2));

    //~ float t = bayer8x8(gl_FragCoord.xy);
    
    float t = bayer4x4(gl_FragCoord.xy);

    color.r = dither4bit(color.r, t);
    color.g = dither4bit(color.g, t);
    color.b = dither4bit(color.b, t);

    color = pow(color, vec3(1.0 / 2.2));

    gl_FragColor = vec4(color, 1.0);
}
`


func createFragmentVertexProgram(vsSource, fsSource string) uint32 {
	compile := func(src string, stype uint32) uint32 {
		sh := gl.CreateShader(stype)
		csources, free := gl.Strs(src + "\x00")
		gl.ShaderSource(sh, 1, csources, nil)
		free()
		gl.CompileShader(sh)

		var status int32
		gl.GetShaderiv(sh, gl.COMPILE_STATUS, &status)
		if status == gl.FALSE {
			var logLen int32
			gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &logLen)
			log := strings.Repeat("\x00", int(logLen+1))
			gl.GetShaderInfoLog(sh, logLen, nil, gl.Str(log))
			panic(log)
		}
		return sh
	}

	vs := compile(vsSource, gl.VERTEX_SHADER)
	fs := compile(fsSource, gl.FRAGMENT_SHADER)

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)

	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		panic("program link failed")
	}

	return prog
}



func createFragmentProgram(fsSource string) uint32 {
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free := gl.Strs(fsSource + "\x00")
	gl.ShaderSource(fs, 1, csources, nil)
	free()
	gl.CompileShader(fs)

	var status int32
	gl.GetShaderiv(fs, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(fs, gl.INFO_LOG_LENGTH, &logLen)
		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(fs, logLen, nil, gl.Str(log))
		panic(log)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)

	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		panic("program link failed")
	}

	return prog
}



var (
    fbo        uint32
    colorTex   uint32
    depthRb    uint32

    fboHalfRatio int = 4
    fboHalf      uint32
    colorTexHalf uint32
    depthRbHalf  uint32
    
    fboAccumRatio int = 4
    fboAccum       uint32
    colorTexAccum  uint32
)

func CreateFBO(width, height int) error {
    if colorTex != 0 {
        gl.DeleteTextures(1, &colorTex);
        colorTex = 0;
    }
    if depthRb != 0 {
        gl.DeleteRenderbuffers(1, &depthRb);
        depthRb = 0;
    }
    if fbo != 0 {
        gl.DeleteFramebuffers(1, &fbo);
        fbo = 0;
    }
    if colorTexHalf != 0 {
        gl.DeleteTextures(1, &colorTexHalf);
        colorTexHalf = 0;
    }
    if depthRbHalf != 0 {
        gl.DeleteRenderbuffers(1, &depthRbHalf);
        depthRbHalf = 0;
    }
    if fboHalf != 0 {
        gl.DeleteFramebuffers(1, &fboHalf);
        fboHalf = 0;
    }
    
    if fboAccum != 0 {
        gl.DeleteRenderbuffers(1, &fboAccum);
        fboAccum = 0;
    }
    if colorTexAccum != 0 {
        gl.DeleteTextures(1, &colorTexAccum);
        colorTexAccum = 0;
    }


    // Full-resolution FBO
    gl.GenFramebuffers(1, &fbo)
    gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)

    gl.GenTextures(1, &colorTex)
    gl.BindTexture(gl.TEXTURE_2D, colorTex)
    gl.TexImage2D(gl.TEXTURE_2D, 0,
        gl.RGBA8, int32(widthInternal), int32(heightInternal), 0,
        gl.RGBA, gl.UNSIGNED_BYTE, nil)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
    gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
        gl.TEXTURE_2D, colorTex, 0)

    gl.GenRenderbuffers(1, &depthRb)
    gl.BindRenderbuffer(gl.RENDERBUFFER, depthRb)
    gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8,
        int32(widthInternal), int32(heightInternal))
    gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT,
        gl.RENDERBUFFER, depthRb)

    status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
    if status != gl.FRAMEBUFFER_COMPLETE {
        return fmt.Errorf("FBO 1 incomplete: 0x%X", status)
    }

    // Lower-resolution FBO
    halfW := int32(widthInternal / fboHalfRatio)
    halfH := int32(heightInternal / fboHalfRatio)

    gl.GenFramebuffers(1, &fboHalf)
    gl.BindFramebuffer(gl.FRAMEBUFFER, fboHalf)

    // Color texture
    gl.GenTextures(1, &colorTexHalf)
    gl.BindTexture(gl.TEXTURE_2D, colorTexHalf)
    
    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8,
        halfW, halfH, 0,
        gl.RGBA, gl.UNSIGNED_BYTE, nil)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

    gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
        gl.TEXTURE_2D, colorTexHalf, 0)

    // Depth/stencil renderbuffer
    gl.GenRenderbuffers(1, &depthRbHalf)
    gl.BindRenderbuffer(gl.RENDERBUFFER, depthRbHalf)
    gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8,
        halfW, halfH)

    gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT,
        gl.RENDERBUFFER, depthRbHalf)

    status = gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
    if status != gl.FRAMEBUFFER_COMPLETE {
        return fmt.Errorf("FBO 2 incomplete: 0x%X", status)
    }


    // Ray casting FBO
    gl.GenFramebuffers(1, &fboAccum)
    gl.BindFramebuffer(gl.FRAMEBUFFER, fboAccum)

    gl.GenTextures(1, &colorTexAccum)
    gl.BindTexture(gl.TEXTURE_2D, colorTexAccum)
    gl.TexImage2D(gl.TEXTURE_2D, 0,
        gl.RGBA8, int32(widthInternal / fboAccumRatio), int32(heightInternal / fboAccumRatio), 0,
        gl.RGBA, gl.UNSIGNED_BYTE, nil)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
    gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0,
        gl.TEXTURE_2D, colorTexAccum, 0)
        
    gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT,
        gl.RENDERBUFFER, depthRb)

    status = gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
    if status != gl.FRAMEBUFFER_COMPLETE {
        return fmt.Errorf("FBO 3 incomplete: 0x%X", status)
    }


    gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
    return nil
}


func RenderFBO(tex uint32, offset_x, offset_y float32) {

    gl.MatrixMode(gl.PROJECTION)
    gl.PushMatrix()
    gl.LoadIdentity()

    gl.MatrixMode(gl.MODELVIEW)
    gl.PushMatrix()
    gl.LoadIdentity()

    // at the current moment this is the only point we can draw a small quad
    // right a the middle of the screen
    // TODO: put this somewhere else
    gl.Begin(gl.QUADS)
        gl.Vertex3f(-1*0.02, -1*0.02, 0);
        gl.Vertex3f( 1*0.02, -1*0.02, 0);
        gl.Vertex3f( 1*0.02,  1*0.02, 0);
        gl.Vertex3f(-1*0.02,  1*0.02, 0);
    gl.End()

    // Enable 2D texturing
    gl.Enable(gl.TEXTURE_2D)
    gl.BindTexture(gl.TEXTURE_2D, tex)

    // Draw screen-sized quad
    gl.Begin(gl.QUADS)

    gl.TexCoord2f(0.0, 0.0); gl.Vertex2f(-1.0 + offset_x, -1.0 + offset_y)
    gl.TexCoord2f(1.0, 0.0); gl.Vertex2f( 1.0 + offset_x, -1.0 + offset_y)
    gl.TexCoord2f(1.0, 1.0); gl.Vertex2f( 1.0 + offset_x,  1.0 + offset_y)
    gl.TexCoord2f(0.0, 1.0); gl.Vertex2f(-1.0 + offset_x,  1.0 + offset_y)

    gl.End()

    // Cleanup
    gl.Disable(gl.TEXTURE_2D)

    // Restore matrices
    gl.MatrixMode(gl.PROJECTION)
    gl.PopMatrix()

    gl.MatrixMode(gl.MODELVIEW)
    gl.PopMatrix()
}
///



type Mat4 [16]float32

func Perspective (fovyDeg, aspect, zNear, zFar float32) Mat4 {
    fovyRad := fovyDeg * (math.Pi / 180.0)
    f := float32(1.0 / math.Tan(float64(fovyRad/2)))

    var m Mat4
    m[0]  = f / aspect
    m[5]  = f
    m[10] = (zFar + zNear) / (zNear - zFar)
    m[11] = -1
    m[14] = (2 * zFar * zNear) / (zNear - zFar)
    return m
}

func Invert (m Mat4) (Mat4, bool) {
    var inv Mat4

    inv[0] = m[5]*m[10]*m[15] -
        m[5]*m[11]*m[14] -
        m[9]*m[6]*m[15] +
        m[9]*m[7]*m[14] +
        m[13]*m[6]*m[11] -
        m[13]*m[7]*m[10]

    inv[4] = -m[4]*m[10]*m[15] +
        m[4]*m[11]*m[14] +
        m[8]*m[6]*m[15] -
        m[8]*m[7]*m[14] -
        m[12]*m[6]*m[11] +
        m[12]*m[7]*m[10]

    inv[8] = m[4]*m[9]*m[15] -
        m[4]*m[11]*m[13] -
        m[8]*m[5]*m[15] +
        m[8]*m[7]*m[13] +
        m[12]*m[5]*m[11] -
        m[12]*m[7]*m[9]

    inv[12] = -m[4]*m[9]*m[14] +
        m[4]*m[10]*m[13] +
        m[8]*m[5]*m[14] -
        m[8]*m[6]*m[13] -
        m[12]*m[5]*m[10] +
        m[12]*m[6]*m[9]

    inv[1] = -m[1]*m[10]*m[15] +
        m[1]*m[11]*m[14] +
        m[9]*m[2]*m[15] -
        m[9]*m[3]*m[14] -
        m[13]*m[2]*m[11] +
        m[13]*m[3]*m[10]

    inv[5] = m[0]*m[10]*m[15] -
        m[0]*m[11]*m[14] -
        m[8]*m[2]*m[15] +
        m[8]*m[3]*m[14] +
        m[12]*m[2]*m[11] -
        m[12]*m[3]*m[10]

    inv[9] = -m[0]*m[9]*m[15] +
        m[0]*m[11]*m[13] +
        m[8]*m[1]*m[15] -
        m[8]*m[3]*m[13] -
        m[12]*m[1]*m[11] +
        m[12]*m[3]*m[9]

    inv[13] = m[0]*m[9]*m[14] -
        m[0]*m[10]*m[13] -
        m[8]*m[1]*m[14] +
        m[8]*m[2]*m[13] +
        m[12]*m[1]*m[10] -
        m[12]*m[2]*m[9]

    inv[2] = m[1]*m[6]*m[15] -
        m[1]*m[7]*m[14] -
        m[5]*m[2]*m[15] +
        m[5]*m[3]*m[14] +
        m[13]*m[2]*m[7] -
        m[13]*m[3]*m[6]

    inv[6] = -m[0]*m[6]*m[15] +
        m[0]*m[7]*m[14] +
        m[4]*m[2]*m[15] -
        m[4]*m[3]*m[14] -
        m[12]*m[2]*m[7] +
        m[12]*m[3]*m[6]

    inv[10] = m[0]*m[5]*m[15] -
        m[0]*m[7]*m[13] -
        m[4]*m[1]*m[15] +
        m[4]*m[3]*m[13] +
        m[12]*m[1]*m[7] -
        m[12]*m[3]*m[5]

    inv[14] = -m[0]*m[5]*m[14] +
        m[0]*m[6]*m[13] +
        m[4]*m[1]*m[14] -
        m[4]*m[2]*m[13] -
        m[12]*m[1]*m[6] +
        m[12]*m[2]*m[5]

    inv[3] = -m[1]*m[6]*m[11] +
        m[1]*m[7]*m[10] +
        m[5]*m[2]*m[11] -
        m[5]*m[3]*m[10] -
        m[9]*m[2]*m[7] +
        m[9]*m[3]*m[6]

    inv[7] = m[0]*m[6]*m[11] -
        m[0]*m[7]*m[10] -
        m[4]*m[2]*m[11] +
        m[4]*m[3]*m[10] +
        m[8]*m[2]*m[7] -
        m[8]*m[3]*m[6]

    inv[11] = -m[0]*m[5]*m[11] +
        m[0]*m[7]*m[9] +
        m[4]*m[1]*m[11] -
        m[4]*m[3]*m[9] -
        m[8]*m[1]*m[7] +
        m[8]*m[3]*m[5]

    inv[15] = m[0]*m[5]*m[10] -
        m[0]*m[6]*m[9] -
        m[4]*m[1]*m[10] +
        m[4]*m[2]*m[9] +
        m[8]*m[1]*m[6] -
        m[8]*m[2]*m[5]

    det := m[0]*inv[0] + m[1]*inv[4] + m[2]*inv[8] + m[3]*inv[12]
    if det == 0 {
        return Mat4{}, false
    }

    det = 1.0 / det
    for i := 0; i < 16; i++ {
        inv[i] *= det
    }

    return inv, true
}
        
        
func UploadMat4 (program uint32, name string, m Mat4) {
    loc := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
    gl.UniformMatrix4fv(loc, 1, false, &m[0])
}
        
        
func AdjustResolution(window *glfw.Window) {
    nw, nh := window.GetFramebufferSize()
    if (nw != width) || (nh != height) {
        width = nw 
        height = nh
        CreateFBO(width, height)
    }
}
        
        
var bloom_on bool = true


func PlayWav(ctx *malgo.AllocatedContext, filename string) {
    var e error
    f := func() {
        file, err := os.Open(filename)
        if err != nil {
            e = err
            return
        }
        defer file.Close()

        w := wav.NewReader(file)

        format, err := w.Format()
        if err != nil {
            e = err
            return
        }

        reader := io.Reader(w)
        channels := uint32(format.NumChannels)
        sampleRate := format.SampleRate

        deviceConfig := malgo.DefaultDeviceConfig(malgo.Playback)
        deviceConfig.Playback.Format = malgo.FormatS16
        deviceConfig.Playback.Channels = channels
        deviceConfig.SampleRate = sampleRate
        deviceConfig.Alsa.NoMMap = 1

        onSamples := func(pOutputSample, _ []byte, frameCount uint32) {
            _, err := io.ReadFull(reader, pOutputSample)
            if err != nil {
                // Silence output on EOF
                for i := range pOutputSample {
                    pOutputSample[i] = 0
                }
            }
        }

        deviceCallbacks := malgo.DeviceCallbacks{
            Data: onSamples,
        }

        device, err := malgo.InitDevice(ctx.Context, deviceConfig, deviceCallbacks)
        if err != nil {
            e = err
            return
        }
        defer device.Uninit()

        if err := device.Start(); err != nil {
            e = err
            return
        }

        // Block until playback finishes
        duration, _ := w.Duration()
        time.Sleep(duration)
        return
    }
    
    go f()
    if e != nil {
        panic(e)
    }
}



func main() {
    

/*
    if len(os.Args) > 1 {
        switch (os.Args[1]) {
            
            case "-h","-help":
fmt.Println(
`Usage: dmc [options]

DMC a block game inspired by Minecraft infdev and Alpha.
Options:
  -h, --help      show available commands
  -s, --save      choose savegame name
  -r, --render    choose render distance`)
                os.Exit(0)
            case "-s","-save":
                os.Exit(0)
            default:
        }
    }
*/


    if err := glfw.Init(); err != nil {
        log.Fatalln("failed to initialize glfw:", err)
    }
    defer glfw.Terminate()
        
    glfw.WindowHint(glfw.ContextVersionMajor, 3);
    glfw.WindowHint(glfw.ContextVersionMinor, 3);
    glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile);
    glfw.WindowHint(glfw.OpenGLForwardCompatible, 0);

    window, err := glfw.CreateWindow(800, 600, "DMC", nil, nil)
    if err != nil {
        log.Fatalln("failed to create window:", err)
    }
    window.MakeContextCurrent()
    if err := gl.Init(); err != nil {
        log.Fatalln("failed to initialize OpenGL:", err)
    }
    //~ gl.Enable(gl.DEBUG_OUTPUT)
    LoadAtlasLOD("assets/atlas.png",16)
    atlas = LoadTexture("assets/atlas.png")
    
    Background = LoadTexture("assets/splash.png")
    TextureSun = LoadTexture("assets/sun.png")
    
    window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
    window.SetCursorPosCallback(mouseCursorCallback)
    
    ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, nil)
    if err != nil {
        panic(err)
    }
    defer ctx.Uninit()

    //~ PlayWav(ctx, "/assets/sine.wav")

    window.SetScrollCallback(func(w *glfw.Window, xoff, yoff float64) {
        if yoff < 0 { if selectedSlot > 0 { selectedSlot-- } }
        if yoff > 0 { if selectedSlot < (len(player_inventory.C)-1) { selectedSlot++ } }
    })
    
    MainMenu(window)
    
    player_pos = Vec3{1000+float32(rand.Uint32()%9000),CHUNK_V,1000+float32(rand.Uint32()%9000)}
    
    if _, err := os.Stat(fmt.Sprintf("%s/level.bin", saveSlot)); err == nil {
        log.Println("loaded save in ",saveSlot)
        saveReadPlayerdata()
    }
    
    playerInventoryInsert(CRAFTING_TABLE, 1)
    playerInventoryInsert(SIGN, 16)
    //~ playerInventoryInsert(XRAY, 16)
    //~ playerInventoryInsert(CANDLE, 1)
    //~ playerInventoryInsert(STICK, 1)
    //~ playerInventoryInsert(GRASS_TUFF, 1)
    //~ playerInventoryInsert(LAVENDER, 1)
    //~ playerInventoryInsert(RAYTRACED_SPHERE, 16)
    //~ playerInventoryInsert(PAINTING, 16)
    
    CreateFBO(width, height)
    orderedDitherProgram = createFragmentProgram(orderedDitherSource)
    subpixelSnapProgram = createFragmentVertexProgram(subpixelSnapSource, minimumFragmentSource)
    loc := gl.GetUniformLocation(subpixelSnapProgram, gl.Str("uViewport\x00"))
    subpixelSnapUniformTextured = gl.GetUniformLocation(subpixelSnapProgram, gl.Str("uBillboard\x00"))
    subpixelSnapUniformSeed = gl.GetUniformLocation(subpixelSnapProgram, gl.Str("uSeed\x00"))

    RaytracingProgram = createFragmentVertexProgram(RaytracingSourceV, RaytracingSourceF)
    
    var t uint32 = 0;
    for !window.ShouldClose() {
        const DEBUG = false
        var times []time.Time
        
        if DEBUG { 
            times = append(times, time.Now())
        }
            
        AdjustResolution(window)
        
        t++
        
        if bloom_on {
            
            /* for now the shader breaks untextured polygons */
            gl.UseProgram(subpixelSnapProgram)
            //~ gl.Uniform2f(loc, float32(width-width%2), float32(height-height%2))
            gl.Uniform2f(loc, float32(widthInternal-widthInternal%2), float32(heightInternal-heightInternal%2))
            // Render into FBO
            gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
            //~ gl.Viewport(0, 0, int32(widthInternal), int32(heightInternal))
            setupScene(window,int32(widthInternal), int32(heightInternal))
        } else {
            setupScene(window, int32(width), int32(height))
        }
    
        gl.Enable(gl.CULL_FACE);
        gl.CullFace(gl.BACK);      // or GL_FRONT
        gl.FrontFace(gl.CCW);     // or GL_CW
        drawScene()
        

        processBlockBreaking()
        
        


        
        if (t%6) == 0 {
            if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
                for key := glfw.KeyA; key <= glfw.KeyZ; key++ {
                    if window.GetKey(key) == glfw.Press {
                        letter := byte(rune('a' + (key - glfw.KeyA)))
                        //~ log.Println("pressed: ",letter,"\n")
                        player_textbuffer = append(player_textbuffer, letter)
                    }
                }
            }
        }
        
        if (t%3) == 0 {
            processBlockUpdates()
            TICK++
        }

        //~ TICK++

        playerCollision(window) 

        if DEBUG { 
            times = append(times, time.Now())
        }
        
        gl.BindFramebuffer(gl.READ_FRAMEBUFFER, fbo);
        gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, fboAccum);
        gl.BlitFramebuffer(
            0, 0, int32(widthInternal), int32(heightInternal),   // destination rect
            0, 0, int32(widthInternal/fboAccumRatio), int32(heightInternal/fboAccumRatio),   // source rect
            gl.DEPTH_BUFFER_BIT,
            gl.NEAREST);
        gl.BindFramebuffer(gl.FRAMEBUFFER, fboAccum)
        gl.Viewport(0, 0, int32(widthInternal/fboAccumRatio), int32(heightInternal/fboAccumRatio))
        gl.Disable(gl.BLEND);
        gl.ClearColor(0,0,0,0)
        gl.Clear(gl.COLOR_BUFFER_BIT)
        drawRaytraced()
        
        if DEBUG { 
            times = append(times, time.Now())
        }
        
        // EXPERIMENTAL
        if bloom_on {
            
            DrawTexturedQuad := func (tex uint32, ox, oy float32) {
                gl.BindTexture(gl.TEXTURE_2D, tex)
                gl.Begin(gl.QUADS)
                    gl.TexCoord2f(0, 0); gl.Vertex2f(-1+ox, -1+oy)
                    gl.TexCoord2f(1, 0); gl.Vertex2f( 1+ox, -1+oy)
                    gl.TexCoord2f(1, 1); gl.Vertex2f( 1+ox,  1+oy)
                    gl.TexCoord2f(0, 1); gl.Vertex2f(-1+ox,  1+oy)
                gl.End()
            }
            
            gl.MatrixMode(gl.PROJECTION)
            gl.LoadIdentity()

            gl.MatrixMode(gl.MODELVIEW)
            gl.LoadIdentity()

            gl.Disable(gl.DEPTH_TEST)
            gl.Enable(gl.TEXTURE_2D)

            gl.BindFramebuffer(gl.FRAMEBUFFER, fboHalf)
            gl.Viewport(0, 0, int32(widthInternal/fboHalfRatio), int32(heightInternal/fboHalfRatio))
            gl.Disable(gl.BLEND)
            gl.Color4f(1, 1, 1, 1)
            DrawTexturedQuad(colorTex, 0, 0)

            gl.BindFramebuffer(gl.FRAMEBUFFER, fbo)
            gl.Viewport(0, 0, int32(widthInternal), int32(heightInternal))
            gl.Enable(gl.BLEND)
            gl.BlendFunc(gl.SRC_COLOR, gl.ONE)

            filterStrength := float32(0.305)
            gl.Color4f(filterStrength, filterStrength, filterStrength, 1)

            blurX := float32(16.0) / float32(widthInternal)
            blurY := float32(16.0) / float32(heightInternal)

            DrawTexturedQuad(colorTexHalf,  blurX, 0)
            DrawTexturedQuad(colorTexHalf, -blurX, 0)
            DrawTexturedQuad(colorTexHalf, 0,  blurY)
            DrawTexturedQuad(colorTexHalf, 0,  -blurY)
            //~ DrawTexturedQuad(colorTexHalf, 0, 0)
            
            gl.Color4f(1, 1, 1, 1)
            gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
            DrawTexturedQuad(colorTexAccum, 0, 0)
            
            gl.Disable(gl.BLEND)

            gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
            gl.Viewport(0, 0, int32(width), int32(height))

            gl.UseProgram(orderedDitherProgram)
            levelsLoc := gl.GetUniformLocation(orderedDitherProgram, gl.Str("levels\x00"))
            gl.Uniform1f(levelsLoc, 31.0)

            DrawTexturedQuad(colorTex, 0, 0)
        }
        gl.UseProgram(0)
        
        gl.Enable(gl.BLEND)
        gl.BlendFunc(gl.ONE_MINUS_DST_COLOR, gl.ONE_MINUS_SRC_COLOR);
        DrawInventory()
        gl.Disable(gl.BLEND)
        
        if DEBUG { 
            times = append(times, time.Now())
            for k, v := range times {
                fmt.Println(k, time.Since(v))
            }
        }
        

        
        if bloom_on {
            gl.Enable(gl.DEPTH_TEST)
        }

        if (window.GetKey(glfw.KeyB) == glfw.Press) && !holding_b {
            bloom_on = !bloom_on;
            holding_b = true
        }
        if (window.GetKey(glfw.KeyB) == glfw.Release) {
            holding_b = false
        }
        
        if DEBUG {
            fmt.Println("total: ",time.Since(times[0]))
        }


        window.SwapBuffers()
    
        glfw.PollEvents()
        if window.GetKey(glfw.KeyEscape) == glfw.Press {
            window.SetShouldClose(true)
        }
    }
}

type Achievement struct {
    name []byte
    icon int
}

var title_d = [][]byte{
    []byte("daemon"),
    []byte("dancing"),
    []byte("dartmouth"),
    []byte("darwin"),
    []byte("data"),
    []byte("dictionary"),
    []byte("date"),
    []byte("doc"),
    []byte("dotcom"),
    []byte("dead"),
    []byte("debug"),
    []byte("decay"),
    []byte("desktop"),
    []byte("degree"),   
    []byte("deflate"),
    []byte("design"),
    []byte("disk"),
    []byte("dither"),
    []byte("domain"),
    []byte("dweeb"),
}

var title_m = [][]byte{
    []byte("machine"),
    []byte("macro"),
    []byte("magic"),
    []byte("magnetic"),
    []byte("mail"),
    []byte("malloc"),
    []byte("mars"),
    []byte("matrix"),
    []byte("meme"),
    []byte("meta"),
    []byte("micro"),
    []byte("miner"),
    []byte("monitor"),
    []byte("munch"),
}

var title_c = [][]byte{
    []byte("call"),
    []byte("calc"),
    []byte("card"),
    []byte("case"),
    []byte("cast"),
    []byte("cell"),
    []byte("char"),
    []byte("coax"),
    []byte("code"),
    []byte("com"),
    []byte("command"),
    []byte("computer"),
    []byte("copper"),
    []byte("clone"),
    []byte("cyber"),
}

func MainMenu (window *glfw.Window) {
    
    var state int = 0
    var selected int = 0
    
    var holding_up bool
    var holding_down bool
    var holding_enter bool
    
    var title string
    newtitle := func() {
        title = string(title_d[rand.Intn(len(title_d))])
        title += " "
        title += string(title_m[rand.Intn(len(title_m))])
        title += " "
        title += string(title_c[rand.Intn(len(title_c))])
    }
    tick := 0
    
    
    for !window.ShouldClose() {

        
        
        if (window.GetKey(glfw.KeyUp) == glfw.Press) && !holding_up {
            selected--
            holding_up = true
        }
        if (window.GetKey(glfw.KeyUp) == glfw.Release) && holding_up {
            holding_up = false
        }

        if (window.GetKey(glfw.KeyDown) == glfw.Press) && !holding_down {
            selected++
            holding_down = true
        }
        if (window.GetKey(glfw.KeyDown) == glfw.Release) && holding_down {
            holding_down = false
        }

        if (window.GetKey(glfw.KeyEnter) == glfw.Press) && !holding_enter {
            holding_enter = true
        }
        if (window.GetKey(glfw.KeyEnter) == glfw.Release) && holding_enter {
            holding_enter = false
        }



        AdjustResolution(window)
        
        gl.Viewport(0, 0, int32(width), int32(height))
        
        gl.MatrixMode(gl.PROJECTION)
        gl.PushMatrix()
        gl.LoadIdentity()

        gl.MatrixMode(gl.MODELVIEW)
        gl.PushMatrix()
        gl.LoadIdentity()

        gl.ClearColor(0.529, 0.808, 1.0, 1.0)
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.Enable(gl.BLEND)
        gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);


        gl.Enable(gl.TEXTURE_2D);
        
        gl.BindTexture(gl.TEXTURE_2D, Background)
        gl.Begin(gl.QUADS)
            gl.TexCoord2f(0.0, 0.0); gl.Vertex2f( 1,  1)
            gl.TexCoord2f(1.0, 0.0); gl.Vertex2f(-1,  1)
            gl.TexCoord2f(1.0, 1.0); gl.Vertex2f(-1, -1)
            gl.TexCoord2f(0.0, 1.0); gl.Vertex2f( 1, -1)
        gl.End()
        
        if selected < 0 {
            selected = 0 
        }
        
        
        if state == 0 {
            if selected > 3 {
               selected = 3 
            }
            
            menu := [][]byte{
                []byte("load"),
                []byte("tutorial"),
                []byte("achievements"),
                []byte("quit"),
            }

            gl.BindTexture(gl.TEXTURE_2D, atlas)        
            gl.Begin(gl.QUADS)
                        
            DrawText( []byte(title), 0, 5,float32(0.05))
            
            if (tick % 60) == 0 {
                newtitle()
            }
            
            for k, _ := range menu {
                if k == selected {
                    gl.Color3f(0.75, 0.75, 0.75)
                } else {
                    gl.Color3f(1.0, 1.0, 1.0)
                }
                DrawText(menu[k], 0, float32(-k),float32(0.05))
            }
            gl.End()

            if (window.GetKey(glfw.KeyEnter) == glfw.Press) {
                if selected == 0 {
                    break
                } else if selected == 1 {
                    saveSlot = "saves/tutorial"
                    break
                } else if selected == 2 {
                    
                } else if selected == 3 {
                    os.Exit(0)
                }
            }

            
        } else if state == 1 {
            
            
            
        }
    
        
        gl.Color3f(1.0, 1.0, 1.0)
        
        // Restore matrices
        gl.MatrixMode(gl.PROJECTION)
        gl.PopMatrix()

        gl.MatrixMode(gl.MODELVIEW)
        gl.PopMatrix()

        window.SwapBuffers()
        glfw.PollEvents()
        
        tick++

    }
}

func drawRaytraced() {
    if schedule_raytraced_blocks_p <= 0 {
        return
    }
    
    // ray tracing 
    gl.UseProgram(RaytracingProgram)

    var view Mat4
    gl.GetFloatv(gl.MODELVIEW_MATRIX, &view[0])

    invView, ok := Invert(view)
    if !ok { panic("view mauInvViewtrix not invertible") }
    //~ log.Println(view)
    UploadMat4(RaytracingProgram, "uInvView", invView)
    proj := Perspective(60.0, float32(width)/float32(height), 0.1, 100.0)
    invProj, ok := Invert(proj)
    if !ok {
        panic("projection matrix not invertible")
    }
    UploadMat4(RaytracingProgram, "uInvProjection", invProj)

    RaytracingUniformuSphereCenter := gl.GetUniformLocation(RaytracingProgram, gl.Str("uSphereCenter\x00"))
    
    gl.Uniform4fv(RaytracingUniformuSphereCenter, int32(schedule_raytraced_blocks_p), &schedule_raytraced_blocks[0])
    RaytracingUniformSeed := gl.GetUniformLocation(RaytracingProgram, gl.Str("Seed\x00"))
    gl.Uniform1f(RaytracingUniformSeed, float32(TICK))

    gl.Begin(gl.QUADS)
        gl.Vertex3f(-1, -1, 0);
        gl.Vertex3f( 1, -1, 0);
        gl.Vertex3f( 1,  1, 0);
        gl.Vertex3f(-1,  1, 0);
    gl.End()
    gl.UseProgram(subpixelSnapProgram)
}


var holding_b bool

var processBlockBreaking_buffer []struct{x, y, z, t, id, elapsed int}
func processBlockBreaking() {
    breakTimes := [N_BLOCK_IDS]int{
        LEAF:        6*3,
        DIRT:       15*3,
        GRASS:      15*3,
        SAND:       15*3,
        GRAVEL:     15*3,
        SNOW:       15*3,
        ICE:        15*3,
        MOON:       15*3,
        CLAY:       15*3,
        WOOD:       60*3,
        STONE:     150*3,
        COBBLE:    150*3,
        ORE_COAL:  300*3,
        ORE_IRON:  300*3,
        FLESH:    1500*3,
        
        CRAFTING_TABLE: 60*3,
        FURNACE:       150*3,
        PLANK: 60*3,
        
        RAYTRACED_SPHERE: 20000000,
        BEDROCK:20000000,
    }
    
    var drawBreaking = func (x, y, z, t int, target BlockID) {
        if breakTimes[target] == 0 {
            return
        }
        
        
        //~ log.Println(x,y,z,t)
        texHPos := func ( face int) float32 {
            return float32(face)*(16.0/1024.0)
        }
        textVPos := func ( id BlockID ) float32 {
            return float32(id)*(16.0/1024.0)
        }
        
        //~ var id BlockID = BREAKING
        var id BlockID = BREAKING
        var n = [6][3]int{
            { 1, 0, 0},
            {-1, 0, 0},
            { 0, 1, 0},
            { 0,-1, 0},
            { 0, 0, 1},
            { 0, 0,-1},
        }
        
        
        l := float32(-0.01)
        h := float32( 1.01)
        var v = [6][4][3]float32{
            {{h, l, h}, {h, l, l}, {h, h, l}, {h, h, h}}, // +X
            {{l, l, l}, {l, l, h}, {l, h, h}, {l, h, l}}, // -X
            {{l, h, l}, {l, h, h}, {h, h, h}, {h, h, l}}, // +Y
            {{l, l, l}, {h, l, l}, {h, l, h}, {l, l, h}}, // -Y
            {{l, l, h}, {h, l, h}, {h, h, h}, {l, h, h}}, // +Z
            {{h, l, l}, {l, l, l}, {l, h, l}, {h, h, l}}, // -Z
        }

        for f := uint32(0); f < 6; f++ {
            nb := getBlockVal(x + n[f][0], y + n[f][1],z + n[f][2])
            

            if (nb != AIR) && (!BlockIsLiquid[nb]) {
                continue
            }
            
            lv := float32(lightLevels - getLightVal(x + n[f][0], y + n[f][1],z + n[f][2]))/float32(lightLevels)
            
            gl.Color3f(lv, lv, lv)
            
            v0 := textVPos(id  )
            v1 := textVPos(id+1)
            //~ h0 := texHPos(  (t / 10))
            //~ h1 := texHPos(1+(t / 10))
            h0 := texHPos(  (t / (breakTimes[target]/6) ))
            h1 := texHPos(1+(t / (breakTimes[target]/6) ))
            
            gl.TexCoord2f(h1, v1);
            gl.Vertex3f(float32(x)+v[f][0][0], float32(y)+v[f][0][1], float32(z)+v[f][0][2])
            gl.TexCoord2f(h0, v1);
            gl.Vertex3f(float32(x)+v[f][1][0], float32(y)+v[f][1][1], float32(z)+v[f][1][2])
            gl.TexCoord2f(h0, v0);
            gl.Vertex3f(float32(x)+v[f][2][0], float32(y)+v[f][2][1], float32(z)+v[f][2][2])
            gl.TexCoord2f(h1, v0);
            gl.Vertex3f(float32(x)+v[f][3][0], float32(y)+v[f][3][1], float32(z)+v[f][3][2])
        }
    }


    processBlockBreaking_buffer = processBlockBreaking_buffer[:0]
    gl.Enable(gl.TEXTURE_2D);
    gl.BindTexture(gl.TEXTURE_2D, atlas)    
    gl.Begin(gl.QUADS)

    for _, v := range breakingBlocks {
        b  := getBlockRef(int(v.x),int(v.y),int(v.z))
        if v.t >= (breakTimes[*b]+1) {
            
            
            if *b == STONE {
                playerInventoryInsert(COBBLE, 1)
            } else if *b == LEAF {
                
                if rand.Uint32()%4 == 0 {
                    playerInventoryInsert(SAPLING, 1)
                }
            } else if *b == ICE {
            } else {
                playerInventoryInsert(*b, 1)
            }
            *b = AIR
            chunkUpdateDisplaylist(int(v.x+1), int(v.y), int(v.z))
            chunkUpdateDisplaylist(int(v.x-1), int(v.y), int(v.z))
            chunkUpdateDisplaylist(int(v.x), int(v.y), int(v.z+1))
            chunkUpdateDisplaylist(int(v.x), int(v.y), int(v.z-1))
            updateSpread(v.x, v.y, v.z)
        } else {
            v.elapsed += 1
            drawBreaking(v.x, v.y, v.z, v.t, *b)
            
            if v.elapsed < 50 {
                processBlockBreaking_buffer = append(processBlockBreaking_buffer, v)
            }
        }
    }
    gl.End()
    gl.Disable(gl.TEXTURE_2D);
    
    breakingBlocks = breakingBlocks[:0]
    breakingBlocks = append(breakingBlocks[:0], processBlockBreaking_buffer...)
}


func MobSpawn(x, y, z int, id EntityID) {
    var foo Entity;
    foo.pos = Vec3{float32(x), float32(y), float32(z)}
    foo.id = id
    Entities = append(Entities, foo)
}


func processMobs() {
    
    //~ log.Println("%d\n", len(Entities))
    
    angle_to_entity := func (pos, dir, target Vec3) float32 {
        v := Vec3{target.x - pos.x, 0, target.z - pos.z}
        dot := float64(dir.x*v.x + dir.z*v.z)
        cross := float64(dir.x*v.z - dir.z*v.x)
        return float32(math.Atan2(cross, dot))
    }
    
    distance := func(p1, p2 Vec3) float32 {
        return (p1.x-p2.x)*(p1.x-p2.x)+(p1.y-p2.y)*(p1.y-p2.y)+(p1.z-p2.z)*(p1.z-p2.z)
    }
    
    w := 0
    for _, v := range Entities {
        if v.hp < 0 {
            continue
        }
        
        
        var walk_speed float32 = 0.1
        
        if v.state == MOB_IDLE {
            if (rand.Uint32()%10) == 0 {
                v.state = MOB_WALKING
            } else if (rand.Uint32()%10) == 0 {
                v.state = MOB_TURNING
            } else if (rand.Uint32()%10) == 0 {
                v.state = MOB_JUMPING
            }
        }
        
        if v.state == MOB_TURNING {
            v.dir = Vec3{walk_speed*(rand.Float32()-0.5), -0.3, walk_speed*(rand.Float32()-0.5)}
            v.state = MOB_IDLE
        }
        
        var valid_x, valid_y, valid_z bool
        
        mob_aabb_callback := func(b *BlockID, x, y, z int) { }
        valid_x, valid_y, valid_z, _ = AABB(v.pos.x,v.pos.y,v.pos.z,v.dir.x,v.dir.y,v.dir.z, 1.0, 0.3, mob_aabb_callback)


        if chunkExists(int(v.pos.x),int(v.pos.y),int(v.pos.z)) {
            if getBlockVal(int(v.pos.x),int(v.pos.y),int(v.pos.z)) == WATER {
                v.dir.y = 0.1
            }
        } else {
            continue
        }

        if v.state == MOB_JUMPING {
            
            if !valid_y {
                v.dir.y = 0.3
            }
            if (rand.Uint32()%4) == 0 {
                if v.dir.y > 0 {
                    v.dir.y = 0
                }
                v.state = MOB_IDLE
            }
        }
        
        if (v.state == MOB_WALKING) {
            //np := Vec3{walk_speed*(rand.Float32()-0.5), -0.3, walk_speed*(rand.Float32()-0.5)}
            
            if valid_x {
                v.pos.x += v.dir.x
            }
            if valid_y {
                //~ v.pos.y += v.dir.y
            }
            if valid_z {
                v.pos.z += v.dir.z
            }
            
            if (rand.Uint32()%10) == 0 {
                v.state = MOB_IDLE
            }
        }
        if valid_y {
            v.pos.y += v.dir.y
        }
        
        {
            light := getLightVal(int(v.pos.x),int(v.pos.y)+1,int(v.pos.z))
            lv := float32(lightLevels - light) / float32(lightLevels)
            gl.Color3f(lv,lv,lv)
        }
        
        if v.id == MOBID_DUCK {
            a := angle_to_entity(v.pos, v.dir, player_pos)
            pi := float32(3.14)
            if ( a < 0.2 ) && ( a > 0.2 ) {
                drawBillboard(DUCK, 0, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else if ( a > 0.2 ) && ( a < (pi-0.2) ) {
                drawBillboard(DUCK, 1, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else if ( a < -0.2 ) && ( a > (-pi+0.2)) {
                drawBillboard(DUCK, 2, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else {
                drawBillboard(DUCK, 3, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            }
        } else if v.id == MOBID_BEE {
            a := angle_to_entity(v.pos, v.dir, player_pos)
            pi := float32(3.14)
            if ( a < 0.2 ) && ( a > 0.2 ) {
                drawBillboard(BEE, 0, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else if ( a > 0.2 ) && ( a < (pi-0.2) ) {
                drawBillboard(BEE, 1, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else if ( a < -0.2 ) && ( a > (-pi+0.2)) {
                drawBillboard(BEE, 2, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            } else {
                drawBillboard(BEE, 3, v.pos.x, v.pos.y-0.5, v.pos.z, 1.0)
            }
        }
        
        if distance(v.pos,player_pos) < (SUBCHUNK_H*SUBCHUNK_H*RENDER_DISTANCE*RENDER_DISTANCE) {
            Entities[w] = v
            w++
        }
        
        gl.Color3f(1.0,1.0,1.0)
        
    }
    Entities = Entities[:w]
    
}

func mouseCursorCallback(window *glfw.Window, xpos, ypos float64) {
	if first_mouse {
		last_mouse_x = xpos
		last_mouse_y = ypos
		first_mouse = false
	}
	sensitivity := float32(0.00125)
	xoffset := float32(xpos - last_mouse_x) * sensitivity
	yoffset := float32(last_mouse_y - ypos) * sensitivity // Reversed Y
	last_mouse_x = xpos
	last_mouse_y = ypos
    
    
    if !openInventory {
        cam_yaw += xoffset
        cam_pitch += yoffset
        if cam_pitch > 1.57 {
            cam_pitch = 1.57
        }
        if cam_pitch < -1.57 {
            cam_pitch = -1.57
        }
    }
}

var AtlasLOD [][][3]uint32
func LoadAtlasLOD(filename string, scale int) {
    file, err := os.Open(filename)
    if err != nil {
        log.Println("INFO: Failed to open texture file:", filename)
        return
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        log.Println("INFO: Failed to decode texture:", filename)
        return
    }

    rgba := image.NewRGBA(img.Bounds())
    
    draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

    width := int32(rgba.Bounds().Dx())
    height := int32(rgba.Bounds().Dy())
    
    
    AtlasLOD = AtlasLOD[:0]
    for i := int32(0); i < height; i+=int32(scale) {
        var arr [][3]uint32
        for j := int32(0); j < width; j+=int32(scale) {
            r, g, b, _ := rgba.At(int(j),int(i)).RGBA()
            arr = append(arr, [3]uint32{r/255,g/255,b/255})
        }
        AtlasLOD = append(AtlasLOD, arr)
    }
    
}


func LoadTexture(filename string) uint32 {
    file, err := os.Open(filename)
    if err != nil {
        log.Println("INFO: Failed to open texture file:", filename)
        return 0
    }
    defer file.Close()

    img, err := png.Decode(file)
    if err != nil {
        log.Println("INFO: Failed to decode texture:", filename)
        return 0
    }

    rgba := image.NewRGBA(img.Bounds())
    
    

    draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

    width := int32(rgba.Bounds().Dx())
    height := int32(rgba.Bounds().Dy())

    var texture uint32
    gl.GenTextures(1, &texture)
    gl.BindTexture(gl.TEXTURE_2D, texture)

    gl.TexImage2D(
        gl.TEXTURE_2D,
        0,
        gl.RGBA,
        width,
        height,
        0,
        gl.RGBA,
        gl.UNSIGNED_BYTE,
        gl.Ptr(rgba.Pix),
    )

    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    log.Println("INFO: Texture loaded successfully", filename, texture)
    
    return texture
}


func updateSpread(x, y, z int) {
    var n = [7][3]int{ 
        {x  , y  , z  },
        {x+1, y  , z  },
        {x-1, y  , z  },
        {x  , y+1, z  },
        {x  , y-1, z  },
        {x  , y  , z+1},
        {x  , y  , z-1},
    }
    for j := 0; j < 7; j++ {
        blockUpdates = append(blockUpdates,[3]int{n[j][0],n[j][1],n[j][2]})
    }
}



var bub [][3]int
func processBlockUpdates() {
    bub = bub[:0]
    
    bub_updateSpread := func (x, y, z int) {
        var n = [7][3]int{ 
            {x  , y  , z  }, 
            {x+1, y  , z  }, 
            {x-1, y  , z  },
            {x  , y+1, z  },
            {x  , y-1, z  },
            {x  , y  , z+1},
            {x  , y  , z-1},
        }
        for j, _ := range n {
            bub = append(bub,[3]int{n[j][0],n[j][1],n[j][2]})
        }
    }

    
    for _, v := range blockUpdates {
        if  b := getBlockRef(v[0], v[1], v[2]); *b == WATER {
            
        } else if a := getBlockRef(v[0], v[1], v[2]); *a == AIR {
            var n = [5][3]int{ 
                {int(v[0]+1), int(v[1]  ), int(v[2]  )},
                {int(v[0]-1), int(v[1]  ), int(v[2]  )},
                {int(v[0]  ), int(v[1]+1), int(v[2]  )},
                {int(v[0]  ), int(v[1]  ), int(v[2]+1)},
                {int(v[0]  ), int(v[1]  ), int(v[2]-1)},
            }            
            var water_neighbors int = 0
            var lava_neighbors int = 0
            for k, _ := range n {
                b := getBlockVal(n[k][0], n[k][1], n[k][2]);
                if b == WATER {
                    water_neighbors++
                }
                if b == LAVA {
                    lava_neighbors++
                }
            }
            var water_can_move_here bool = false
            var lava_can_move_here bool = false
            switch water_neighbors {
            case 3:
                water_can_move_here = true
            default:
            }
            switch lava_neighbors {
            case 3:
                lava_can_move_here = true
            default:
            }
            
            if water_can_move_here && lava_can_move_here {
                *a = COBBLE
                bub_updateSpread(v[0], v[1], v[2])
            } else if water_can_move_here {
                *a = WATER
                bub_updateSpread(v[0], v[1], v[2])
            } else if lava_can_move_here {
                *a = LAVA
                bub_updateSpread(v[0], v[1], v[2])
            }
            
            
        } else if a := getBlockRef(v[0], v[1], v[2]); *a == SAND {
            b := getBlockRef(v[0], v[1]-1, v[2])
            if *b == AIR {
                *a = AIR
                *b = SAND
                bub = append(bub, [3]int{v[0], v[1]-1, v[2]})
                chunkUpdateDisplaylist( v[0], v[1], v[2])
            }
        } else if a := getBlockRef(v[0], v[1], v[2]); *a == GRAVEL {
            b := getBlockRef(v[0], v[1]-1, v[2])
            if *b == AIR {
                *a = AIR
                *b = GRAVEL
                bub = append(bub, [3]int{v[0], v[1]-1, v[2]})
                chunkUpdateDisplaylist( v[0], v[1], v[2])
            }
        } else if a := getBlockRef(v[0], v[1], v[2]); *a == MOON {
            b := getBlockRef(v[0], v[1]-1, v[2])
            if *b == AIR {
                *a = AIR
                *b = MOON
                bub = append(bub, [3]int{v[0], v[1]-1, v[2]})
                chunkUpdateDisplaylist( v[0], v[1], v[2])
            }
        }
    }
    blockUpdates = blockUpdates[:0]
    blockUpdates = append(blockUpdates, bub...)
}



func DrawText (str []byte, x, y, s float32) {
    d := func (id BlockID, offset, x, y, z float32) {
        texHPos := func(face uint32) float32 {
            return float32(face) * (16.0 / 1024.0)
        }
        textVPos := func(id BlockID) float32 {
            return float32(id) * (16.0 / 1024.0)
        }
        var v = [4][3]float32{{0.5, -0.5, 0.0}, {-0.5, -0.5, 0.0}, {-0.5, 0.5, 0.0}, { 0.5, 0.5, 0.0}}
        v0 := textVPos(id)
        v1 := textVPos(id + 1)
        h0 := texHPos(uint32(offset))
        h1 := texHPos(uint32(offset + 1))    
        gl.TexCoord2f(h1, v1)
        gl.Vertex3f(x+v[0][0]*s, y+v[0][1]*s, z+v[0][2])
        gl.TexCoord2f(h0, v1)
        gl.Vertex3f(x+v[1][0]*s, y+v[1][1]*s, z+v[1][2])
        gl.TexCoord2f(h0, v0)
        gl.Vertex3f(x+v[2][0]*s, y+v[2][1]*s, z+v[2][2])
        gl.TexCoord2f(h1, v0)
        gl.Vertex3f(x+v[3][0]*s, y+v[3][1]*s, z+v[3][2])
    }
    
    midpadh := (float32(len(str)-1)*s)/2.0
    
    for k, v := range str {
        if v == ' ' {
            continue
        } else if  (v >= '0') && (v <= '9') {
            //~ d(AIR, float32(FONT_0 + (v - '0')), 0, y, float32(k))
            d(AIR, float32(FONT_0 + (v - '0')), float32(k)*s, y, 0)
        } else if (v >= 'a') && (v <= 'z')  {
            //~ d(AIR, float32(FONT_A + (v - 'a')), 0, y, float32(k))
            d(AIR, float32(FONT_A + (v - 'a')), float32(k)*s - midpadh, y*s, 0)
        }
    }
}
    
var drop_amount int = 1
var player_textbuffer []byte

func DrawInventory() {
    gl.Disable(gl.DEPTH_TEST)
    gl.Disable(gl.CULL_FACE);
    var Names = [N_BLOCK_IDS][]byte{
        STONE:          []byte("stone"),
        GRASS:          []byte("grass"),
        DIRT:           []byte("dirt"),
        FLOWER:         []byte("flower"),
        SAND:           []byte("sand"),
        GRAVEL:         []byte("gravel"),
        CLAY:           []byte("clay"),
        MOON:           []byte("moon"),
        SNOW:           []byte("snow"),
        ORE_COAL:       []byte("coal"),
        ORE_IRON:       []byte("iron ore"),
        WOOD:           []byte("wood"),
        PLANK:          []byte("plank"),
        CRAFTING_TABLE: []byte("craft table"),
        PAINTING:       []byte("painting"),
        
        SHOVEL:         []byte("shovel"),
        PICKAXE:        []byte("pickaxe"),
        AXE:            []byte("axe"),
        
        STICK:          []byte("stick"),
        SAPLING:        []byte("sapling"),
        CANDLE:         []byte("candle"),
        WAX:            []byte("wax"),
        SIGN:           []byte("sign"),
        RAYTRACED_SPHERE: []byte("weird matter"),
    }
    
    var Emotions = [][]byte{
        []byte("happiness"),
        []byte("sadness"),
        []byte("anger"),
        []byte("fear"),
        []byte("relief"),
        []byte("humber"),
        []byte("nage"),
        []byte("dorceness"),
        []byte("andric"),
        []byte("varition"),
        []byte("ponnish"),
        []byte("trantiness"),
        []byte("teluge"),
        []byte("loric"),
        []byte("naritrance"),
    }
    /*
    var Effects = [][]byte{
        []byte("infection"),
        []byte("dread"),
        []byte("slumber"),
        []byte("petrification"),
        []byte("curse"),
        []byte("insanity"),
        []byte("blindness"),
    }
    */

    gl.MatrixMode(gl.PROJECTION)
    //~ gl.PushMatrix() 
    gl.LoadIdentity()
    
    aspect := float64(width) / float64(height)
    near := 0.1
    far := float64(SUBCHUNK_H*RENDER_DISTANCE)
    top := near * math.Tan(fov * PI / 360.0)
    bottom := -top
    right := top * aspect
    left := -right
    gl.Frustum(left, right, bottom, top, near, far)
    
    // Modelview
    gl.MatrixMode(gl.MODELVIEW)
    gl.PushMatrix() 
    gl.LoadIdentity()
    
    //~ gl.Viewport(0, 0, int32(widthInternal)/2, int32(heightInternal)/2)
    gl.Viewport(0, 0, int32(width/2),int32(height/2))
    LookAt(Vec3{0,0,40}, Vec3{0,0,0} )
    
    gl.Enable(gl.TEXTURE_2D);
    gl.BindTexture(gl.TEXTURE_2D, atlas)
    
    gl.Begin(gl.QUADS)
    
    gl.Color3f(1.0, 1.0, 1.0)
    
    
    DrawText([]byte(fmt.Sprintf("health %d",30)),0,  -1, 1)
    DrawText([]byte(fmt.Sprintf("drop amount %d",drop_amount)),0, -2, 1)
    DrawText([]byte(fmt.Sprintf("mood %s",Emotions[0])),0,  -3, 1)
    DrawText([]byte(fmt.Sprintf("speak %s",player_textbuffer)),0,  -4,1)
    
    
    for k, v := range player_inventory.ID {
        if k == selectedSlot {
            gl.Color3f(0.0, 1.0, 0.0)
        }
        if Names[v] != nil {
            DrawText([]byte(fmt.Sprintf("%-10s %d",Names[v], player_inventory.C[k])), 0, 1+float32(k),1)
        } else {
            DrawText([]byte(fmt.Sprintf("unknown %d",v)), 0, 1+float32(k),1)
        }
        
        if k == selectedSlot {
            gl.Color3f(1.0, 1.0, 1.0)
        }
    }
    
    
    gl.End()
    gl.Disable(gl.TEXTURE_2D);
    
    gl.MatrixMode(gl.PROJECTION)
    gl.PopMatrix()
    
    gl.MatrixMode(gl.MODELVIEW)
    gl.PopMatrix()
    
    gl.Enable(gl.DEPTH_TEST)
    gl.Enable(gl.CULL_FACE);
    
}


var getLightValL *[SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]uint8
var getLightValP iVec2
func getLightVal (x, y, z int) uint8 {
    if (getLightValP.x == (x-x%SUBCHUNK_H)) && (getLightValP.y == (z-z%SUBCHUNK_H)) {
        return getLightValL[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
    }
    c_p := iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}
    c := World[c_p]
    getLightValP = c_p
    getLightValL = &c.light
    return c.light[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
}

var setLightValL *[SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]uint8
var setLightValP iVec2
//~ var setLightValP uint32
func setLightVal (x, y, z int, v uint8) {
    if (setLightValP.x == (x-x%SUBCHUNK_H)) && (setLightValP.y == (z-z%SUBCHUNK_H)) {
    //~ if (setLightValP == psrng2d(x,z)) {
        setLightValL[x%SUBCHUNK_H][y][z%SUBCHUNK_H] = v
        return
    }
    c_p := iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}
    c := World[c_p]
    setLightValP = c_p
    //~ setLightValP = psrng2d(x,z)
    setLightValL = &c.light
    c.light[x%SUBCHUNK_H][y][z%SUBCHUNK_H] = v
}



var samples time.Duration
var total time.Duration

func BlockRandomUpdates(c *Chunk, p iVec2) {
    const BlockUpdatesPerFrame = 1
    
    if rand.Intn(60) == 0 {
        //~ for l := 0; l < BlockUpdatesPerFrame; l++ {
            i := rand.Intn(SUBCHUNK_H)
            j := 1+rand.Intn(CHUNK_V-2)
            k := rand.Intn(SUBCHUNK_H)
            b := &c.block[i][j][k]
            switch *b {
            case DIRT:
                if c.block[i][j+1][k] == AIR {
                    *b = GRASS
                    ChunkUpdateDisplayListRef(c);
                }
            case GRASS:
                if c.block[i][j+1][k] == AIR {
                    if len(Entities) < 30 {
                        
                        switch rand.Intn(2) {
                        case 0:
                            for i := 0; i < 6; i++ {
                                MobSpawn(p.x + i, j + 10, p.y + k, MOBID_DUCK)
                            }                            
                        case 1:
                            for i := 0; i < 6; i++ {
                                MobSpawn(p.x + i, j + 10, p.y + k, MOBID_BEE)
                            }
                        }
                    }
                }
                
            case OBSCURE:
                //~ if c.block[i][j+1][k] == AIR {
                    //~ *b = AIR;
                //~ } else if c.block[i][j-1][k] == AIR {
                    //~ *b = AIR;
                //~ } else if getBlockVal(p.x+i+1,j,p.y+k) == AIR {
                    //~ *b = AIR;
                //~ } else if getBlockVal(p.x+i-1,j,p.y+k) == AIR {
                    //~ *b = AIR;
                //~ } else if getBlockVal(p.x+i,j,p.y+k+1) == AIR {
                    //~ *b = AIR;
                //~ } else if getBlockVal(p.x+i,j,p.y+k-1) == AIR {
                    //~ *b = AIR;
                //~ }
                //~ ChunkUpdateDisplayListRef(c);
            case FLESH:
                if c.block[i][j+1][k] == AIR {
                    c.block[i][j+1][k] = FLESH
                    ChunkUpdateDisplayListRef(c);
                }
            case SAPLING:
                putTree(p.x+i,j,p.y+k)
            }
        //~ }
    }
}

func lightUpdate(c *Chunk, p iVec2, sunlight uint8) {
    
    const DEBUG = false
    var start time.Time
    
    if DEBUG { start = time.Now() }
    
    const overdraw = lightLevels
    var AirIndexes [SUBCHUNK_H+overdraw*2][SUBCHUNK_H+overdraw*2][]int
    var LightSources [][3]int
    
    for i := -overdraw; i < (SUBCHUNK_H+overdraw); i++ {
        for j := -overdraw; j < (SUBCHUNK_H+overdraw); j++ {
            hitopaque := false
            for k := CHUNK_V-1; k > 0; k-- {
                dx := p.x+i
                dy := k
                dz := p.y+j
                b := getBlockVal(dx,dy,dz)

                if (b == CANDLE) || (b == LILY) || (b == PIXIE) {
                    setLightVal(dx, dy, dz, 0)

                    if (i > 0) && (i < SUBCHUNK_H) && (j > 0) && (j < SUBCHUNK_H) {
                        LightSources = append(LightSources, [3]int{i, j, k})                        
                    }
                } else if !hitopaque {
                    setLightVal(dx, dy, dz, sunlight)
                    if (b != AIR) && (b != ICE) && (b != FLOWER) {
                        hitopaque = true
                    }
                } else {
                    setLightVal(dx, dy, dz, lightLevels-2)
                    if (b == AIR) || BlockIsLiquid[b] || (b == FLOWER) || (b == ICE) {
                        AirIndexes[i+overdraw][j+overdraw] = append(AirIndexes[i+overdraw][j+overdraw], k)
                    }
                }
            }
        }
    }
    
    
    /* add an lightlevels x lightlevels x lightlevels area to the 
     * calculated light values around artificial light sources 
    */
    for _, v := range LightSources {
        for i := -lightLevels+1; i < (lightLevels-1); i++ {
            for j := -lightLevels+1; j < (lightLevels-1); j++ {
                for k := -lightLevels+1; k < (lightLevels-1); k++ {
                    if (i == 0) && (j == 0) && (k == 0) {
                        continue
                    }
                    AirIndexes[v[0]+i+overdraw][v[1]+j+overdraw] = append(AirIndexes[v[0]+i+overdraw][v[1]+j+overdraw], v[2]+k)
                }
            }
        }
    }
    
    
    var mid time.Time
    if DEBUG {
        mid = time.Now()
        log.Println("First Step: ", time.Since(start) )
    }   
    
    Min6 := func (a, b, c, d, e, f, g uint8) uint8 {
        min := a
        if b < min { min = b }
        if c < min { min = c }
        if d < min { min = d }
        if e < min { min = e }
        if f < min { min = f }
        if g < min { min = g }
        return min
    }
    const passes = lightLevels
    for l := 0; l  < passes; l++ {
        for i := -overdraw; i < (SUBCHUNK_H+overdraw); i++ {
            for j := -overdraw; j < (SUBCHUNK_H+overdraw); j++ {
                dx := p.x + i
                dz := p.y + j
                for _, v := range AirIndexes[i+overdraw][j+overdraw] {
                    dy := v
                    setLightVal(dx,dy,dz, Min6(
                        getLightVal(dx+1, dy  , dz  ),
                        getLightVal(dx-1, dy  , dz  ),
                        getLightVal(dx  , dy+1, dz  ),
                        getLightVal(dx  , dy-1, dz  ),
                        getLightVal(dx  , dy  , dz+1),
                        getLightVal(dx  , dy  , dz-1),
                        lightLevels-1,
                    )+1)
                }
            }
        }
    }
    for i := -overdraw; i < (SUBCHUNK_H+overdraw); i++ {
        for j := -overdraw; j < (SUBCHUNK_H+overdraw); j++ {
            AirIndexes[i+overdraw][j+overdraw] = AirIndexes[i+overdraw][j+overdraw][:0]
        }
    }
    if DEBUG {
        log.Println("Second step:", time.Since(mid), p )
        log.Println()
    }
}




func putPillar(id BlockID, x, y, z int) {
    h := int(rand.Uint32() % 10)
    for i := 0; i < h; i++ {
        if (y + h) > (CHUNK_V - 1) {
            break
        }
        b := getBlockRef(x, y + i, z)
        *b = id
    }
}


func putTree(x, y, z int) {
    var tree0 = [5][7][5]int{
        {   {0,0,1,0,0},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {0,0,1,0,0},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {1,1,1,1,1},
            {1,1,1,1,1},
            {1,1,2,1,1},
            {1,1,2,1,1},
            {0,0,2,0,0},
            {0,0,2,0,0},
            {0,0,2,0,0},
        },
        {   {0,0,1,0,0},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {0,0,1,0,0},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {1,1,1,1,1},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
    }
    
    var tree1 = [5][9][5]int{
        {   {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,1,1,1,0},
            {0,0,1,0,0},
            {0,1,1,1,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {0,0,1,0,0},
            {0,1,2,1,0},
            {0,0,2,0,0},
            {0,1,2,1,0},
            {1,1,2,1,1},
            {0,1,2,1,0},
            {1,1,2,1,1},
            {0,0,2,0,0},
            {0,0,2,0,0},
        },
        {   {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,1,1,1,0},
            {0,0,1,0,0},
            {0,1,1,1,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
        {   {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,1,0,0},
            {0,0,0,0,0},
            {0,0,0,0,0},
        },
    }
    
    
    if rand.Uint32() % 2 == 0 {
        for i := 0; i < 5; i++ {
            for j := 0; j < 7; j++ {
                for k := 0; k < 5; k++ {
                    if (y+7-j) >= CHUNK_V {
                        continue
                    }
                    if tree0[i][j][k] == 1 {
                        *getBlockRef(x+k-2, y+7-j, z+i-2) = LEAF
                    } else if tree0[i][j][k] == 2 {
                        *getBlockRef(x+k-2, y+7-j, z+i-2) = WOOD
                    }
                }
            }
        }
    } else {
        for i, _ := range tree1 {
            for j, _ := range tree1[i] {
                for k, _ := range tree1[i][j] {
                    if (y+7-j + 1) >= CHUNK_V {
                        continue
                    }
                    if tree1[i][j][k] == 1 {
                        *getBlockRef(x+k-2, y+7-j + 1, z+i-2) = LEAF
                    } else if tree1[i][j][k] == 2 {
                        *getBlockRef(x+k-2, y+7-j + 1, z+i-2) = WOOD
                    }
                }
            }
        }
    }
}


var scheduleTree [][3]int
var dirt_blocks []int
func terrain(p iVec2) {
    
    const DEBUG = false
    var start time.Time
    if DEBUG { start = time.Now() }
    c := &Chunk{}
    lerp := func (from, to, t float32) float32 {
        return from + t*(to-from)
    }
    Norm := func (f float32) float32 {
        f -= 0.5
        f *= 2.0
        f *= f
        return f
    }
    absAndOffset := func (f float32) float32 {
        f = Norm(f)
        //~ f *= 6.0
        f *= 1.5
        if f < 0.0 {
            f = -f
        }
        if f > 1.0 {
            f = 0.99
        }
        return f
    }
    
    
    c.display_list = 0
    World[p] = c

    filename := fmt.Sprintf("%s/%d,%d.bin",saveSlot,p.x, p.y)
    if _, err := os.Stat(filename); err == nil {
        saveReadChunk(c, p)
        return
    }


    /* the entire codebase assumes positive world coordinates */
    if (p.x < 1000) || (p.y < 1000) {
        for i := 0;  i < SUBCHUNK_H; i++ {
            for j := 0;  j < CHUNK_V; j++ {
                for k := 0;  k < SUBCHUNK_H; k++ {
                    c.block[i][j][k] = AIR
                }
            }
        }
        return
    }


    /* brick pillar */
    if (psrng2d(p.x, p.y) % 2000) == 0 {        
        for i := 0;  i < SUBCHUNK_H; i++ {
            for j := 0;  j < CHUNK_V; j++ {
                for k := 0;  k < SUBCHUNK_H; k++ {
                    if j > 3 {
                        c.block[i][j][k] = BRICK
                    } else {
                        c.block[i][j][k] = BEDROCK
                    }
                }
            }
        }
        return
    } 
    
    vegetationNoise := func (nx, nz int) (uint32, bool) {
        var foo uint32
        var bar bool = false
        const folliage_scale = 1000.0
        n := noise2d((float32(nx-8322))/folliage_scale, float32(nz+123)/folliage_scale)
        n -= 0.5
        n *= 6.0
        if n < 0.0 {
            n = -n
        }
        if n > 1.0 {
            n = 1.0
        }
        
        foo = uint32(256 - 200*n)
        if foo > 250 {
            bar = true
        }
        return foo, bar
    }
    
    /* true == snow, false == no snow */
    temperatureNoise := func (x, z int) (bool) {
        var scale float32 = 2000.0
        n := noise2d((float32(x+8172))/scale, float32(z-12323)/scale)
        //~ if n > 0.673 {
        if n > 0.875 {
            return true
        }
        return false
    }
    const offset = 20.0
    
    
    seaLevelVariation := func(x, z int) (int, float32) {
        var scale float32 = 1573.0
        var n float32 = fractalNoise2D(float32(x+97898)/scale, float32(z-21323)/scale, 3, 4, 0.25)
        var variation float32 = offset
        var continentality float32 = 0
        
        var fract float32 = (n*variation) - float32(int((n*variation)))
        //~ if fract > 0.5 {
            //~ fract -= 0.75
            //~ if fract < 0 {
                //~ fract = -fract
            //~ }
            //~ fract *= 4.0
            
            //~ continentality = (10*(fract))
        //~ }
        
        fract -= 0.5
        if fract < 0 {
            fract = -fract
        }
        fract *= 2.0
        //~ continentality = (offset*(fract))
        continentality = (n * variation)*(fract)
        
        return int(n * variation), continentality
    }
    

    var seaBottom int = (CHUNK_V/2) - offset*2
    
    
    var seaLevel int
    
    for i := 0; i < SUBCHUNK_H; i++ {
        for k := 0; k < SUBCHUNK_H; k++ {
            
            var erosion bool = noise2d(float32(p.x+i)/10.0,float32(p.y+k)/10.0) > 0.80
            
            var nx = p.x + i
            var nz = p.y + k
            
            var continentality int
            {
                n, c := seaLevelVariation(nx, nz)
                seaLevel = (CHUNK_V/2) + n - offset*2
                continentality = int(c)
            }
            
            tree_distribution, wasteland := vegetationNoise(nx, nz)
            snow := temperatureNoise(nx, nz)
            dirt_blocks = dirt_blocks[:0]
            
            
            v_scale_variation := 0.75 + noise2d(float32(p.x+i+12312)/1000.0,float32(p.y+k)/1000.0)*0.25
            
            for j := 0; j < (CHUNK_V-5); j++ {
                
                placeStone := func() {
                    cavescale := float32(50)
                    n := noise3d(float32(nx+1333)/cavescale, float32(j)/(cavescale/3.0), float32(nz-1888)/cavescale)
                    nn := float32(0.02)
                    if (n > (0.5 - nn)) && (n < (0.5 + nn)) && (j > 3) {
                        c.block[i][j][k] = AIR
                        return
                    }
                    
                    
                    const (
                        orescaleh = 12.0
                        orescalev =  6.0
                        seed = 12312
                    )
                    //~ if j == 6 {
                        //~ c.block[i][j][k] = LAVA
                        //~ return
                    //~ }
                    if j < 5 {
                        c.block[i][j][k] = BEDROCK
                        return
                    }
                    if (psrng3d(i,j,k) % 8) == 0 {
                        if Norm(noise3d((float32(nx)+seed)/orescaleh, float32(j)/orescalev, float32(nz)/orescaleh)) > 0.05 {
                            c.block[i][j][k] = ORE_IRON
                            return
                        }
                        if Norm(noise3d((float32(nx)+(seed*2))/orescaleh, float32(j)/orescalev, float32(nz)/orescaleh)) > 0.05 {
                            c.block[i][j][k] = ORE_COAL
                            return
                        }
                    }
                    c.block[i][j][k] = STONE
                }
                
                if j < (seaBottom) {
                    placeStone()
                    continue
                }
                
                var hscale float32 = 256.0*v_scale_variation //default: 68
                const order = 3
                                
                var lownoise float32
                var highnoise float32
                
                var selectnoise = fractalNoise2D( float32(nx-10000.0)/600.0, float32(nz+10000.0)/600.0, 3, 4.0, 0.25 )
                //~ var selectnoise = fractalNoise3D( float32(nx-10000.0)/600.0, float32(j)/600.0, float32(nz+10000.0)/600.0, 2, 4.0, 0.25 )
                var h float32;
                
                //~ lownoise = absAndOffset( fractalNoise3D(float32(nx)/hscale, float32(j)/(hscale/2.0), float32(nz)/hscale, order,  4, 0.25) )
                //~ highnoise = absAndOffset( fractalNoise3D(float32(nx+10000.0)/hscale, float32(j)/(hscale/2.0), float32(nz+10000.0)/hscale, order, 4, 0.25) )
                lownoise  = absAndOffset( fractalNoise2D(float32(nx)/hscale, float32(nz)/hscale, order,  4, 0.25) )
                highnoise = absAndOffset( fractalNoise2D(float32(nx+10000.0)/hscale, float32(nz+10000.0)/hscale, order, 4, 0.25) )
                
                //~ if selectnoise < 0.33 {
                    //~ h = float32(continentality) + float32(seaBottom) + lownoise*(CHUNK_V/2);
                //~ } else if selectnoise < 0.66 {
                    h = float32(continentality) + float32(seaBottom) + lerp(highnoise, lownoise, 1.0 - selectnoise)*(CHUNK_V/2);
                //~ } else {
                    //~ h = float32(continentality) + float32(seaBottom) + highnoise*(CHUNK_V/2);
                //~ }
                
                if j > (seaLevel+3) {
                    overhangscale := float32(30.0)
                    overhangnoise := 2*fractalNoise3D(float32(nx)/overhangscale, float32(j)/overhangscale, float32(nz)/overhangscale, 2,  4, 0.25)
                    h *= float32(overhangnoise)
                }
                
                
                
                if (int(j) < int(h)) {
                    placeStone()
                //~ } else if (int(j) == int(h)) {
                } else if !erosion && ((int(j+1) >= int(h)) && (int(j-1) <= int(h))) {
                        /* skip on an erosion patch */
                    if  (j > (seaLevel+1)) {
                        dirt_blocks = append(dirt_blocks, j)
                        c.block[i][j][k] = DIRT;
                    } else {
                        if selectnoise < 0.354 { //~ 0.000 ... 0.354
                            c.block[i][j][k] = SAND
                        } else if selectnoise < 0.500 { //~ 0.354 ... 0.500
                            c.block[i][j][k] = GRAVEL
                        } else if selectnoise < 0.646 { //~ 0.500 ... 0.646
                            c.block[i][j][k] = MOON
                        } else { //~ 0.646 ... 1.000
                            c.block[i][j][k] = CLAY
                        }
                    }
                } else {
                    if (j == seaLevel) && snow {
                        c.block[i][j][k] = ICE;
                    } else if j <= seaLevel {
                        //~ c.block[i][j][k] = LAVA; //testing lava 
                        c.block[i][j][k] = WATER; 
                    }else{
                        if (c.block[i][j-1][k] == WATER) && ((rand.Uint32() % 1024) == 0) {
                            c.block[i][j][k] = LILY
                        } else {
                            c.block[i][j][k] = AIR;
                        }
                    }
                }
            }
            
            for _, v := range dirt_blocks {
                if c.block[i][v+1][k] == AIR {
                    if snow {
                        c.block[i][v+1][k] = SNOW;
                    } 
                    c.block[i][v][k] = GRASS;
                }
                if (i > 4) &&
                (i < (SUBCHUNK_H-4)) &&
                (k > 4) &&
                (k < (SUBCHUNK_H-4)) {

                    r := psrng3d(nx, v, nz)
                    
                    if wasteland {
                        switch (r % (tree_distribution/2)) {
                        case 0:
                            *getBlockRef(nx, v+1, nz) = FLOWER
                        case 1:
                            putPillar(FLESH,nx, v+1, nz) 
                        default:
                        }
                    } else {
                        switch (r % tree_distribution) {
                        case 0:
                            *getBlockRef(nx, v+1, nz) = FLOWER
                        case 1:
                            scheduleTree = append(scheduleTree, [3]int{nx, v, nz})
                        default:
                        }
                    }
                    
                    if (r%1777) == 0 {
                        c.block[i][v+1][k] = RAYTRACED_SPHERE
                    }                    
                }
            }
            if (rand.Uint32() % 1024) == 0 {
                c.block[i][CHUNK_V-15- rand.Uint32()%10 ][k] = PIXIE
            }
            
        }
    }
    
    var mid time.Time
    if DEBUG {
        mid = time.Now()
        log.Println("First Step: ", time.Since(start) )
    }   
    
    for i := 0; i < SUBCHUNK_H; i++ {
        for j := 0; j < SUBCHUNK_H; j++ {
            var nx = p.x + i
            var nz = p.y + j
            n := noise2d(float32(nx)/10.0, float32(nz)/10.0)
            /*
            if n > 0.60 {
                c.cloud[i][j] = true
            } else {
                c.cloud[i][j] = false
            }
            */
            c.cloud[i][j] = byte(255.0*n)
            
        }
    }


    /* catacomb generation */
    {
        scale := float32(100.0)
        if noise2d((float32(p.x+12323))/scale, float32(p.y-12323)/scale) > 0.75 {
            
            h := 16
            y_level := (CHUNK_V/2) - 40 - h
            
            for i := 1;  i < (SUBCHUNK_H-1); i++ {
                for j := 0;  j < h; j++ {
                    for k := 1;  k < (SUBCHUNK_H-1); k++ {
                        
                        c.block[i][j+y_level][k] = AIR
                    }
                }
            }
            c.block[SUBCHUNK_H/2][y_level  ][0] = AIR
            c.block[SUBCHUNK_H/2][y_level+1][0] = AIR
            
            c.block[SUBCHUNK_H/2][y_level][SUBCHUNK_H-1] = AIR
            c.block[SUBCHUNK_H/2][y_level+1][SUBCHUNK_H-1] = AIR
            
            c.block[0][y_level  ][SUBCHUNK_H/2] = AIR
            c.block[0][y_level+1][SUBCHUNK_H/2] = AIR

            c.block[SUBCHUNK_H-1][y_level  ][SUBCHUNK_H/2] = AIR
            c.block[SUBCHUNK_H-1][y_level+1][SUBCHUNK_H/2] = AIR

            
        }
    }



    for _, v := range scheduleTree {
        
        trunk := int(rand.Uint32()) %10
        for i := 0; i <= trunk; i++ {
            *getBlockRef(v[0], v[1]+i, v[2]) = WOOD
        }
        
        putTree(v[0], v[1]+int(trunk), v[2])

    }
    scheduleTree = scheduleTree[:0]
    
    if DEBUG {
        log.Println("Second step:", time.Since(mid), p )
        log.Println()
    }
}


func psrng3d(x, y, z int) uint32 {
	//~ h := uint32(x)*374761393 + uint32(y)*668265263 + uint32(z)*1446645773
	//~ h = (h ^ (h >> 13)) * 1274126177
	//~ return h ^ (h >> 16)
    
    return uint32(x)*374761393 + uint32(y)*668265263 + uint32(z)*1446645773
}

func psrng2d(x, y int) uint32 {
    //~ h := uint32(x)*374761393 + uint32(y)*374761393
    //~ h = (h ^ (h >> 13)) * 1274126177
    //~ return h ^ (h >> 16)
    
    return uint32(x)*374761393 + uint32(y)*668265263
}

func psrng1d(x int) uint32 {
	//~ h := uint32(x) * 374761393
	//~ h = (h ^ (h >> 13)) * 1274126177
	//~ return h ^ (h >> 16)
    
    //~ return uint32(x*374761393 ^ (x>>13))

	h := int32(x) * int32(0x79B9)
	return uint32(h ^ (h >> 16))
}



var terrain_newseed_state uint32 = 0
func terrain_newseed() uint32 {
    terrain_newseed_state = psrng1d(int(terrain_newseed_state))
    return terrain_newseed_state
}

    
func noise2d(x, y float32) float32 {
    floor := func(f float32) int {
        i := int(f)
        if f < 0 && (float32(i) != f) {
            return i - 1
        }
        return i
    }
	hash := func(x, y int) uint32 {
		h := uint32(x)*374761393 + uint32(y)*668265263
		h = (h ^ (h >> 13)) * 1274126177
		return h ^ (h >> 16)
	}
	xi := floor(x)
	yi := floor(y)
	xf := x - float32(xi)
	yf := y - float32(yi)
	v00 := float32(hash(xi, yi)) / 4294967296.0
	v10 := float32(hash(xi+1, yi)) / 4294967296.0
	v01 := float32(hash(xi, yi+1)) / 4294967296.0
	v11 := float32(hash(xi+1, yi+1)) / 4294967296.0
	u := xf * xf * (3 - 2*xf)
	v := yf * yf * (3 - 2*yf)
	x1 := v00 + u*(v10-v00)
	x2 := v01 + u*(v11-v01)
	return x1 + v*(x2-x1)
}


func noise3d(x, y, z float32) float32 {
    const norm = 1.0 / 4294967296.0

    xi := int(x)
    yi := int(y)
    zi := int(z)

    xf := x - float32(xi)
    yf := y - float32(yi)
    zf := z - float32(zi)

    // Keep hash returning uint32
    v000 := float32(psrng3d(xi,   yi,   zi  )) * norm
    v100 := float32(psrng3d(xi+1, yi,   zi  )) * norm
    v010 := float32(psrng3d(xi,   yi+1, zi  )) * norm
    v110 := float32(psrng3d(xi+1, yi+1, zi  )) * norm
    v001 := float32(psrng3d(xi,   yi,   zi+1)) * norm
    v101 := float32(psrng3d(xi+1, yi,   zi+1)) * norm
    v011 := float32(psrng3d(xi,   yi+1, zi+1)) * norm
    v111 := float32(psrng3d(xi+1, yi+1, zi+1)) * norm

    // Smoothstep
    u := xf * xf * (3 - 2*xf)
    v := yf * yf * (3 - 2*yf)
    w := zf * zf * (3 - 2*zf)

    x00 := v000 + u*(v100 - v000)
    x10 := v010 + u*(v110 - v010)
    x01 := v001 + u*(v101 - v001)
    x11 := v011 + u*(v111 - v011)

    y0 := x00 + v*(x10 - x00)
    y1 := x01 + v*(x11 - x01)

    return y0 + w*(y1 - y0)
}

func fractalNoise2D(x, y float32, octaves int, lacunarity, persistence float32) float32 {
	var total float32
	var amplitude float32 = 1.0
	var frequency float32 = 1.0
	var maxAmplitude float32 = 0.0

	for i := 0; i < octaves; i++ {
		total += noise2d(x*frequency, y*frequency) * amplitude
		maxAmplitude += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}

	return total / maxAmplitude
}

func fractalNoise3D(x, y, z float32, octaves int, lacunarity, persistence float32) float32 {
	var total float32
	var amplitude float32 = 1.0
	var frequency float32 = 1.0
	var maxAmplitude float32 = 0.0
	for i := 0; i < octaves; i++ {
		total += noise3d(x*frequency, y*frequency, z*frequency) * amplitude
		maxAmplitude += amplitude
		amplitude *= persistence
		frequency *= lacunarity
	}
    return total / maxAmplitude
}


var dir Vec3
func setupScene(window *glfw.Window, w, h int32) {
    AddVec3 := func ( a, b Vec3) Vec3 {
        return Vec3{
            x: a.x+b.x, 
            y: a.y+b.y, 
            z: a.z+b.z,
        }
    }
    aspect := float64(w) / float64(h)
    gl.Viewport(0, 0, int32(w), int32(h))
    gl.Enable(gl.DEPTH_TEST)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    near := 0.1
    far := float64(SUBCHUNK_H*RENDER_DISTANCE)
    top := near * math.Tan(fov * PI / 360.0)
    bottom := -top
    right := top * aspect
    left := -right
    gl.Frustum(left, right, bottom, top, near, far)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    LookAt(Vec3{player_pos.x, player_pos.y+ sinf(walking_bob)*bob_range, player_pos.z}, AddVec3(player_pos, dir))
}

var light uint8 = 0


var n000, n001, n010, n011, n100, n101, n110, n111 float32
var nBase Vec3 
func noise3d_fast_set(x, y, z, h, v float32) {
    
    n000 = noise3d(x,       y,       z)
    n001 = noise3d(x,       y,       z + h)
    n010 = noise3d(x,       y + v,   z)
    n011 = noise3d(x,       y + v,   z + h)
    
    n100 = noise3d(x + h,   y,       z)
    n101 = noise3d(x + h,   y,       z + h)
    n110 = noise3d(x + h,   y + v,   z)
    n111 = noise3d(x + h,   y + v,   z + h)
}

/* x y z are values between 0 and 1 */
func noise3d_fast_get(x, y, z float32) float32 {
    nx0 := n000 + (n100-n000)*x
    nx1 := n001 + (n101-n001)*x
    nx2 := n010 + (n110-n010)*x
    nx3 := n011 + (n111-n011)*x

    ny0 := nx0 + (nx2-nx0)*y
    ny1 := nx1 + (nx3-nx1)*y

    return ny0 + (ny1-ny0)*z
}
    

const (
    WEATHER_NORMAL = iota
    WEATHER_CLEAR
    WEATHER_FOGGY
    WEATHER_RAINY
    WEATHER_EERIE
    WEATHER__N
)
var Forecast uint32

var drawMeLater []struct{
    id BlockID
    x, y, z int }
var chunk_render_order []struct{
    i, j int
    d float32 }
    
var lastLightUpdate uint32


var schedule_raytraced_blocks_p int
var schedule_raytraced_blocks [256*4]float32


func DrawBlock ( id BlockID, x, y, z int) {
    const MAGIC = 0xFF;
    
    var occlusion bool = false
    var underwater bool = false
    var underlava bool = false
    var color [4]float32
    var v [][4][3]float32

    switch id {
    case FLOWER:
        occlusion = false
    default:
        occlusion = true
    }

    texHPos := func ( face uint32) float32 {
        return float32(face)*(16.0/1024.0)
    }
    texHPosInv := func ( face uint32) float32 {
        return 1.0 - float32(face)*(16.0/1024.0)
    }
    textVPos := func ( id BlockID ) float32 {
        return float32(id)*(16.0/1024.0)
    }
    textVPosInv := func ( id BlockID ) float32 {
        return 1.0 - float32(id)*(16.0/1024.0)
    }

    biome_noise_function := func(x, z int) float32 {
        const factor = 2.5
        const scale = 800.0
        
        foo := noise2d(float32(x-x%(SUBCHUNK_H/2))/scale, float32(z-z%(SUBCHUNK_H/2)+10230)/scale)
        
        foo -= 0.5
        foo *= factor
        foo += 0.5
        if foo < 0.0 {
            foo = 0
        }
        if foo > 1.0 {
            foo = 1
        }
        return foo
    }

    lerp := func (from, to, t float32) float32 {
        return from + t*(to-from)
    }

    lerp4 := func (a, b, c, d, n float32) float32 {
        if n < 0.33333334 {
            return lerp(a, b, n*3) /* [0] to [1/3] */
        }
        if n < 0.6666667 {
            return lerp(b, c, (n-1.0/3.0)*3) /* [1/3] to [2/3] */
        }
        return lerp(c, d, (n-2.0/3.0)*3) /* [2/3] to [1] */
    }

    biome_noise := biome_noise_function(x, z)
    c := psrng3d(x, y, z)
    var r, g, b float32
    if (c&1) != 0 { r = 255.0/255.0 } else { r = 224.0/255.0 }
    if (c&2) != 0 { g = 255.0/255.0 } else { g = 224.0/255.0 }
    if (c&4) != 0 { b = 255.0/255.0 } else { b = 224.0/255.0 }

    occlusion_check := func (nb BlockID) bool {
        underwater = (nb == WATER)
        underlava = (nb == LAVA)
        switch nb {
        case ICE, AIR, LAVA, WATER, RAYTRACED_SPHERE:
            return false
        default:
            if nb > ITEM_0 {
                return false
            }
            return true
        }
    }
    
    face_color := func (nb BlockID) {
        switch id {
        case GRASS:
            color = [4]float32{
                //grass:    green          teal       purple       yellow
                lerp4( 85.0/255.0,  85.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise),
                lerp4(255.0/255.0, 235.0/255.0,  85.0/255.0, 255.0/255.0, biome_noise),
                lerp4( 85.0/255.0, 235.0/255.0,  85.0/255.0,  32.0/255.0, biome_noise), 1 }
        case LEAF:
            color = [4]float32{
                lerp4( 00.0/255.0, 255.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise),
                lerp4(255.0/255.0,  85.0/255.0,  85.0/255.0, 235.0/255.0, biome_noise),
                lerp4( 00.0/255.0,  85.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise), 1 }
        case WOOD:
            color = [4]float32{
                lerp4(168.0/255.0, 255.0/255.0, 168.0/255.0, 168.0/255.0, biome_noise),
                lerp4( 84.0/255.0, 255.0/255.0,  86.0/255.0,  84.0/255.0, biome_noise),
                lerp4(  0.0/255.0, 255.0/255.0,   0.0/255.0,   0.0/255.0, biome_noise), 1 }
                
        case STONE, COBBLE, ORE_COAL, ORE_IRON:
            color = [4]float32{r, g, b, 1}
        default:
        }
        
        if underwater {
            color[0] *= 0.5
            color[1] *= 0.5
            color[3] = 0.5
        } else if underlava {
            color[0] *= 2.0
            color[0] += 0.2
            color[2] *= 0.1
            color[3] = 0.1
        }
    }
    
    var n = [6][3]int{
        { 1, 0, 0},
        {-1, 0, 0},
        { 0, 1, 0},
        { 0,-1, 0},
        { 0, 0, 1},
        { 0, 0,-1},
    }


    if id == SIGN {
        s := float32((7.0/8.0)/2.0)
        v = [][4][3]float32{
            {{0, 0.5, 0.5}, {1, 0.5, 0.5}, {1, 1, 0.5}, {0, 1, 0.5}}, // +Z
            {{1, 0.5, 0.5}, {0, 0.5, 0.5}, {0, 1, 0.5}, {1, 1, 0.5}}, // -Z

            {{0+s, 0.0, 0.5}, {1-s, 0.0, 0.5}, {1-s, 0.5, 0.5}, {0+s, 0.5, 0.5}}, // +Z
            {{1-s, 0.0, 0.5}, {0+s, 0.0, 0.5}, {0+s, 0.5, 0.5}, {1-s, 0.5, 0.5}}, // -Z
            
        }
    } else if id == PAINTING {
        
        v = [][4][3]float32{
            {{1, 0, 1}, {1, 0, 0}, {1, 1, 0}, {1, 1, 1}}, // +X
            {{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}}, // -X
            {{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}}, // +Y
            {{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // -Y
            {{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, // +Z
            {{1, 0, 0}, {0, 0, 0}, {0, 1, 0}, {1, 1, 0}}, // -Z
        }
    } else if id >= ITEM_0 {
        for f := uint32(0); f < 6; f++ {
            var spacing float32 = float32(f)/24.0
            v = append(v, [4][3]float32{ 
                {0.0, 0.01 + spacing, 0.0},
                {0.0, 0.01 + spacing, 1.0},
                {1.0, 0.01 + spacing, 1.0},
                {1.0, 0.01 + spacing, 0.0},})
        }
    } else if (id == WATER) || (id == LAVA) {
        height := float32(0.875)
        v = [][4][3]float32{
            {{1,     0, 1}, {1,     0, 0}, {1, height, 0}, {1, height, 1}}, // +X
            {{0,     0, 0}, {0,     0, 1}, {0, height, 1}, {0, height, 0}}, // -X
            {{0,height, 0}, {0,height, 1}, {1, height, 1}, {1, height, 0}}, // +Y
            {{0,     0, 0}, {1,     0, 0}, {1,      0, 1}, {0,      0, 1}}, // -Y
            {{0,     0, 1}, {1,     0, 1}, {1, height, 1}, {0, height, 1}}, // +Z
            {{1,     0, 0}, {0,     0, 0}, {0, height, 0}, {1, height, 0}}, // -Z
        }
    } else {
        v = [][4][3]float32{
            {{1, 0, 1}, {1, 0, 0}, {1, 1, 0}, {1, 1, 1}}, // +X
            {{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}}, // -X
            {{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}}, // +Y
            {{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // -Y
            {{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, // +Z
            {{1, 0, 0}, {0, 0, 0}, {0, 1, 0}, {1, 1, 0}}, // -Z
        }
    }

    for f := uint32(0); f < uint32(len(v)); f++ {

        var nb BlockID
        var lv float32
        color = [4]float32{1, 1, 1, 1}
        if id < ITEM_0 {
            nb = getBlockVal(x + n[f][0], y + n[f][1],z + n[f][2])
            if occlusion && occlusion_check(nb) {
                continue
            }
            face_color(nb)
            lv = float32(lightLevels - getLightVal(x + n[f][0], y + n[f][1],z + n[f][2]))/float32(lightLevels)
        } else {
            lv = float32(lightLevels - getLightVal(x, y,z))/float32(lightLevels)            
        }
        
        gl.Color4f(color[0]*lv, color[1]*lv, color[2]*lv, color[3])

        var v0, v1, h0, h1 float32 
        if id >= ITEM_0 {
            v0, v1 = textVPosInv((id - ITEM_0)+0), textVPosInv((id - ITEM_0)+1)
            h0, h1 = texHPosInv(0), texHPosInv(1)
        } else {
            v0 = textVPos(id  )
            v1 = textVPos(id+1)
            h0 = texHPos(  f)
            h1 = texHPos(f+1)
        }
        gl.TexCoord2f(h1, v1);
        gl.Vertex3f(float32(x)+v[f][0][0], float32(y)+v[f][0][1], float32(z)+v[f][0][2])
        gl.TexCoord2f(h0, v1);
        gl.Vertex3f(float32(x)+v[f][1][0], float32(y)+v[f][1][1], float32(z)+v[f][1][2])
        gl.TexCoord2f(h0, v0);
        gl.Vertex3f(float32(x)+v[f][2][0], float32(y)+v[f][2][1], float32(z)+v[f][2][2])
        gl.TexCoord2f(h1, v0);
        gl.Vertex3f(float32(x)+v[f][3][0], float32(y)+v[f][3][1], float32(z)+v[f][3][2])
    }
}

func drawScene() {
    
    biome_noise_function := func(x, z int) float32 {
        const factor = 2.5
        const scale = 800.0
        
        foo := noise2d(float32(x-x%(SUBCHUNK_H/2))/scale, float32(z-z%(SUBCHUNK_H/2)+10230)/scale)
        
        foo -= 0.5
        foo *= factor
        foo += 0.5
        if foo < 0.0 {
            foo = 0
        }
        if foo > 1.0 {
            foo = 1
        }
        return foo
    }
    
    createDisplayLists := func (v *Chunk, p iVec2) {
       lerp := func (from, to, t float32) float32 {
            return from + t*(to-from)
        }

        lerp4 := func (a, b, c, d, n float32) float32 {
            if n < 0.33333334 {
                return lerp(a, b, n*3) /* [0] to [1/3] */
            }
            if n < 0.6666667 {
                return lerp(b, c, (n-1.0/3.0)*3) /* [1/3] to [2/3] */
            }
            return lerp(c, d, (n-2.0/3.0)*3) /* [2/3] to [1] */
        }
        
        texHPos := func ( face uint32) float32 {
            return float32(face)*(16.0/1024.0)
        }
        textVPos := func ( id BlockID ) float32 {
            return float32(id)*(16.0/1024.0)
        }
        
        //~ var drawPainting = func( id BlockID, x, y, z int) {
            //~ const MAGIC = 0xFF;     /*controls the direction the painting is facing*/
            
            //~ var n = [6][3]int{
                //~ { 1, 0, 0},
                //~ {-1, 0, 0},
                //~ { 0, 1, 0},
                //~ { 0,-1, 0},
                //~ { 0, 0, 1},
                //~ { 0, 0,-1},
            //~ }
            
            //~ var v = [6][4][3]float32{
                //~ {{1, 0, 1}, {1, 0, 0}, {1, 1, 0}, {1, 1, 1}}, // +X
                //~ {{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}}, // -X
                //~ {{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}}, // +Y
                //~ {{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // -Y
                //~ {{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, // +Z
                //~ {{1, 0, 0}, {0, 0, 0}, {0, 1, 0}, {1, 1, 0}}, // -Z
            //~ }
            
            //~ f := uint32(MAGIC);
            //~ for ff := uint32(0); ff < 6; ff++ {
                //~ nb := getBlockVal(x + n[ff][0], y + n[ff][1],z + n[ff][2])
                //~ if (nb == AIR) && (nb == WATER) {
                    //~ f = ff;
                    //~ break;
                //~ }
            //~ }
            //~ if f == MAGIC {
                //~ return
            //~ }

            //~ lv := float32(lightLevels - getLightVal(x + n[f][0], y + n[f][1],z + n[f][2]))/float32(lightLevels)

            //~ gl.Color3f(lv, lv, lv)
            //~ v0 := textVPos(id  )
            //~ v1 := textVPos(id+1)
            //~ h0 := texHPos(  f)
            //~ h1 := texHPos(f+1)
            //~ gl.TexCoord2f(h1, v1);

            //~ gl.Vertex3f(float32(x)+v[f][0][0], float32(y)+v[f][0][1], float32(z)+v[f][0][2])
            //~ gl.TexCoord2f(h0, v1);
            //~ gl.Vertex3f(float32(x)+v[f][1][0], float32(y)+v[f][1][1], float32(z)+v[f][1][2])
            //~ gl.TexCoord2f(h0, v0);
            //~ gl.Vertex3f(float32(x)+v[f][2][0], float32(y)+v[f][2][1], float32(z)+v[f][2][2])
            //~ gl.TexCoord2f(h1, v0);
            //~ gl.Vertex3f(float32(x)+v[f][3][0], float32(y)+v[f][3][1], float32(z)+v[f][3][2])
        //~ }
        
        v.uniform_index = 0
        
        if v.display_list == 0 {
            
            //~ start := time.Now() // capture start
            list := gl.GenLists(1)

            gl.NewList(list, gl.COMPILE);
            gl.Enable(gl.TEXTURE_2D);
            gl.BindTexture(gl.TEXTURE_2D, atlas)
            gl.Begin(gl.QUADS)

            for i := 0; i < len(v.uniform_array); i++ {
                v.uniform_array[i] = 0.0
            }
            
            var drawmelater [][6]int

            for i := 0; i < SUBCHUNK_H; i++ {
                for j := 1; j < (CHUNK_V-1); j++ {
                    for k := 0; k < SUBCHUNK_H; k++ {
                        if (v.block[i][j][k] == RAYTRACED_SPHERE){
                            /* we will raytrace this */
                            v.uniform_array[v.uniform_index] = float32(p.x + i)
                            v.uniform_index++
                            v.uniform_array[v.uniform_index] = float32(j)
                            v.uniform_index++
                            v.uniform_array[v.uniform_index] = float32(p.y + k)
                            v.uniform_index++
                        //~ } else if (v.block[i][j][k] == PAINTING){
                            //~ drawPainting(v.block[i][j][k],p.x + i, j,p.y + k)
                        } else if BlockIsLiquid[v.block[i][j][k]] {
                        } else if v.block[i][j][k] >= ITEM_0 {
                            drawmelater = append(drawmelater, [6]int{i, j, k, p.x + i, j,p.y + k})
                        } else if v.block[i][j][k] != AIR  {
                            DrawBlock(v.block[i][j][k],p.x + i, j,p.y + k)
                        }
                    }
                }
            }

            for _, p := range drawmelater {
                DrawBlock(v.block[p[0]][p[1]][p[2]], p[3], p[4], p[5])
            }

            var drawCloud = func ( x, y, z int) {
                var id BlockID
                id = CLOUD
                var v = [6][4][3]float32{
                    {{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // -Y
                }
                lv := float32(lightLevels - light) / float32(lightLevels)
                gl.Color3f(lv,lv,lv)
                for f := uint32(0); f < 1; f++ {
                    v0 := textVPos(id  )
                    v1 := textVPos(id+1)
                    h0 := texHPos(  f)
                    h1 := texHPos(f+1)
                    gl.TexCoord2f(h1, v1); gl.Vertex3f(float32(x)+v[f][0][0], float32(y)+v[f][0][1], float32(z)+v[f][0][2])
                    gl.TexCoord2f(h0, v1); gl.Vertex3f(float32(x)+v[f][1][0], float32(y)+v[f][1][1], float32(z)+v[f][1][2])
                    gl.TexCoord2f(h0, v0); gl.Vertex3f(float32(x)+v[f][2][0], float32(y)+v[f][2][1], float32(z)+v[f][2][2])
                    gl.TexCoord2f(h1, v0); gl.Vertex3f(float32(x)+v[f][3][0], float32(y)+v[f][3][1], float32(z)+v[f][3][2])
                }
            }
                
            gl.Color3f(1.0,1.0,1.0)
            
            var cloudyness byte = 50+byte(psrng1d(int((TICK/1000)+1))%150)
            //~ var cloudyness byte = byte(100)
            for i := 0; i < SUBCHUNK_H; i++ {
                for j := 0; j < SUBCHUNK_H; j++ {
                    if v.cloud[i][j] > cloudyness {
                        drawCloud(p.x + i, CHUNK_V-10,p.y + j)
                    }
                }
            }

            gl.End()
            gl.Disable(gl.TEXTURE_2D);
            gl.EndList()
            v.display_list = list
                        
            var drawLod = func ( id BlockID, x, y, z int) {
                biome_noise := biome_noise_function(x, z)
                underwater := false
                var n = [6][3]int{
                    { 1, 0, 0},
                    {-1, 0, 0},
                    { 0, 1, 0},
                    { 0,-1, 0},
                    { 0, 0, 1},
                    { 0, 0,-1},
                }
                oc := true
                lv := float32(1.0);
                for f := uint32(0); f < 6; f++ {
                    nb := getBlockVal(x + n[f][0], y + n[f][1],z + n[f][2])
                    l := float32(lightLevels - getLightVal(x + n[f][0], y + n[f][1],z + n[f][2]))/float32(lightLevels)
                    if (l < lv) && ((nb == AIR) || (nb == WATER)){
                        lv = l
                    }
                    if (nb == AIR) || ((id != WATER) && nb == WATER) {
                        oc = false
                        break
                    }
                    if nb == WATER {
                        underwater = true
                    }
                    
                }
                if oc {
                    return
                }
                
                switch id {
                case GRASS:
                    gl.Color3f(
                        //grass:    green          teal       purple       yellow
                        lerp4( 85.0/255.0,  85.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise)*lv,
                        lerp4(255.0/255.0, 235.0/255.0,  85.0/255.0, 255.0/255.0, biome_noise)*lv,
                        lerp4( 85.0/255.0, 235.0/255.0,  85.0/255.0,  32.0/255.0, biome_noise)*lv )
                case LEAF:
                    gl.Color3f(
                        lerp4( 00.0/255.0, 255.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise)*lv,
                        lerp4(255.0/255.0,  85.0/255.0,  85.0/255.0, 235.0/255.0, biome_noise)*lv,
                        lerp4( 00.0/255.0,  85.0/255.0, 255.0/255.0, 235.0/255.0, biome_noise)*lv )
                case WOOD:
                    gl.Color3f(
                        lerp4(168.0/255.0, 255.0/255.0, 168.0/255.0, 168.0/255.0, biome_noise)*lv,
                        lerp4( 84.0/255.0, 255.0/255.0,  86.0/255.0,  84.0/255.0, biome_noise)*lv,
                        lerp4(  0.0/255.0, 255.0/255.0,   0.0/255.0,   0.0/255.0, biome_noise)*lv )
                default:
                    r := float32(AtlasLOD[id%64][0][0]) / 255.0
                    g := float32(AtlasLOD[id%64][0][1]) / 255.0
                    b := float32(AtlasLOD[id%64][0][2]) / 255.0
                    gl.Color3f(r*lv, g*lv, b*lv)
                }
                
                c := float32(1.0)
                if underwater {
                    c = 0.5
                }
                gl.Vertex4f(float32(x)+0.5, float32(y)+0.5, float32(z)+0.5, c)
            }

            lodlist := gl.GenLists(1)
            gl.NewList(lodlist, gl.COMPILE)
            gl.Enable(gl.PROGRAM_POINT_SIZE)
            gl.Begin(gl.POINTS)
            for i := 0; i < SUBCHUNK_H; i++ {
                for j := 1; j < (CHUNK_V-1); j++ {
                    for k := 0; k < SUBCHUNK_H; k++ {
                        if (v.block[i][j][k] != AIR) && (!BlockIsLiquid[v.block[i][j][k]]) && (v.block[i][j][k] < ITEM_0)  {
                            drawLod(v.block[i][j][k],p.x + i, j,p.y + k)
                        }
                    }
                }
            }
            for i := 0; i < SUBCHUNK_H; i++ {
                for j := 0; j < SUBCHUNK_H; j++ {
                    if v.cloud[i][j] > cloudyness {
                        drawLod(CLOUD,p.x + i, CHUNK_V-10,p.y + j)
                    }
                }
            }

            gl.End()
            gl.EndList()
            v.LOD1 = lodlist
            
            
        }
    }

    var neighbor [RENDER_DISTANCE*2 + 1][RENDER_DISTANCE*2 + 1]*Chunk
    var neighborp [RENDER_DISTANCE*2 + 1][RENDER_DISTANCE*2 + 1]iVec2

    for i := 0; i < (RENDER_DISTANCE*2+1); i++ {
        for j := 0; j < (RENDER_DISTANCE*2+1); j++ {
            
            var p = iVec2{
                (int(player_pos.x) - int(player_pos.x)%16) + (i - RENDER_DISTANCE)*16,
                (int(player_pos.z) - int(player_pos.z)%16) + (j - RENDER_DISTANCE)*16,
            };
            if World[p] != nil {
                neighbor[i][j] = World[p]
                neighborp[i][j] = p
            } else {
                neighbor[i][j] = nil
                terrain(p)
            }
        }
    }
    
    tick := TICK%DAY_LEN
    if tick == 0 {
        Forecast = uint32(rand.Intn(int(WEATHER__N)))
    }
    //                   0 1 2 3 4 5 6 7 6 5 4 3 2 1 0
    day_cycle := []uint8{7,6,5,4,3,2,1,0,1,2,3,4,5,6,7}
    //~ day_cycle := []uint8{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
    light = day_cycle[ int(tick%DAY_LEN) / int(DAY_LEN/len(day_cycle)) ] 
    
    //~ if lastLightUpdate > (DAY_LEN/uint32(len(day_cycle))) {
    
    if (TICK % (DAY_LEN/uint32(len(day_cycle)))) == 0 {
        for _, v := range World {
            v.scheduledForLightUpdate = true
        }
        //~ lastLightUpdate = 0
    }
    //~ lastLightUpdate++
    
    start := time.Now() // capture start
	duration := 7 * time.Millisecond
    
    for i := 1; i < (RENDER_DISTANCE*2); i++ {
        for j := 1; j < (RENDER_DISTANCE*2); j++ {
            if time.Since(start) > duration {
                continue;
            }
            if neighbor[i][j] != nil {
                if neighbor[i][j].scheduledForLightUpdate {
                    ChunkUpdateDisplayListRef(neighbor[i][j])
                    neighbor[i][j].scheduledForLightUpdate = false
                }
            }
            if ( neighbor[i+1][j  ] != nil ) &&
               ( neighbor[i-1][j  ] != nil ) &&
               ( neighbor[i  ][j+1] != nil ) &&
               ( neighbor[i  ][j-1] != nil ) {
                BlockRandomUpdates(neighbor[i][j], iVec2{ neighborp[i][j].x,neighborp[i][j].y});
                if neighbor[i][j].display_list == 0 {
                    lightUpdate(neighbor[i][j], neighborp[i][j], light)
                    createDisplayLists(neighbor[i][j], neighborp[i][j])
                }
                scheduleCrafting(neighbor[i][j], neighborp[i][j])
            }
        }
    }

    // Bind texture each frame
    gl.Enable(gl.TEXTURE_2D)
    gl.BindTexture(gl.TEXTURE_2D, atlas)

    clearcolor := [4]float32{0.529, 0.808, 0.922, 1.0}
    if Forecast == WEATHER_EERIE {
        clearcolor[0] = 0.71
        clearcolor[1] = 0.75
        clearcolor[2] = 1.00
    }
    if Forecast == WEATHER_FOGGY {
        clearcolor[0] = 1.00
        clearcolor[1] = 1.00
        clearcolor[2] = 1.00
    }
    clearcolor[0] = clearcolor[0] * (float32(lightLevels-light)/float32(lightLevels))
    clearcolor[1] = clearcolor[1] * (float32(lightLevels-light)/float32(lightLevels))
    clearcolor[2] = clearcolor[2] * (float32(lightLevels-light)/float32(lightLevels))
    
    
    // Clear BEFORE draw
    gl.ClearColor(clearcolor[0],clearcolor[1],clearcolor[2],clearcolor[3])
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    
    const sun_radius = 120.0
    
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

    gl.DepthMask(false)
    gl.Enable(gl.TEXTURE_2D)
    gl.BindTexture(gl.TEXTURE_2D, TextureSun)
    
    
    gl.Color3f(1.0,1.0,1.0)
    drawSkybox(0, player_pos.x, player_pos.y, player_pos.z, 200, 0,0, (PI * (float32(TICK)/float32(DAY_LEN)) )+(PI/2.0)+(PI/2.0)+(PI/2.0) )

    gl.Disable(gl.TEXTURE_2D)
    gl.DepthMask(true)

    gl.Enable(gl.FOG);
    gl.Fogfv(gl.FOG_COLOR, &clearcolor[0])
    gl.Fogi(gl.FOG_MODE, gl.LINEAR)
    
    gl.Fogf(gl.FOG_START, float32((RENDER_DISTANCE/4)*SUBCHUNK_H) )
    
    gl.Fogf(gl.FOG_END, float32( RENDER_DISTANCE*SUBCHUNK_H ) )

    {
        cosHalfFOV := float32(math.Cos(float64(fov/2.0)))
        const CHUNK_SIZE = 16.0

        chunkVisible := func (cx, cz float32, px, pz float32, dx, dz float32, cosHalfFOV float32) bool {
            
            if (cam_pitch < -1.3) || (cam_pitch > 1.3) {
                return true
            }
            
            // Chunk corners
            corners := [4][2]float32{
                {cx, cz},
                {cx + CHUNK_SIZE, cz},
                {cx, cz + CHUNK_SIZE},
                {cx + CHUNK_SIZE, cz + CHUNK_SIZE},
            }
            for _, c := range corners {
                vx := c[0] - px
                vz := c[1] - pz
                // Reject chunks fully behind camera
                dot := vx*dx + vz*dz
                if dot <= 0 {
                    continue
                }
                // cos(theta) = dot / |v|
                len := float32(math.Sqrt(float64(vx*vx + vz*vz)))
                if dot/len >= cosHalfFOV {
                    return true
                }
            }
            return false
        }

        chunk_render_order = chunk_render_order[:0]
        sqr := func(f float32) float32 { return f * f }

        for i := 1; i < (RENDER_DISTANCE * 2); i++ {
            for j := 1; j < (RENDER_DISTANCE * 2); j++ {

                cx := float32(neighborp[i][j].x)
                cz := float32(neighborp[i][j].y)

                if !chunkVisible(
                    cx, cz,
                    player_pos.x, player_pos.z,
                    dir.x, dir.z,
                    cosHalfFOV,
                ) {
                    continue
                }

                chunk_render_order = append(chunk_render_order, struct {
                    i, j int
                    d    float32
                }{
                    i, j,
                    sqr(cx+(SUBCHUNK_H/2)-player_pos.x-dir.x) +
                        sqr(cz+(SUBCHUNK_H/2)-player_pos.z-dir.z),
                })
            }
        }

        sort.Slice(chunk_render_order, func(a, b int) bool {
            return chunk_render_order[a].d < chunk_render_order[b].d
        })

        schedule_raytraced_blocks_p = 0

        //~ gl.Disable(gl.BLEND) 
        gl.Enable(gl.BLEND)
        
        gl.Uniform1i(subpixelSnapUniformTextured, 1)
        gl.Uniform1i(subpixelSnapUniformSeed, int32(TICK))
        
        for _,v := range chunk_render_order {
            i := v.i
            j := v.j
            if neighbor[i][j] != nil {
                if neighbor[i][j].display_list != 0 {
                    const LOD_THRESHOLD = (RENDER_DISTANCE/2)*SUBCHUNK_H
                    if v.d < (LOD_THRESHOLD*LOD_THRESHOLD) {
                        gl.CallList(neighbor[i][j].display_list)
                        for l := 0; l < neighbor[i][j].uniform_index; l+=3 {
                            schedule_raytraced_blocks[schedule_raytraced_blocks_p] = neighbor[i][j].uniform_array[l]+0.5
                            schedule_raytraced_blocks_p++
                            schedule_raytraced_blocks[schedule_raytraced_blocks_p] = neighbor[i][j].uniform_array[l+1]+0.5
                            schedule_raytraced_blocks_p++
                            schedule_raytraced_blocks[schedule_raytraced_blocks_p] = neighbor[i][j].uniform_array[l+2]+0.5
                            schedule_raytraced_blocks_p++
                            //~ schedule_raytraced_blocks[schedule_raytraced_blocks_p] = RaytracingSource_donut
                            shapes := [3]float32{RaytracingSource_donut, RaytracingSource_octa, RaytracingSource_sphere}
                            //~ schedule_raytraced_blocks[schedule_raytraced_blocks_p] = RaytracingSource_octa
                            schedule_raytraced_blocks[schedule_raytraced_blocks_p] = shapes[psrng3d(
                                    int(neighbor[i][j].uniform_array[l]),
                                    int(neighbor[i][j].uniform_array[l+1]),
                                    int(neighbor[i][j].uniform_array[l+2]))%uint32(len(shapes))]
                            schedule_raytraced_blocks_p++
                        }                        
                    }
                }
            }
        }



        for _,v := range chunk_render_order {
            i := v.i
            j := v.j
            if neighbor[i][j] != nil {
                if neighbor[i][j].display_list != 0 {
                    const LOD_THRESHOLD = (RENDER_DISTANCE/2)*SUBCHUNK_H
                    if v.d < (LOD_THRESHOLD*LOD_THRESHOLD) {
                        /* rendering sign text out of the chunk drawing loop for correct blending */
                        for key, _ := range neighbor[i][j].blockText {                            
                            if getBlockVal(int(neighborp[i][j].x)+int(key.x), int(key.y), int(neighborp[i][j].y)+int(key.z)) != SIGN {
                                neighbor[i][j].blockText[key] = nil
                                continue
                            }
                            t_pos := Vec3{float32(neighborp[i][j].x)+float32(key.x), float32(key.y), float32(neighborp[i][j].y)+float32(key.z)}
                            a := Angle2D(
                                player_pos.x, player_pos.z, 
                                t_pos.x, t_pos.z )
                            DrawText3D(neighbor[i][j].blockText[key], t_pos.x, t_pos.y+1.5, t_pos.z+0.5, 0.25, 0, (3.14*270.0/180.0)-a, 0)
                        }
                    }
                }
            }
        }
        

        gl.Uniform1i(subpixelSnapUniformTextured, 0)
        for _,v := range chunk_render_order {
            i := v.i
            j := v.j
            if neighbor[i][j] != nil {
                if neighbor[i][j].display_list != 0 {
                    
                    const LOD_THRESHOLD = (RENDER_DISTANCE/2)*SUBCHUNK_H
                    if v.d < (LOD_THRESHOLD*LOD_THRESHOLD) {
                    } else {
                        gl.CallList(neighbor[i][j].LOD1)
                    }
                }
            }
        }    
        gl.Enable(gl.BLEND)
        gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
        gl.Uniform1i(subpixelSnapUniformTextured, 1)
    }
    
    /* TODO: PUT THIS IN THE MAIN FUNCTION */
    processMobs() 
    DrawItemDrops()
    
    gl.Disable(gl.FOG);
}

func scheduleCrafting(c *Chunk, p iVec2) {
    checkForRecipe := func(r Recipe, inv Inventory) bool {
        for kr, vr := range r.from {
            foo := false
            for ki, vi := range inv.ID {
                if vi == vr {
                    if inv.C[ki] >= r.fromc[kr] {
                        foo = true
                        break
                    }
                }
            }
            if foo == false {
                return false
            }
        }
        return true
    }
    subtractForRecipe := func(r Recipe, inv Inventory) {
        for kr, vr := range r.from {
            for ki, vi := range inv.ID {
                if vi == vr {
                    if inv.C[ki] >= r.fromc[kr] {
                        inv.C[ki] -= r.fromc[kr]
                    }
                }
            }
        }
    }
    subtractFromIventory := func(id BlockID, count int32, inv Inventory) {
        for ki, vi := range inv.ID {
            if vi == id {
                inv.C[ki] -= count
            }
        }
    }
    if c.blockInventories != nil {
        for k, v :=  range c.blockInventories {
            if (c.block[k.x][k.y][k.z] == CRAFTING_TABLE) || (c.block[k.x][k.y][k.z] == FURNACE) {
                if len(c.blockInventories[k].ID) == 0 {
                    continue
                }
                if c.blockInventories[k].T > 100 {
                    bar := c.blockInventories[k]
                    
                    var recipesource []Recipe
                    
                    if c.block[k.x][k.y][k.z] == CRAFTING_TABLE {
                        recipesource = crafting_table_recipes
                    } else {
                        recipesource = furnace_recipes
                    }
                    
                    for _, rv := range recipesource {
                        if checkForRecipe(rv, v) {
                            bar := c.blockInventories[k]
                            bar.T = 0
                            c.blockInventories[k] = bar
                            subtractForRecipe(rv,v)
                            for rc, _ := range rv.to {
                                itemDrops = append(itemDrops, Drop{
                                    pos: Vec3{float32(p.x + int(k.x)), float32(k.y)+2.5, float32(p.y + int(k.z))},
                                    dir: Vec3{rand.Float32()-0.5,0.5,rand.Float32()-0.5},
                                    id: rv.to[rc],
                                    c: rv.toc[rc] })
                            }
                            goto processed_crafting
                        }
                    }
                    for ik, _ := range c.blockInventories[k].ID {
                        itemDrops = append(itemDrops, Drop{
                            pos: Vec3{float32(p.x + int(k.x)), float32(k.y)+2.5, float32(p.y + int(k.z))},
                            dir: Vec3{rand.Float32()-0.5,0.5,rand.Float32()-0.5},
                            id: c.blockInventories[k].ID[ik],
                            c: c.blockInventories[k].C[ik] })
                        
                        subtractFromIventory(bar.ID[ik], bar.C[ik], bar)
                    }
                    processed_crafting:
                    inventoryClearEmpty(&bar)
                    c.blockInventories[k] = bar
                    
                } else {
                    bar := c.blockInventories[k]
                    bar.T++
                    c.blockInventories[k] = bar
                }
                
            } else {
                /* if we are neither an furnace nor a crafting table drop everything */
                bar := c.blockInventories[k]
                for ik, _ := range c.blockInventories[k].ID {
                    itemDrops = append(itemDrops, Drop{
                        pos: Vec3{float32(p.x + int(k.x)), float32(k.y)+2.5, float32(p.y + int(k.z))},
                        dir: Vec3{rand.Float32()-0.5,0.5,rand.Float32()-0.5},
                        id: c.blockInventories[k].ID[ik],
                        c: c.blockInventories[k].C[ik] })
                    
                    subtractFromIventory(bar.ID[ik], bar.C[ik], bar)
                }
            }
        }
    }
}


// Angle2D returns the angle between p1 and p2 in radians.
func Angle2D (x1, y1, x2, y2 float32) float32 {
    dx := float64(x2 - x1)
    dy := float64(y2 - y1)
    return float32(math.Atan2(dy, dx)) // angle relative to X axis
}

func DrawItemDrops() {
    
    absf := func(f float32) float32 {
        if f < 0 { return -f }
        return f
    }
    absi := func(i int) int {
        if i < 0 { return -i }
        return i
    }
    
    gl.Enable(gl.TEXTURE_2D)
    gl.BindTexture(gl.TEXTURE_2D, atlas)
    gl.Begin(gl.QUADS)
    filter := itemDrops[:0]
    for key, v := range itemDrops {
        
        drawCube(v.id, v.pos.x, v.pos.y + sinf(float32(TICK)/10.0)/50.0, v.pos.z, 0.25, 0, float32(TICK)/50.0, 0)
        
        Taxicab := func(a, b Vec3) float32 {
            return absf(a.x-b.x)+absf(a.y-b.y)+absf(a.z-b.z)
        }
        Taxicabi := func(a, b [3]int) int {
            return absi(a[0]-b[0])+absi(a[1]-b[1])+absi(a[2]-b[2])
        }
        
        a := Angle2D( player_pos.x, player_pos.z, v.pos.x, v.pos.z )
        DrawNumbers( v.c, v.pos.x, v.pos.y + 0.255, v.pos.z, 0.25, 0, (3.14*270.0/180.0)-a, 0)
        
        var npbb BoundingBox
        var ibox BoundingBox        
        var valid_position_x = true
        var valid_position_y = true
        var valid_position_z = true
        var bb_range = 2
        const vel = 0.2
        np := Vec3{itemDrops[key].dir.x*vel,itemDrops[key].dir.y*vel,itemDrops[key].dir.z*vel}

        var nearbyinventoryblock bool
        nearbyinventoryblock = false
        var nearestDistance int
        nearestDistance = 99999
        var nearestCraftingTable [3]int

        traversible := [N_BLOCK_IDS]bool{
            AIR:true,
            WATER:true,
            LAVA:true,
        }

        xi := int(itemDrops[key].pos.x)
        yi := int(itemDrops[key].pos.y)
        zi := int(itemDrops[key].pos.z)
        
        for i := -bb_range; i < bb_range; i++ {
            for j := -bb_range; j < bb_range; j++ {
                for k := -bb_range; k < bb_range; k++ {
                    nx := xi+i
                    ny := yi+j
                    nz := zi+k
                    
                    if ny > (CHUNK_V - bb_range - 1) {
                        continue
                    }
                    if ny < bb_range {
                        continue
                    }
                    
                    if b := getBlockVal(nx, ny, nz); (b == CRAFTING_TABLE) || (b == FURNACE) {
                        nearbyinventoryblock = true
                        d := Taxicabi([3]int{nx,ny,nz}, [3]int{xi,yi,zi})
                        if nearestDistance > d {
                           nearestDistance = d
                           nearestCraftingTable = [3]int{nx,ny,nz}
                        }
                    }
                    
                    ibox.X = itemDrops[key].pos.x + np.x - 0.125
                    ibox.Y = itemDrops[key].pos.y - 0.125
                    ibox.Z = itemDrops[key].pos.z - 0.125
                    ibox.W = 0.25
                    ibox.H = 0.25
                    ibox.L = 0.25
                    npbb.X = float32(int(itemDrops[key].pos.x + np.x)+i)
                    npbb.Y = float32(int(itemDrops[key].pos.y       )+j)
                    npbb.Z = float32(int(itemDrops[key].pos.z       )+k)
                    npbb.W = 1.0
                    npbb.H = 1.0
                    npbb.L = 1.0
                    if chunkExists(
                    int(itemDrops[key].pos.x + np.x)+i,
                    int(itemDrops[key].pos.y       )+j,
                    int(itemDrops[key].pos.z       )+k ){
                        if !traversible[getBlockVal(
                        int(itemDrops[key].pos.x + np.x)+i,
                        int(itemDrops[key].pos.y       )+j,
                        int(itemDrops[key].pos.z       )+k )] {
                            if intersect(ibox, npbb) {
                                valid_position_x = false
                            }
                        }
                    }
                    ibox.X = itemDrops[key].pos.x - 0.125
                    ibox.Y = itemDrops[key].pos.y + np.y - 0.125
                    ibox.Z = itemDrops[key].pos.z - 0.125
                    ibox.W = 0.25
                    ibox.H = 0.25
                    ibox.L = 0.25
                    npbb.X = float32(int(itemDrops[key].pos.x       )+i)
                    npbb.Y = float32(int(itemDrops[key].pos.y + np.y)+j)
                    npbb.Z = float32(int(itemDrops[key].pos.z       )+k)
                    npbb.W = 1.0
                    npbb.H = 1.0
                    npbb.L = 1.0
                    if chunkExists(
                    int(itemDrops[key].pos.x       )+i,
                    int(itemDrops[key].pos.y + np.y)+j,
                    int(itemDrops[key].pos.z       )+k ){
                        if !traversible[getBlockVal(
                        int(itemDrops[key].pos.x       )+i,
                        int(itemDrops[key].pos.y + np.y)+j,
                        int(itemDrops[key].pos.z       )+k )] {
                            if intersect(ibox, npbb) {
                                valid_position_y = false
                            }
                        }
                    }
                    ibox.X = itemDrops[key].pos.x - 0.125
                    ibox.Y = itemDrops[key].pos.y - 0.125
                    ibox.Z = itemDrops[key].pos.z + np.z - 0.125
                    ibox.W = 0.25
                    ibox.H = 0.25
                    ibox.L = 0.25
                    npbb.X = float32(int(itemDrops[key].pos.x       )+i)
                    npbb.Y = float32(int(itemDrops[key].pos.y       )+j)
                    npbb.Z = float32(int(itemDrops[key].pos.z + np.z)+k)
                    npbb.W = 1.0
                    npbb.H = 1.0
                    npbb.L = 1.0
                    if chunkExists(
                    int(itemDrops[key].pos.x       )+i,
                    int(itemDrops[key].pos.y       )+j,
                    int(itemDrops[key].pos.z + np.z)+k ){
                        if !traversible[getBlockVal(
                        int(itemDrops[key].pos.x       )+i,
                        int(itemDrops[key].pos.y       )+j,
                        int(itemDrops[key].pos.z + np.z)+k )] {
                            if intersect(ibox, npbb) {
                                valid_position_z = false
                            }
                        }
                    }
                }
            }
        }
        
        if valid_position_x {
            itemDrops[key].pos.x += np.x
        } else {
            itemDrops[key].dir.x = 0
        }
        if valid_position_y {
            itemDrops[key].pos.y += np.y
            if getBlockVal(xi, yi, zi) == WATER {
                itemDrops[key].dir.y = 0.16
            } else {
                itemDrops[key].dir.y -= 0.02
            }
        } else {
            itemDrops[key].dir.x *= 0.8
            itemDrops[key].dir.z *= 0.8
            itemDrops[key].dir.y = 0
        }

        if valid_position_z {
            itemDrops[key].pos.z += np.z
        } else {
            itemDrops[key].dir.z = 0
        }
        
        //~ if !valid_position_y {
            pv := player_pos
            pv.y -= 1.5
            if (nearbyinventoryblock) && (nearestDistance == 1) {
                insertBlockInventory(
                    itemDrops[key].id,
                    nearestCraftingTable[0], nearestCraftingTable[1], nearestCraftingTable[2], 
                    Vec3{itemDrops[key].dir.x, itemDrops[key].dir.y, itemDrops[key].dir.z} )
            }else if Taxicab(itemDrops[key].pos, pv) > 1.0 {
                filter = append(filter, itemDrops[key])
            } else {
                playerInventoryInsert(v.id, v.c)
            }
        //~ } else {
            //~ filter = append(filter, itemDrops[key])
        //~ }
    }
    itemDrops = filter
    
    gl.End()
    gl.Disable(gl.TEXTURE_2D)
}

var getBlockValA *[SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]BlockID
var getBlockValP iVec2
func getBlockVal(x, y, z int) BlockID {
    if (getBlockValP.x == (x-x%SUBCHUNK_H)) && (getBlockValP.y == (z-z%SUBCHUNK_H)) {
        return getBlockValA[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
    }
    getBlockValP = iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}
    getBlockValA = &World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].block
    return World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].block[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
}

var getBlockRefA *[SUBCHUNK_H][CHUNK_V][SUBCHUNK_H]BlockID
var getBlockRefP iVec2
func getBlockRef(x, y, z int) *BlockID {
    if (getBlockRefP.x == (x-x%SUBCHUNK_H)) && (getBlockRefP.y == (z-z%SUBCHUNK_H)) {
        return &getBlockRefA[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
    }
    getBlockRefP = iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}
    getBlockRefA = &World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].block
    return &World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}].block[x%SUBCHUNK_H][y][z%SUBCHUNK_H]
}

func insertBlockInventory(id BlockID, x, y, z int, direction Vec3) {
    var c *Chunk 
    c = World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}] 
    if c.blockInventories == nil {
        c.blockInventories = make(map[nearVec]Inventory)
    }
    var nv nearVec
    nv.x = uint8(x%SUBCHUNK_H)
    nv.y = uint8(y)
    nv.z = uint8(z%SUBCHUNK_H)
    
    inv := c.blockInventories[nv] // copy from map
    for k, v := range inv.ID {
        if v == id {
            inv.C[k]++
            c.blockInventories[nv] = inv // store back
            return
        }
    }
    inv.ID = append(inv.ID, id)
    inv.C = append(inv.C, 1)
    inv.T = 0
    c.blockInventories[nv] = inv // store back
        
    log.Println(c.blockInventories[nv])
}

func chunkExists(x, y, z int) bool {
    if y > (CHUNK_V-1) { return false }
    if y < (0) { return false }
    return World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}] != nil
}
func chunkUpdateDisplaylist(x, y, z int) {
    c := World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}]
    
    if c.display_list != 0 {
        gl.DeleteLists(c.display_list, 1);
        c.display_list = 0;
    }

    if c.LOD1 != 0 {
        gl.DeleteLists(c.LOD1, 1);
        c.display_list = 0;
    }
    
}
func ChunkUpdateDisplayListRef(c *Chunk) {
    if c.display_list != 0 {
        gl.DeleteLists(c.display_list, 1);
        c.display_list = 0;
    }
    if c.LOD1 != 0 {
        gl.DeleteLists(c.LOD1, 1);
        c.display_list = 0;
    }
}

func drawCube(id BlockID, x, y, z float32, size float32, angleX, angleY, angleZ float32) {
    
    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }

    // Original cube vertex positions (unit cube)
    var v = [6][4][3]float32{
        {{ 0.5, -0.5,  0.5}, { 0.5, -0.5, -0.5}, { 0.5,  0.5, -0.5}, { 0.5,  0.5,  0.5}}, // +X
        {{-0.5, -0.5, -0.5}, {-0.5, -0.5,  0.5}, {-0.5,  0.5,  0.5}, {-0.5,  0.5, -0.5}}, // -X
        {{-0.5,  0.5, -0.5}, {-0.5,  0.5,  0.5}, { 0.5,  0.5,  0.5}, { 0.5,  0.5, -0.5}}, // +Y
        {{-0.5, -0.5, -0.5}, { 0.5, -0.5, -0.5}, { 0.5, -0.5,  0.5}, {-0.5, -0.5,  0.5}}, // -Y
        {{-0.5, -0.5,  0.5}, { 0.5, -0.5,  0.5}, { 0.5,  0.5,  0.5}, {-0.5,  0.5,  0.5}}, // +Z
        {{ 0.5, -0.5, -0.5}, {-0.5, -0.5, -0.5}, {-0.5,  0.5, -0.5}, { 0.5,  0.5, -0.5}}, // -Z
    }

    // Precompute sine/cosine
    sx, cx := float32(0.0), float32(1.0)
    if angleX != 0 {
        sx, cx = float32(math.Sin(float64(angleX))), float32(math.Cos(float64(angleX)))
    }
    
    sy, cy := float32(0.0), float32(1.0)
    if angleY != 0 {
        sy, cy = float32(math.Sin(float64(angleY))), float32(math.Cos(float64(angleY)))
    }
    
    sz, cz := float32(0.0), float32(1.0)
    if angleZ != 0 {
        sz, cz = float32(math.Sin(float64(angleZ))), float32(math.Cos(float64(angleZ)))
    }

    rotate := func(px, py, pz float32) (float32, float32, float32) {
        // Rotate X
        y1 := py*cx - pz*sx
        z1 := py*sx + pz*cx
        x1 := px
        // Rotate Y
        x2 := x1*cy + z1*sy
        z2 := -x1*sy + z1*cy
        y2 := y1
        // Rotate Z
        x3 := x2*cz - y2*sz
        y3 := x2*sz + y2*cz
        z3 := z2
        return x3, y3, z3
    }

    // Precompute all rotated vertices once
    var rotated [6][4][3]float32
    for f := 0; f < 6; f++ {
        for i := 0; i < 4; i++ {
            rx, ry, rz := rotate(v[f][i][0]*size, v[f][i][1]*size, v[f][i][2]*size)
            rotated[f][i][0] = rx
            rotated[f][i][1] = ry
            rotated[f][i][2] = rz
        }
    }

    v0 := textVPos(id)
    v1 := textVPos(id + 1)
    
    gl.Begin(gl.QUADS)
    
    for f := 0; f < 6; f++ {
        
        h0 := texHPos(uint32(f))
        h1 := texHPos(uint32(f + 1))    
        gl.TexCoord2f(h1, v1)
        gl.Vertex3f(x+rotated[f][0][0], y+rotated[f][0][1], z+rotated[f][0][2])
        gl.TexCoord2f(h0, v1)
        gl.Vertex3f(x+rotated[f][1][0], y+rotated[f][1][1], z+rotated[f][1][2])
        gl.TexCoord2f(h0, v0)
        gl.Vertex3f(x+rotated[f][2][0], y+rotated[f][2][1], z+rotated[f][2][2])
        gl.TexCoord2f(h1, v0)
        gl.Vertex3f(x+rotated[f][3][0], y+rotated[f][3][1], z+rotated[f][3][2])
        
    }

    gl.End()
}


func drawSkybox(id BlockID, x, y, z float32, size float32, angleX, angleY, angleZ float32) {
    
    // Original cube vertex positions (unit cube)
    var v = [6][4][3]float32{
        {{ 0.5, -0.5,  0.5}, { 0.5,  0.5,  0.5}, { 0.5,  0.5, -0.5}, { 0.5, -0.5, -0.5}}, // +X
        {{-0.5, -0.5, -0.5}, {-0.5,  0.5, -0.5}, {-0.5,  0.5,  0.5}, {-0.5, -0.5,  0.5}}, // -X
        {{-0.5,  0.5, -0.5}, { 0.5,  0.5, -0.5}, { 0.5,  0.5,  0.5}, {-0.5,  0.5,  0.5}}, // +Y
        {{-0.5, -0.5, -0.5}, {-0.5, -0.5,  0.5}, { 0.5, -0.5,  0.5}, { 0.5, -0.5, -0.5}}, // -Y
        {{-0.5, -0.5,  0.5}, {-0.5,  0.5,  0.5}, { 0.5,  0.5,  0.5}, { 0.5, -0.5,  0.5}}, // +Z
        {{ 0.5, -0.5, -0.5}, { 0.5,  0.5, -0.5}, {-0.5,  0.5, -0.5}, {-0.5, -0.5, -0.5}}, // -Z
    }

    // Precompute sine/cosine
    sx, cx := float32(0.0), float32(1.0)
    if angleX != 0 {
        sx, cx = float32(math.Sin(float64(angleX))), float32(math.Cos(float64(angleX)))
    }
    
    sy, cy := float32(0.0), float32(1.0)
    if angleY != 0 {
        sy, cy = float32(math.Sin(float64(angleY))), float32(math.Cos(float64(angleY)))
    }
    
    sz, cz := float32(0.0), float32(1.0)
    if angleZ != 0 {
        sz, cz = float32(math.Sin(float64(angleZ))), float32(math.Cos(float64(angleZ)))
    }

    rotate := func(px, py, pz float32) (float32, float32, float32) {
        // Rotate X
        y1 := py*cx - pz*sx
        z1 := py*sx + pz*cx
        x1 := px
        // Rotate Y
        x2 := x1*cy + z1*sy
        z2 := -x1*sy + z1*cy
        y2 := y1
        // Rotate Z
        x3 := x2*cz - y2*sz
        y3 := x2*sz + y2*cz
        z3 := z2
        return x3, y3, z3
    }

    // Precompute all rotated vertices once
    var rotated [6][4][3]float32
    for f := 0; f < 6; f++ {
        for i := 0; i < 4; i++ {
            rx, ry, rz := rotate(v[f][i][0]*size, v[f][i][1]*size, v[f][i][2]*size)
            rotated[f][i][0] = rx
            rotated[f][i][1] = ry
            rotated[f][i][2] = rz
        }
    }

    v0 := float32(0.0)
    v1 := float32(1.0)
    
    gl.Begin(gl.QUADS)
    
    //~ for f := 0; f < 6; f++ {
        var f int 
        var h0, h1 float32
        
        f = 2
        h0 = float32(0.0)
        h1 = float32(0.5)
        gl.TexCoord2f(h1, v1); gl.Vertex3f(x+rotated[f][0][0], y+rotated[f][0][1], z+rotated[f][0][2])
        gl.TexCoord2f(h0, v1); gl.Vertex3f(x+rotated[f][1][0], y+rotated[f][1][1], z+rotated[f][1][2])
        gl.TexCoord2f(h0, v0); gl.Vertex3f(x+rotated[f][2][0], y+rotated[f][2][1], z+rotated[f][2][2])
        gl.TexCoord2f(h1, v0); gl.Vertex3f(x+rotated[f][3][0], y+rotated[f][3][1], z+rotated[f][3][2])

        f = 3
        h0 = float32(0.5)
        h1 = float32(1.0)
        gl.TexCoord2f(h1, v1); gl.Vertex3f(x+rotated[f][0][0], y+rotated[f][0][1], z+rotated[f][0][2])
        gl.TexCoord2f(h0, v1); gl.Vertex3f(x+rotated[f][1][0], y+rotated[f][1][1], z+rotated[f][1][2])
        gl.TexCoord2f(h0, v0); gl.Vertex3f(x+rotated[f][2][0], y+rotated[f][2][1], z+rotated[f][2][2])
        gl.TexCoord2f(h1, v0); gl.Vertex3f(x+rotated[f][3][0], y+rotated[f][3][1], z+rotated[f][3][2])
        
    //~ }

    gl.End()
}


func drawExtruded(id BlockID, x, y, z float32, size float32, angleX, angleY, angleZ float32) {
    
    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }

    var v = [6][4][3]float32{
        {{ 0.5+0.0, -0.5,  0.5}, { 0.5+0.0, -0.5, -0.5}, { 0.5+0.0,  0.5, -0.5}, { 0.5+0.0,  0.5,  0.5}}, // +X
        {{ 0.5+0.1, -0.5,  0.5}, { 0.5+0.1, -0.5, -0.5}, { 0.5+0.1,  0.5, -0.5}, { 0.5+0.1,  0.5,  0.5}}, // +X
        {{ 0.5+0.2, -0.5,  0.5}, { 0.5+0.2, -0.5, -0.5}, { 0.5+0.2,  0.5, -0.5}, { 0.5+0.2,  0.5,  0.5}}, // +X
        {{ 0.5+0.3, -0.5,  0.5}, { 0.5+0.3, -0.5, -0.5}, { 0.5+0.3,  0.5, -0.5}, { 0.5+0.3,  0.5,  0.5}}, // +X
        {{ 0.5+0.4, -0.5,  0.5}, { 0.5+0.4, -0.5, -0.5}, { 0.5+0.4,  0.5, -0.5}, { 0.5+0.4,  0.5,  0.5}}, // +X
        {{ 0.5+0.5, -0.5,  0.5}, { 0.5+0.5, -0.5, -0.5}, { 0.5+0.5,  0.5, -0.5}, { 0.5+0.5,  0.5,  0.5}}, // +X
    }

    sx, cx := float32(0.0), float32(1.0)
    if angleX != 0 {
        sx, cx = float32(math.Sin(float64(angleX))), float32(math.Cos(float64(angleX)))
    }
    
    sy, cy := float32(0.0), float32(1.0)
    if angleY != 0 {
        sy, cy = float32(math.Sin(float64(angleY))), float32(math.Cos(float64(angleY)))
    }
    
    sz, cz := float32(0.0), float32(1.0)
    if angleZ != 0 {
        sz, cz = float32(math.Sin(float64(angleZ))), float32(math.Cos(float64(angleZ)))
    }

    rotate := func(px, py, pz float32) (float32, float32, float32) {
        // Rotate X
        y1 := py*cx - pz*sx
        z1 := py*sx + pz*cx
        x1 := px
        // Rotate Y
        x2 := x1*cy + z1*sy
        z2 := -x1*sy + z1*cy
        y2 := y1
        // Rotate Z
        x3 := x2*cz - y2*sz
        y3 := x2*sz + y2*cz
        z3 := z2
        return x3, y3, z3
    }

    var rotated [6][4][3]float32
    for f := 0; f < 6; f++ {
        for i := 0; i < 4; i++ {
            rx, ry, rz := rotate(v[f][i][0]*size, v[f][i][1]*size, v[f][i][2]*size)
            rotated[f][i][0] = rx
            rotated[f][i][1] = ry
            rotated[f][i][2] = rz
        }
    }

    v0 := textVPos(id)
    v1 := textVPos(id + 1)
    
    gl.Begin(gl.QUADS)
    
    for f := 0; f < 6; f++ {
        
        h0 := texHPos(uint32(0))
        h1 := texHPos(uint32(0 + 1))    
        gl.TexCoord2f(h1, v1)
        gl.Vertex3f(x+rotated[f][0][0], y+rotated[f][0][1], z+rotated[f][0][2])
        gl.TexCoord2f(h0, v1)
        gl.Vertex3f(x+rotated[f][1][0], y+rotated[f][1][1], z+rotated[f][1][2])
        gl.TexCoord2f(h0, v0)
        gl.Vertex3f(x+rotated[f][2][0], y+rotated[f][2][1], z+rotated[f][2][2])
        gl.TexCoord2f(h1, v0)
        gl.Vertex3f(x+rotated[f][3][0], y+rotated[f][3][1], z+rotated[f][3][2])
        
    }

    gl.End()
}

/* draws a billboard without creating a new view matrix */
func drawBillboard_fast(
        id BlockID,
        horizontal_offset uint32,
        x, y, z float32,
        size float32) {
            
        
    var camDirX, camDirZ float32 = dir.x, dir.z

    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }

    v0 := textVPos(id)
    v1 := textVPos(id + 1)
    h0 := texHPos(0 + horizontal_offset)
    h1 := texHPos(1 + horizontal_offset)

    var vx float32 = camDirX
    var vy float32 = 0
    var vz float32 = camDirZ

    len := float32(math.Sqrt(float64(vx*vx + vy*vy + vz*vz)))
    if len < 0.00001 {
        vx, vz = 0, 1 // fallback
    } else {
        inv := float32(1.0 / len)
        vx *= inv
        vz *= inv
    }

    var rightX float32 =  vz
    var rightY float32 =  0
    var rightZ float32 = -vx

    upX := float32(0)
    upY := float32(1)
    upZ := float32(0)

    var half float32 = size * 0.5

    px0 := x + (-half)*rightX + (-half)*upX
    py0 := y + (-half)*rightY + (-half)*upY
    pz0 := z + (-half)*rightZ + (-half)*upZ

    px1 := x + ( half)*rightX + (-half)*upX
    py1 := y + ( half)*rightY + (-half)*upY
    pz1 := z + ( half)*rightZ + (-half)*upZ

    px2 := x + ( half)*rightX + ( half)*upX
    py2 := y + ( half)*rightY + ( half)*upY
    pz2 := z + ( half)*rightZ + ( half)*upZ

    px3 := x + (-half)*rightX + ( half)*upX
    py3 := y + (-half)*rightY + ( half)*upY
    pz3 := z + (-half)*rightZ + ( half)*upZ

    gl.Begin(gl.QUADS)

        gl.TexCoord2f(h0, v1)
        gl.Vertex3f(px0, py0, pz0)

        gl.TexCoord2f(h1, v1)
        gl.Vertex3f(px1, py1, pz1)

        gl.TexCoord2f(h1, v0)
        gl.Vertex3f(px2, py2, pz2)

        gl.TexCoord2f(h0, v0)
        gl.Vertex3f(px3, py3, pz3)

    gl.End()
}



func drawBillboard(id BlockID, horizontal_offset uint32, x, y, z, size float32) {
    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }

    v0 := textVPos(id)
    v1 := textVPos(id + 1)
    h0 := texHPos(0+horizontal_offset)
    h1 := texHPos(1+horizontal_offset)

    gl.PushMatrix()
    gl.Translatef(x, y, z)
    
    var modelview [16]float32
    gl.GetFloatv(gl.MODELVIEW_MATRIX, &modelview[0])

    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if i == j {
                modelview[i*4+j] = 1.0
            } else {
                modelview[i*4+j] = 0.0
            }
        }
    }
    gl.LoadMatrixf(&modelview[0])

    half := size * 0.5
    gl.Begin(gl.QUADS)
        gl.TexCoord2f(h0, v1)
        gl.Vertex3f(-half, -half, 0)

        gl.TexCoord2f(h1, v1)
        gl.Vertex3f( half, -half, 0)

        gl.TexCoord2f(h1, v0)
        gl.Vertex3f( half,  half, 0)

        gl.TexCoord2f(h0, v0)
        gl.Vertex3f(-half,  half, 0)
    gl.End()

    gl.PopMatrix()
}

func DrawNumbers(num int32, x, y, z, size, ax, ay, az float32) {
    
    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }
    
    // Original cube vertex positions (unit cube)
    var v = [6][4][3]float32{
        {{ 0.5,-0.5, 0.5}, { 0.5,-0.5,-0.5}, { 0.5, 0.5,-0.5}, { 0.5, 0.5, 0.5}}, // +X
        {{-0.5,-0.5,-0.5}, {-0.5,-0.5, 0.5}, {-0.5, 0.5, 0.5}, {-0.5, 0.5,-0.5}}, // -X
        {{-0.5, 0.5,-0.5}, {-0.5, 0.5, 0.5}, { 0.5, 0.5, 0.5}, { 0.5, 0.5,-0.5}}, // +Y
        {{-0.5,-0.5,-0.5}, { 0.5,-0.5,-0.5}, { 0.5,-0.5, 0.5}, {-0.5,-0.5, 0.5}}, // -Y
        {{-0.5,-0.5, 0.5}, { 0.5,-0.5, 0.5}, { 0.5, 0.5, 0.5}, {-0.5, 0.5, 0.5}}, // +Z
        {{ 0.5,-0.5,-0.5}, {-0.5,-0.5,-0.5}, {-0.5, 0.5,-0.5}, { 0.5, 0.5,-0.5}}, // -Z
    }

    // Precompute sine/cosine
    sx, cx := float32(math.Sin(float64(ax))), float32(math.Cos(float64(ax)))
    sy, cy := float32(math.Sin(float64(ay))), float32(math.Cos(float64(ay)))
    sz, cz := float32(math.Sin(float64(az))), float32(math.Cos(float64(az)))

    rotate := func(px, py, pz float32) (float32, float32, float32) {
        // Rotate X
        y1 := py*cx - pz*sx
        z1 := py*sx + pz*cx
        x1 := px
        // Rotate Y
        x2 := x1*cy + z1*sy
        z2 := -x1*sy + z1*cy
        y2 := y1
        // Rotate Z
        x3 := x2*cz - y2*sz
        y3 := x2*sz + y2*cz
        z3 := z2
        return x3, y3, z3
    }
    
    v2 := textVPos(0)
    v3 := textVPos(1)

    // Precompute all rotated vertices once
    var rotated [6][4][3]float32
    for f := 0; f < 6; f++ {
        for i := 0; i < 4; i++ {
            rx, ry, rz := rotate(v[f][i][0]*size, v[f][i][1]*size, v[f][i][2]*size)
            rotated[f][i][0] = rx
            rotated[f][i][1] = ry
            rotated[f][i][2] = rz
        }
    }
    
    n := num
    i := 0
    t := float32(int(math.Log10(float64(num))))
    if t == 0 {
        t = 1
    }

    gl.Begin(gl.QUADS)
    
    for {
        m := n%10
        const camside = 4
        n0 := texHPos(uint32(6 + m + 1))
        n1 := texHPos(uint32(6 + m))
        gl.TexCoord2f(n1, v3); gl.Vertex3f(x+rotated[camside][0][0]/t-float32(i)*size, y+rotated[camside][0][1], z+rotated[camside][0][2]/t)
        gl.TexCoord2f(n0, v3); gl.Vertex3f(x+rotated[camside][1][0]/t-float32(i)*size, y+rotated[camside][1][1], z+rotated[camside][1][2]/t)
        gl.TexCoord2f(n0, v2); gl.Vertex3f(x+rotated[camside][2][0]/t-float32(i)*size, y+rotated[camside][2][1], z+rotated[camside][2][2]/t)
        gl.TexCoord2f(n1, v2); gl.Vertex3f(x+rotated[camside][3][0]/t-float32(i)*size, y+rotated[camside][3][1], z+rotated[camside][3][2]/t)
        n /= 10
        i++
        if n == 0 {
            break
        }
    }    
    
    gl.End()
}

func DrawText3D(str []byte, x, y, z, size, ax, ay, az float32) {
    
    texHPos := func(face uint32) float32 {
        return float32(face) * (16.0 / 1024.0)
    }
    textVPos := func(id BlockID) float32 {
        return float32(id) * (16.0 / 1024.0)
    }
    
    // Original cube vertex positions (unit cube)
    var v = [6][4][3]float32{
        {{ 0.5,-0.5, 0.5}, { 0.5,-0.5,-0.5}, { 0.5, 0.5,-0.5}, { 0.5, 0.5, 0.5}}, // +X
        {{-0.5,-0.5,-0.5}, {-0.5,-0.5, 0.5}, {-0.5, 0.5, 0.5}, {-0.5, 0.5,-0.5}}, // -X
        {{-0.5, 0.5,-0.5}, {-0.5, 0.5, 0.5}, { 0.5, 0.5, 0.5}, { 0.5, 0.5,-0.5}}, // +Y
        {{-0.5,-0.5,-0.5}, { 0.5,-0.5,-0.5}, { 0.5,-0.5, 0.5}, {-0.5,-0.5, 0.5}}, // -Y
        {{-0.5,-0.5, 0.5}, { 0.5,-0.5, 0.5}, { 0.5, 0.5, 0.5}, {-0.5, 0.5, 0.5}}, // +Z
        {{ 0.5,-0.5,-0.5}, {-0.5,-0.5,-0.5}, {-0.5, 0.5,-0.5}, { 0.5, 0.5,-0.5}}, // -Z
    }

    // Precompute sine/cosine
    sx, cx := float32(math.Sin(float64(ax))), float32(math.Cos(float64(ax)))
    sy, cy := float32(math.Sin(float64(ay))), float32(math.Cos(float64(ay)))
    sz, cz := float32(math.Sin(float64(az))), float32(math.Cos(float64(az)))

    rotate := func(px, py, pz float32) (float32, float32, float32) {
        // Rotate X
        y1 := py*cx - pz*sx
        z1 := py*sx + pz*cx
        x1 := px
        // Rotate Y
        x2 := x1*cy + z1*sy
        z2 := -x1*sy + z1*cy
        y2 := y1
        // Rotate Z
        x3 := x2*cz - y2*sz
        y3 := x2*sz + y2*cz
        z3 := z2
        return x3, y3, z3
    }
    
    v2 := textVPos(0)
    v3 := textVPos(1)

    // Precompute all rotated vertices once
    var rotated [6][4][3]float32
    for f := 0; f < 6; f++ {
        for i := 0; i < 4; i++ {
            rx, ry, rz := rotate(v[f][i][0]*size, v[f][i][1]*size, v[f][i][2]*size)
            rotated[f][i][0] = rx
            rotated[f][i][1] = ry
            rotated[f][i][2] = rz
        }
    }
    
    //~ n := num
    //~ i := 0
    //~ t := float32(int(math.Log10(float64(num))))
    //~ if t == 0 {
        //~ t = 1
    //~ }

    gl.Begin(gl.QUADS)
    
    //~ for {
        //~ m := n%10
        //~ const camside = 4
        //~ n0 := texHPos(uint32(6 + m + 1))
        //~ n1 := texHPos(uint32(6 + m))
        //~ gl.TexCoord2f(n1, v3); gl.Vertex3f(x+rotated[camside][0][0]/t-float32(i)*size, y+rotated[camside][0][1], z+rotated[camside][0][2]/t)
        //~ gl.TexCoord2f(n0, v3); gl.Vertex3f(x+rotated[camside][1][0]/t-float32(i)*size, y+rotated[camside][1][1], z+rotated[camside][1][2]/t)
        //~ gl.TexCoord2f(n0, v2); gl.Vertex3f(x+rotated[camside][2][0]/t-float32(i)*size, y+rotated[camside][2][1], z+rotated[camside][2][2]/t)
        //~ gl.TexCoord2f(n1, v2); gl.Vertex3f(x+rotated[camside][3][0]/t-float32(i)*size, y+rotated[camside][3][1], z+rotated[camside][3][2]/t)
        //~ n /= 10
        //~ i++
        //~ if n == 0 {
            //~ break
        //~ }
    //~ }

    for k, v := range str {
        const camside = 4

        var n0 float32
        var n1 float32
        
        if v == ' ' {
            continue
        } else if  (v >= '0') && (v <= '9') {
            //~ d(AIR, float32(FONT_0 + (v - '0')), float32(k)*s, y, 0)
            
            n0 = texHPos( uint32(FONT_0 + (v - '0')+1))
            n1 = texHPos( uint32(FONT_0 + (v - '0')))
            
        } else if (v >= 'a') && (v <= 'z')  {
            //~ d(AIR, float32(FONT_A + (v - 'a')), float32(k)*s - midpadh, y*s, 0)
            
            n0 = texHPos( uint32(FONT_A + (v - 'a')+1))
            n1 = texHPos( uint32(FONT_A + (v - 'a')))
        }
        
        kf := float32(1)
        gl.TexCoord2f(n1, v3); gl.Vertex3f(x+rotated[camside][0][0]/kf+float32(k)*size, y+rotated[camside][0][1], z+rotated[camside][0][2]/kf)
        gl.TexCoord2f(n0, v3); gl.Vertex3f(x+rotated[camside][1][0]/kf+float32(k)*size, y+rotated[camside][1][1], z+rotated[camside][1][2]/kf)
        gl.TexCoord2f(n0, v2); gl.Vertex3f(x+rotated[camside][2][0]/kf+float32(k)*size, y+rotated[camside][2][1], z+rotated[camside][2][2]/kf)
        gl.TexCoord2f(n1, v2); gl.Vertex3f(x+rotated[camside][3][0]/kf+float32(k)*size, y+rotated[camside][3][1], z+rotated[camside][3][2]/kf)
        
    }

    
    gl.End()
}

func LookAt(eye, center Vec3) {
    dot := func (a, b Vec3) float32 {
        return a.x*b.x + a.y*b.y + a.z*b.z
    }
    normalize := func (v Vec3) Vec3 {
        len := float32(math.Sqrt(float64(v.x*v.x + v.y*v.y + v.z*v.z)))
        return Vec3{v.x / len, v.y / len, v.z / len}
    }
    cross := func(a, b Vec3) Vec3 {
        return Vec3{
            a.y*b.z - a.z*b.y,
            a.z*b.x - a.x*b.z,
            a.x*b.y - a.y*b.x,
        }
    }
    up := Vec3{0, 1, 0}
    f := normalize(Vec3{
        x: center.x - eye.x,
        y: center.y - eye.y,
        z: center.z - eye.z,
    })
    s := normalize(cross(f, up))
    u := cross(s, f)
    view := [16]float32{
        s.x, u.x, -f.x, 0,
        s.y, u.y, -f.y, 0,
        s.z, u.z, -f.z, 0,
        -dot(s, eye), -dot(u, eye), dot(f, eye), 1,
    }
    gl.MultMatrixf(&view[0])
}

var selectedSlot = 0
var mouse_holding_right bool
var mouse_holding_left bool
var holding_q bool
var holding_p bool
var holding_l bool
var player_velocity Vec3
func playerInventoryInsert (id BlockID, c int32) {
    if id == AIR {
        return
    }
    for k, v := range player_inventory.ID {
        if v == id {
            player_inventory.C[k] += c
            return
        }
    }
    player_inventory.ID = append(player_inventory.ID, id)
    player_inventory.C = append(player_inventory.C, c)
}

func inventoryClearEmpty(inv *Inventory) {
    nID := inv.ID[:0]
    nC := inv.C[:0]
    for k, v := range inv.C {
        if v > 0 {
            nID = append(nID, inv.ID[k])
            nC = append(nC, v)
        }
    }
    inv.C = nC
    inv.ID = nID
}

func AABB(x, y, z, dx, dy, dz, height, fat float32, callback func(*BlockID, int,int,int) ) (bool, bool, bool, int) {
    const bb_range = 3
    var npbb BoundingBox
    var npbox BoundingBox
    canWalkThrough := [N_BLOCK_IDS]bool{
        AIR: true,
        LAVA: true,
        WATER: true,
        FLOWER: true,
    }
    vertical_collision_type := 0
    var bx, by, bz bool = true, true, true
    for i := -bb_range; i < bb_range; i++ {
        for j := -bb_range; j < bb_range; j++ {
            for k := -bb_range; k < bb_range; k++ {
                // collision through x
                npbox.X = x + dx - fat
                npbox.Y = y - height
                npbox.Z = z - fat
                npbox.W = 2*fat
                npbox.H = height
                npbox.L = 2*fat
                npbb.X = float32(int(x + dx)+i)
                npbb.Y = float32(int(y     )+j)
                npbb.Z = float32(int(z     )+k)
                npbb.W = 1.0 + fat
                npbb.H = 1.0
                npbb.L = 1.0 + fat
                if chunkExists( int(x + dx)+i, int(y)+j, int(z)+k ) {
                    b := getBlockRef(int(x+dx)+i, int(y)+j, int(z)+k)
                    if !canWalkThrough[*b] {
                        callback(b, int(x+dx)+i, int(y)+j, int(z)+k)
                        if intersect(npbox, npbb) {
                            bx = false
                        }
                    }
                }
                
                // collision through y
                npbox.X = x - fat
                npbox.Y = y + dy - height
                npbox.Z = z - fat
                npbox.W = 2*fat
                npbox.H = height
                npbox.L = 2*fat
                npbb.X = float32(int(x     )+i)
                npbb.Y = float32(int(y + dy)+j)
                npbb.Z = float32(int(z     )+k)
                npbb.W = 1.0 + fat
                npbb.H = 1.0
                npbb.L = 1.0 + fat
                if chunkExists( int(x)+i, int(y + dy)+j, int(z)+k ) {
                    b := getBlockRef(int(x)+i, int(y + dy)+j, int(z)+k)
                    if !canWalkThrough[*b] {
                        callback(b, int(x)+i, int(y + dy)+j, int(z)+k)
                        if intersect(npbox, npbb) {
                            if j < 0 {
                                vertical_collision_type = 1
                            } else {
                                vertical_collision_type = -1
                            }
                            by = false 
                        }
                    }
                }
                
                // collision through z
                npbox.X = x - fat
                npbox.Y = y - height
                npbox.Z = z + dz - fat
                npbox.W = 2*fat
                npbox.H = height
                npbox.L = 2*fat
                npbb.X = float32(int(x     )+i)
                npbb.Y = float32(int(y     )+j)
                npbb.Z = float32(int(z + dz)+k)
                npbb.W = 1.0 + fat
                npbb.H = 1.0 
                npbb.L = 1.0 + fat
                if chunkExists( int(x)+i, int(y)+j, int(z + dz)+k ) {
                    b := getBlockRef(int(x)+i, int(y)+j, int(z+dz)+k)
                    if !canWalkThrough[*b] {
                        callback(b, int(x)+i, int(y)+j, int(z+dz)+k)
                        if intersect(npbox, npbb) {
                            bz = false
                        }
                    }
                }
                
                // collision through xz
                npbox.X = x - fat
                npbox.Y = y - height
                npbox.Z = z + dz - fat
                npbox.W = 2*fat
                npbox.H = height
                npbox.L = 2*fat
                npbb.X = float32(int(x + dx)+i)
                npbb.Y = float32(int(y     )+j)
                npbb.Z = float32(int(z + dz)+k)
                npbb.W = 1.0 + fat
                npbb.H = 1.0 
                npbb.L = 1.0 + fat
                if chunkExists( int(x + dx)+i, int(y)+j, int(z + dz)+k ) {
                    b := getBlockRef(int(x+dx)+i, int(y)+j, int(z+dz)+k)
                    if !canWalkThrough[*b] {
                        callback(b, int(x+dx)+i, int(y)+j, int(z+dz)+k)
                        if intersect(npbox, npbb) {
                            bx = false
                            bz = false
                        }
                    }
                }
            }
        }
    }
    return bx, by, bz, vertical_collision_type
}

const bob_speed =  2.0
//~ const bob_range = 0.005
const bob_range = 0.01
var walking_bob float32


var holding_period bool
var holding_comma bool


func setBlockText(x, y, z int, s []byte) {
    c := World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}]
    
    if c.blockText == nil {
        c.blockText = make(map[nearVec][]byte)
    }
    c.blockText[nearVec{uint8(x%SUBCHUNK_H), uint8(y), uint8(z%SUBCHUNK_H)}] = s
}

func getBlockText(x, y, z int) []byte {
    c := World[iVec2{x-x%SUBCHUNK_H,z-z%SUBCHUNK_H}]
    v, e := c.blockText[nearVec{uint8(x%SUBCHUNK_H), uint8(y), uint8(z%SUBCHUNK_H)}]
    if e {
        return v
    } else {
        return []byte("")
    }
}


/* this used to be just about player collision... */
func playerCollision(window *glfw.Window) {



    if (window.GetKey(glfw.KeyPeriod) == glfw.Press) && !holding_period {
        drop_amount++
        holding_period = true
    }
    if (window.GetKey(glfw.KeyPeriod) == glfw.Release) && holding_period {
        
        holding_period = false
    }
    
    if (window.GetKey(glfw.KeyComma) == glfw.Press) && !holding_comma {
        drop_amount--
        if drop_amount <= 0 {
            drop_amount = 1
        }
        holding_comma = true
    }
    if (window.GetKey(glfw.KeyComma) == glfw.Release) && holding_comma {
        holding_comma = false
    }

    const speed = 0.14*0.90
    dir = Vec3{
        x: cosf(cam_yaw) * cosf(cam_pitch),
        y: sinf(cam_pitch),
        z: sinf(cam_yaw) * cosf(cam_pitch),
    }
    
    np := Vec3{player_velocity.x,player_velocity.y,player_velocity.z}
    
    
    if window.GetKey(glfw.KeyLeftControl) == glfw.Press {
        /* for now the handling of key ctrl stays on main */
    } else {
        if window.GetKey(glfw.KeyW) == glfw.Press {
            np.x += cosf(cam_yaw) * speed
            np.z += sinf(cam_yaw) * speed
        }
        
        if window.GetKey(glfw.KeyD) == glfw.Press {
            np.x += cosf(cam_yaw+1.57) * (speed/1.2)
            np.z += sinf(cam_yaw+1.57) * (speed/1.2)
        }
        if window.GetKey(glfw.KeyA) == glfw.Press {
            np.x += cosf(cam_yaw-1.57) * (speed/1.2)
            np.z += sinf(cam_yaw-1.57) * (speed/1.2)
        }
        
        if window.GetKey(glfw.KeyS) == glfw.Press {
            np.x += -cosf(cam_yaw) * speed
            np.z += -sinf(cam_yaw) * speed
        }
    }
    
    
    if (window.GetKey(glfw.KeyQ) == glfw.Press) && !holding_q {
        if (len(player_inventory.C) > 0) {
            min := func (a, b int32) int32 {
                if a > b { return b }
                return a
            }
            
            itemDrops = append(itemDrops, Drop{
                c: min(int32(drop_amount), player_inventory.C[selectedSlot]),
                dir: dir,
                id: player_inventory.ID[selectedSlot],
                pos: player_pos,
                rot: 0,
            })
            player_inventory.C[selectedSlot]-= min(int32(drop_amount), player_inventory.C[selectedSlot])
            if player_inventory.C[selectedSlot] == -1 {
                selectedSlot = 0
            }
            
            inventoryClearEmpty(&player_inventory)
        }
        holding_q = true
    }
    
    if (window.GetKey(glfw.KeyQ) == glfw.Release) && holding_q {
        holding_q = false
    }
    
    if (window.GetKey(glfw.KeyP) == glfw.Press) && !holding_p {
        saveWritePlayerdata()
        for k, v := range World {
            saveWriteChunk(v, k)
        }
        holding_p = true
    }
    if (window.GetKey(glfw.KeyP) == glfw.Release) && holding_p {
        holding_p = false
    }    

    if (window.GetKey(glfw.KeyL) == glfw.Press) && !holding_l {
        saveReadPlayerdata()
        holding_l = true
    }
    if (window.GetKey(glfw.KeyL) == glfw.Release) && holding_l {
        holding_l = false
    }
    
    vertical_collision_type := 0
    var valid_position_x, valid_position_y, valid_position_z = true, true, true
    var underwater bool = false
    
    player_velocity.y -= 0.025
    if chunkExists( int(player_pos.x), int(player_pos.y-1), int(player_pos.z) ) {
        if getBlockVal(int(player_pos.x),int(player_pos.y-1),int(player_pos.z)) == WATER {
            np.y += 0.15
            underwater = true
        }
    }


    player_aabb_callback := func( b *BlockID, x, y, z int) {
        if *b == OBSCURE {
            *b = AIR
        }
        //~ chunkUpdateDisplaylist(x, y, z)
    }
    
    valid_position_x, valid_position_y, valid_position_z, vertical_collision_type = AABB(
        player_pos.x,
        player_pos.y,
        player_pos.z,
        np.x,
        np.y,
        np.z,
        1.5,
        0.125, player_aabb_callback) /* previous fat: 0.25 */
        
    if (valid_position_x || valid_position_z) && (!valid_position_y) {
        abs := func(f float32) float32 {
            if f < 0 {
                return -f
            }
            return f
        }
    
        if (window.GetKey(glfw.KeyW) == glfw.Press) ||
        (window.GetKey(glfw.KeyA) == glfw.Press) ||
        (window.GetKey(glfw.KeyS) == glfw.Press) ||
        (window.GetKey(glfw.KeyD) == glfw.Press) {
        }
        
        if abs(np.x) > abs(np.z) {
            walking_bob += bob_speed*abs(np.x)
        } else {
            walking_bob += bob_speed*abs(np.z)
        }
    }
    
    player_velocity = Vec3{0, player_velocity.y*0.975, 0}
    
    if valid_position_x {
        player_pos.x += np.x
    }
    if valid_position_y {
        player_pos.y += np.y
    } else {
        player_velocity.y = -0.001
    }
    if window.GetKey(glfw.KeySpace) == glfw.Press {
        if underwater {
            player_velocity.y += 0.30
        }
        if vertical_collision_type == 1 {
            player_velocity.y = 0.30
        }
    }
    if valid_position_z {
        player_pos.z += np.z
    }
    
    //~ gl.Begin(gl.QUADS)
    
    for i := 0; i < 20; i++ {
        
        v := Vec3{
            x: player_pos.x + dir.x * float32(i)/4.0,
            y: player_pos.y + dir.y * float32(i)/4.0,
            z: player_pos.z + dir.z * float32(i)/4.0,
        }
        
        
        var entityInteraction = func (ek int) {
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press) && !mouse_holding_right {
                mouse_holding_right = true  
            }
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press) {
                //~ log.Println("teste\n")
                if Entities[ek].id == MOBID_BEE {
                    playerInventoryInsert(WAX, 1)
                }
                Entities[ek].hp = -1.0
            }
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Release) && mouse_holding_right {
                mouse_holding_right = false
            }
        }
        var blockInteraction = func () {
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press) && !mouse_holding_right {
                mouse_holding_right = true  
            }
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press) {
                new := true
                blockid := int(getBlockVal(int(v.x), int(v.y), int(v.z)))
                for bk, bv := range breakingBlocks {
                    if (bv.x == int(v.x)) && (bv.y == int(v.y)) && (bv.z == int(v.z)) {

                        if len(player_inventory.ID) > 0 {
                            //~ if (player_inventory.ID[selectedSlot] == PICKAXE) && ( (blockid == STONE) || (blockid == COBBLE) || (blockid == FLESH) || (blockid == ORE_COAL) || (blockid == ORE_IRON) ) {
                                //~ breakingBlocks[bk].t+=2
                            //~ } else {
                                //~ breakingBlocks[bk].t++
                            //~ }
                            
                            if (player_inventory.ID[selectedSlot] == PICKAXE) {
                                switch blockid {
                                case BRICK, FURNACE, STONE, COBBLE, FLESH, ORE_COAL, ORE_IRON:
                                    breakingBlocks[bk].t+=2
                                default:
                                    breakingBlocks[bk].t++
                                }
                                
                            } else if (player_inventory.ID[selectedSlot] == SHOVEL) {
                                switch blockid {
                                case DIRT, MOON, CLAY, SAND:
                                    breakingBlocks[bk].t+=2
                                default:
                                    breakingBlocks[bk].t++
                                }
                                
                            } else if (player_inventory.ID[selectedSlot] == AXE) {
                                switch blockid {
                                case SIGN, WOOD, PLANK, CRAFTING_TABLE:
                                    breakingBlocks[bk].t+=2
                                default:
                                    breakingBlocks[bk].t++
                                }
                            } else {
                                breakingBlocks[bk].t++
                                
                            }
                            
                        } else {
                            breakingBlocks[bk].t++
                        }
                        

                        breakingBlocks[bk].elapsed = 0
                        new = false
                        break
                    }
                }
                
                
                for ek, ev := range Entities {
                    if (int(ev.pos.x) == int(v.x)) && (int(ev.pos.y) == int(v.y)) && (int(ev.pos.z) == int(v.z)) {
                        //~ log.Println("teste\n")
                        Entities[ek].hp = -1.0
                    }
                }
                
                
                if new {
                    breakingBlocks = append(breakingBlocks, struct{x, y, z, t, id, elapsed int}{ int(v.x), int(v.y), int(v.z), 0, blockid, int(0) })
                }
            }
            if (window.GetMouseButton(glfw.MouseButtonRight) == glfw.Release) && mouse_holding_right {
                mouse_holding_right = false
            }
            
            
            if (window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press) && !mouse_holding_left {
                
                ClosestToInt := func (a, b, c float32) float32 {
                    dist := func(x float32) float32 {
                        return float32(math.Abs(float64(x - float32(math.Round(float64(x))))))
                    }
                    dA := dist(a)
                    dB := dist(b)
                    dC := dist(c)
                    closest := a
                    minDist := dA
                    if dB < minDist {
                        closest = b
                        minDist = dB
                    }
                    if dC < minDist {
                        closest = c
                    }
                    return closest
                }
                
                norm := func (f float32) float32 {
                    return f - float32(math.Round(float64(f)))
                }
                
                var c float32
                c = ClosestToInt(v.x,v.y,v.z)
                d := [3]int{0,0,0}
                
                if c == v.x {
                    if norm(c) > 0 {
                        d = [3]int{-1,0,0}
                    } else {
                        d = [3]int{1,0,0}
                    }
                }
                if c == v.z {
                    if norm(c) > 0 {
                        d = [3]int{0,0,-1}
                    } else {
                        d = [3]int{0,0,1}
                    }
                }
                if c == v.y {
                    if norm(c) > 0 {
                        d = [3]int{0,-1,0}
                    } else {
                        d = [3]int{0,1,0}
                    }
                }
                b := getBlockRef(int(v.x)+d[0],int(v.y)+d[1],int(v.z)+d[2])
                
                
                if *b == AIR {
                    if (len(player_inventory.C)>0) {
                        *b = player_inventory.ID[selectedSlot]
                        updateSpread(int(v.x  ), int(v.y  ), int(v.z  ))
                        chunkUpdateDisplaylist(int(v.x)+d[0], int(v.y)+d[1], int(v.z)+d[2])
                    
                        player_inventory.C[selectedSlot]--
                        inventoryClearEmpty(&player_inventory)
                        
                        /* special block placing actions... */
                        if player_inventory.ID[selectedSlot] == SIGN {
                            setBlockText(int(v.x)+d[0],int(v.y)+d[1],int(v.z)+d[2], player_textbuffer)
                            player_textbuffer = player_textbuffer[:0]
                        }

                        if selectedSlot >= len(player_inventory.ID) {
                           selectedSlot =  len(player_inventory.ID) - 1
                        }
                        if selectedSlot < 0 {
                            selectedSlot = 0
                        }
                    }
                }
                mouse_holding_left = true
            }
            if (window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Release) && mouse_holding_left {
                mouse_holding_left = false
            }
        }
        
        if chunkExists(int(v.x),int(v.y),int(v.z)) {
            b := getBlockVal(int(v.x),int(v.y),int(v.z))
            if i <= 2 {
                if b == WATER {
                    drawBillboard(WATER, 0, v.x, v.y, v.z, 2.0)
                }
                if b == LAVA {
                    drawBillboard(LAVA, 0, v.x, v.y, v.z, 2.0)
                }
            }
            
            for ek, ev := range Entities {
                if (int(ev.pos.x) == int(v.x)) && (int(ev.pos.y) == int(v.y)) && (int(ev.pos.z) == int(v.z)) {
                    entityInteraction(ek)
                    goto exit_raycast
                }
            }
            if (b != AIR) && (b != WATER) && (b != LAVA) {
                blockInteraction()
                goto exit_raycast
            }
        }
    }
    exit_raycast:
}

func cosf(f float32) float32 {
    return float32(math.Cos(float64(f)))
}
func sinf(f float32) float32 {
    return float32(math.Sin(float64(f)))
}

type BoundingBox struct {
    X, Y, Z float32
    W, H, L float32
}
func intersect(a, b BoundingBox) bool {
    aminX, aminY, aminZ := a.X, a.Y, a.Z
    amaxX, amaxY, amaxZ := a.X+a.W, a.Y+a.H, a.Z+a.L
    bminX, bminY, bminZ := b.X, b.Y, b.Z
    bmaxX, bmaxY, bmaxZ := b.X+b.W, b.Y+b.H, b.Z+b.L
    return aminX <= bmaxX && amaxX >= bminX &&
           aminY <= bmaxY && amaxY >= bminY &&
           aminZ <= bmaxZ && amaxZ >= bminZ
}

var saveSlot = "saves/World1"
func saveWriteChunk(c *Chunk, p iVec2) {
    filename := fmt.Sprintf("%s/%d,%d.bin", saveSlot, p.x, p.y)

    file, err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    writer, err := flate.NewWriter(file, flate.BestSpeed)
    if err != nil {
        log.Fatal(err)
    }
    defer writer.Close() // flush compressed data on function exit
    
    for i := 0; i < SUBCHUNK_H; i++ {
        for j := 0; j < CHUNK_V; j++ {
            if err := binary.Write(writer, binary.LittleEndian, c.block[i][j] );  err != nil {
                log.Fatal(err)
            }
        }
    }
    for i := 0; i < SUBCHUNK_H; i++ {
        if err := binary.Write(writer, binary.LittleEndian, c.cloud[i] );  err != nil {
            log.Fatal(err)
        }
    }

    //~ binary.Write(writer, binary.LittleEndian, int32(len(c.blockInventories)))
    //~ binary.Write(writer, binary.LittleEndian, c.blockInventories)
    writer.Close()
    file.Close()    
}

func saveWritePlayerdata() {
    filename := fmt.Sprintf("%s/level.bin",saveSlot)
    file, err := os.Create(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    binary.Write(file, binary.LittleEndian, player_pos.x )
    binary.Write(file, binary.LittleEndian, player_pos.y )
    binary.Write(file, binary.LittleEndian, player_pos.z )
    binary.Write(file, binary.LittleEndian, cam_pitch )
    binary.Write(file, binary.LittleEndian, cam_yaw )
    binary.Write(file, binary.LittleEndian, TICK )
    if err = binary.Write(file, binary.LittleEndian, int32(len(player_inventory.C)) ); err != nil { log.Fatal(err) }
    if err = binary.Write(file, binary.LittleEndian, player_inventory.C ); err != nil { log.Fatal(err) }
    if err = binary.Write(file, binary.LittleEndian, player_inventory.ID ); err != nil { log.Fatal(err) }
    file.Close()
}


func saveReadChunk(c *Chunk, p iVec2) *Chunk {
    filename := fmt.Sprintf("%s/%d,%d.bin",saveSlot,p.x, p.y)
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    reader := flate.NewReader(file)
    defer reader.Close()
    for i := 0; i < SUBCHUNK_H; i++ {
        for j := 0; j < CHUNK_V; j++ {
            if err := binary.Read(reader, binary.LittleEndian, &c.block[i][j]);  err != nil {
                log.Fatal(err)
            }
        }
    }
    for i := 0; i < SUBCHUNK_H; i++ {
        if err := binary.Read(reader, binary.LittleEndian, &c.cloud[i]);  err != nil {
            log.Fatal(err)
        }
    }
    file.Close()
    reader.Close()
    
    ChunkUpdateDisplayListRef(c)
    return c
}

func saveReadPlayerdata() {
    filename := fmt.Sprintf("%s/level.bin", saveSlot)
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    binary.Read(file, binary.LittleEndian, &player_pos.x);
    binary.Read(file, binary.LittleEndian, &player_pos.y);
    binary.Read(file, binary.LittleEndian, &player_pos.z);
    binary.Read(file, binary.LittleEndian, &cam_pitch);
    binary.Read(file, binary.LittleEndian, &cam_yaw);
    binary.Read(file, binary.LittleEndian, &TICK);
    var length int32
    if err = binary.Read(file, binary.LittleEndian,  &length ); err != nil { log.Fatal(err) }
    player_inventory.C = player_inventory.C[:0]
    player_inventory.ID = player_inventory.ID[:0]
    for i := int32(0); i < length; i++ {
        var c int32
        if err = binary.Read(file, binary.LittleEndian,  &c ); err != nil { log.Fatal(err) }
        player_inventory.C = append(player_inventory.C, c)
    }
    for i := int32(0); i < length; i++ {
        var id BlockID
        if err = binary.Read(file, binary.LittleEndian,  &id ); err != nil { log.Fatal(err) }
        player_inventory.ID = append(player_inventory.ID, id)
    }
    file.Close()
}
