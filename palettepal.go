
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "os"
    "encoding/csv"
    "log"
    "strconv"
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

func load_sieve_csv(file string) {
    // parse the given csv file and populate the sieve
    // with these colors values found in the ultra palette
    fd, err := os.Open(file)
    if err != nil {
        log.Fatalln("Unable to open csv file", err)
    }
    defer fd.Close()

    records, err := csv.NewReader(fd).ReadAll()
    if err != nil {
        log.Fatalln("Unable to parse csv file", err)
    }

    for _, record := range records {
        if cc, err := strconv.ParseUint(record[0], 16, 16); err == nil {
            sieve[uint16(cc)] = true
        }
    }
}

func pick_phase_pair() ([16]uint8, [16]uint8) {
    // populate the phase palette sets by randomly selecting
    // values found in the unsat_v6 master palette
    phaseA := [16]uint8{}
    phaseB := [16]uint8{}
    for _, i := range phaseA {
        phaseA[i] = uint8(rand.Intn(64))
    }
    for _, i := range phaseB {
        phaseB[i] = uint8(rand.Intn(64))
    }
    return phaseA, phaseB
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
    // calculate the root mean square of two 8bit integers
    xf := float64(x)
    yf := float64(y)

    xfs := math.Pow(xf, 2)
    yfs := math.Pow(yf, 2)

    res := math.Sqrt( (xfs + yfs ) / 2 )
    return uint8( math.Round(res) )
}

func blend(p, q RGB) RGB {
    // pairwise blend the RGB values for
    // two given colors using root mean square
    v := RGB{}
    v.R = rms(p.R, q.R)
    v.G = rms(p.G, q.G)
    v.B = rms(p.B, q.B)
    return v
}

func yield_vps(a, b [16]uint8) [256]uint16 {
    // calculate the full product of a phase pair and
    // let the 16x16 result be known as a 'virtual palette set'
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

    sieve[0x0921] = true
    sieve[0x0a05] = true

    load_sieve_csv("background.csv")

    rand.Seed(time.Now().UnixNano())
    for i := 0; i < 10000; i++ {
        a, b := pick_phase_pair()
        distill(a, b)
    }
}

