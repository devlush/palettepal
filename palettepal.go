
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
    "strings"
)

type RGB struct {
    R, G, B uint8
}

type Specimen struct {
    PhaseA, PhaseB *[16]uint8
    Ensemble string
    Score map[string]int
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

func load_sieve_csv(file string, rank_edict string) {
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

    rank_map := make(map[string]uint8)
    rank_map["green"] = 0x04
    rank_map["yellow"] = 0x06

    rew := rank_map[rank_edict]  // rew: rank edict weight

    for _, record := range records {
        if cc, err := strconv.ParseUint(record[0], 16, 16); err == nil {
            if rr := strings.ToLower(record[1]); rank_map[rr] <= rew {
                sieve[uint16(cc)] = true
            }
        }  // cc: color code, rr: render rank
    }
}

func pick_phase_pair() ([16]uint8, [16]uint8) {
    // populate the phase palette sets by randomly selecting
    // values found in the unsat_v6 master palette
    phaseA := [16]uint8{}
    phaseB := [16]uint8{}
    for i := range phaseA {
        phaseA[i] = uint8(rand.Intn(64))
    }
    for i := range phaseB {
        phaseB[i] = uint8(rand.Intn(64))
    }
    return phaseA, phaseB
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

func yield_specimen(a, b *[16]uint8) *Specimen {
    // initialize a specimen wrapper given a phase pair
    specimen := Specimen{}
    specimen.PhaseA = a
    specimen.PhaseB = b
    specimen.Score = make(map[string]int)

    var ensemble string = ""
    for _, cc := range *a {
        ensemble += fmt.Sprintf("%02x", cc)
    }

    for _, cc := range *b {
        ensemble += fmt.Sprintf("%02x", cc)
    }

    specimen.Ensemble = ensemble
    return &specimen
}

func yield_vps_full(a, b *[16]uint8) [256]uint16 {
    // calculate the full product of a phase pair and
    // let the 16x16 result be known as a 'virtual palette set'
    vps := [256]uint16 {}
    for i := 0; i < 16; i++ {
        for j := 0; j < 16; j++ {
            left := uint16((*a)[i]) << 8
            right := uint16((*b)[j])
            vps[i*16+j] = left | right
        }
    }
    return vps
}

func print_vps(vps *[256]uint16) {
    // helper function for printing the color codes
    // composing a virtual palette set (16x16)
    fmt.Printf("\n vps:")
    for i := 0; i < 16; i++ {
        for j := 0; j < 16; j++ {
            fmt.Printf(" %04x", vps[i*16+j])
        }
        fmt.Printf("\n     ")
    }
    fmt.Println("\n")
}

func print_phase_pair(a, b *[16]uint8) {
    // helper function for printing the color codes
    // composing a phase pair (A 4x4, B 4x4)
    fmt.Printf("\n A:")
    for j := 0; j < 4; j++ {
        fmt.Printf(" %02x", a[j])
    }
    fmt.Printf("   B:")
    for j := 0; j < 4; j++ {
        fmt.Printf(" %02x", b[j])
    }
    for i := 1; i < 4; i++ {
        fmt.Printf("\n   ")
        for j := 0; j < 4; j++ {
            fmt.Printf(" %02x", a[i*4+j])
        }
        fmt.Printf("     ")
        for j := 0; j < 4; j++ {
            fmt.Printf(" %02x", b[i*4+j])
        }
    }
    fmt.Println("\n")
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

func appraise_specimen(specimen *Specimen) {

    specimen.Score["colors_available"] = 0
    specimen.Score["contrast_amount"] = 0
    specimen.Score["largest_virtual_palette"] = 0

    vps_full := yield_vps_full(specimen.PhaseA, specimen.PhaseB)
    for _, cc := range vps_full {
        if sieve[cc] {
            specimen.Score["colors_available"] += 1
        }
    }
}

func adjudicate_specimen(specimen *Specimen) bool {
    // judge a specimen according to criteria specified
    // and determine whether it should be saved

    // if the number of available colors meets
    // a threshold, return affirmative
    if specimen.Score["colors_available"] > 22 {
        return true
    }
    return false
}

func main() {

    build_ultra()

    sieve[0x0921] = true
    sieve[0x0a05] = true

    load_sieve_csv("background.csv", "yellow")

    rand.Seed(time.Now().UnixNano())
    for i := 0; i < 10000; i++ {
        a, b := pick_phase_pair()
        x := yield_specimen(&a, &b)
        appraise_specimen(x)
        is_worthy := adjudicate_specimen(x)
        if is_worthy {
            print_phase_pair(x.PhaseA, x.PhaseB)
            fmt.Println(x.Score["colors_available"])
        }
    }
}

