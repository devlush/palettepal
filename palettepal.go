
package main

import (
    "fmt"
    "math"
)

type RGB struct {
    R, G, B uint8
}

var palette_master = []RGB {
    {107, 107, 107}, {135,  30,   0}, {150,  11,  31}, {135,  12,  59},
    { 97,  13,  89}, { 40,   5,  94}, {  0,  17,  85}, {  0,  27,  70},
    {  0,  50,  48}, {  0,  72,  10}, {  0,  78,   0}, { 25,  70,   0},
    { 88,  58,   0}, {  0,   0,   0}, {  0,   0,   0}, {  0,   0,   0},

    {178, 178, 178}, {209,  83,  26}, {238,  53,  72}, {236,  35, 113},
    {183,  30, 154}, { 98,  30, 165}, { 25,  45, 165}, {  0,  75, 135},
    {  0, 105, 103}, {  0, 132,  41}, {  0, 139,   3}, { 64, 130,   0},
    {145, 120,   0}, {  0,   0,   0}, {  0,   0,   0}, {  0,   0,   0},

    {255, 255, 255}, {253, 173,  99}, {254, 138, 144}, {252, 119, 185},
    {254, 113, 231}, {201, 111, 247}, {106, 131, 245}, { 41, 156, 221},
    {  7, 184, 189}, {  7, 209, 132}, { 59, 220,  91}, {125, 215,  72},
    {206, 204,  72}, { 85,  85,  85}, {  0,   0,   0}, {  0,   0,   0},

    {255, 255, 255}, {254, 227, 196}, {254, 213, 215}, {254, 205, 230},
    {254, 202, 249}, {240, 201, 254}, {199, 209, 254}, {172, 220, 247},
    {156, 232, 232}, {157, 242, 209}, {177, 244, 191}, {205, 245, 183},
    {238, 240, 183}, {190, 190, 190}, {  0,   0,   0}, {  0,   0,   0} }

var palette_ultra = [64][64]RGB {}

func build_ultra() {
    for i, ci := range palette_master {
        for j, cj := range palette_master {
            palette_ultra[i][j] = blend(ci, cj)
        }
    }
}

func rms(x, y uint8) uint8 {
    xf := float64(x)
    yf := float64(y)

    xfs := math.Pow(xf, 2)
    yfs := math.Pow(yf, 2)

    res := math.Sqrt( (xfs + yfs ) / 2 )
    return uint8( math.Round(res) )
}

func blend(p, q RGB) RGB {
    v := RGB{}
    v.R = rms(p.R, q.R)
    v.G = rms(p.G, q.G)
    v.B = rms(p.B, q.B)
    return v
}

func print_master() {
    fmt.Println(palette_master[0x11])
    for i := 0; i < 4; i++ {
        for j := 0; j < 16; j++ {
            fmt.Print(palette_master[i*16+j])
            fmt.Print(" ")
        }
        fmt.Println()
    }
}

func main() {

    pm := &palette_master
    t1 := blend( (*pm)[0x11], (*pm)[0x12] )
    t2 := blend( (*pm)[0x28], (*pm)[0x03] )

    build_ultra()

    fmt.Println(palette_master[0x11].B)
    fmt.Println(t1)
    fmt.Println(t2)

    fmt.Println(palette_ultra[0x11][0x12])
    fmt.Println(palette_ultra[0x28][0x03])
}

