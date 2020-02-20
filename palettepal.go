
package main

import (
    "fmt"
    "math"
)

type RGB struct {
    R, G, B uint8
}

var sieve = make(map[uint16]bool)

var palette_unsat_v6 = []RGB {
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
    for i, ci := range palette_unsat_v6 {
        for j, cj := range palette_unsat_v6 {
            palette_ultra[i][j] = blend(ci, cj)
        }
    }
}

func distill(a, b [16]uint8) {

    vps := yield_vps(a, b)
    for _, cc := range vps {
        if sieve[cc] {
            fmt.Printf(" %04x\n", cc)
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

func yield_vps(a, b [16]uint8) [256]uint16 {

    vps := [256]uint16 {}
    for i := 0; i < 16; i++ {
        for j := 0; j < 16; j++ {
            left := uint16(a[i]) << 8
            right := uint16(b[j])
            vps[i*16+j] = left | right
        }
    }
    return vps
}

func print_master() {
    fmt.Println(palette_unsat_v6[0x11])
    for i := 0; i < 4; i++ {
        for j := 0; j < 16; j++ {
            fmt.Print(palette_unsat_v6[i*16+j])
            fmt.Print(" ")
        }
        fmt.Println()
    }
}

func main() {

    build_ultra()

    sieve[0x3afb] = true
    sieve[0x0aab] = true

    a := [16]uint8 {
        0x0a, 0x1a, 0x2a, 0x3a, 0x4a, 0x5a, 0x6a, 0x7a,
        0x8a, 0x9a, 0xaa, 0xba, 0xca, 0xda, 0xea, 0xfa }

    b := [16]uint8 {
        0x0b, 0x1b, 0x2b, 0x3b, 0x4b, 0x5b, 0x6b, 0x7b,
        0x8b, 0x9b, 0xab, 0xbb, 0xcb, 0xdb, 0xeb, 0xfb }

    distill(a, b)
}

