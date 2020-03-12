
package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "os"
    "encoding/csv"
    "encoding/hex"
    "log"
    "strconv"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
    //"io/ioutil"
)

const (
    //host     = "172.18.0.3"
    host     = "repono"
    user     = "postgres"
    password = "example"
    dbname   = "palettepal"
)

type RGB struct {
    R, G, B uint8
}

type Specimen struct {
    PhaseA, PhaseB *[16]uint8
    Ensemble string
    Score map[string]int
}

var log_file *os.File
var gerr error
var filter = make(map[uint16]bool)
var filter_desc string
var target_desc string
var worker_count int64
var rounds_total int64
var run_id string

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

func load_filter_csv(file string, rank_edict string) {
    // parse the given csv file and populate the filter
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
                filter[uint16(cc)] = true
            }
        }  // cc: color code, rr: render rank
    }
}

func pick_phase_pair() ([16]uint8, [16]uint8) {
    // populate the phase palette sets by randomly selecting
    // values found in the unsat_v6 master palette
    phaseA := [16]uint8{}
    phaseB := [16]uint8{}

    phaseA[0] = uint8(rand.Intn(64))
    phaseA[1] = uint8(rand.Intn(64))
    phaseA[2] = uint8(rand.Intn(64))
    phaseA[3] = uint8(rand.Intn(64))
    phaseA[4] = phaseA[0]
    phaseA[5] = uint8(rand.Intn(64))
    phaseA[6] = uint8(rand.Intn(64))
    phaseA[7] = uint8(rand.Intn(64))
    phaseA[8] = phaseA[0]
    phaseA[9] = uint8(rand.Intn(64))
    phaseA[10] = uint8(rand.Intn(64))
    phaseA[11] = uint8(rand.Intn(64))
    phaseA[12] = phaseA[0]
    phaseA[13] = uint8(rand.Intn(64))
    phaseA[14] = uint8(rand.Intn(64))
    phaseA[15] = uint8(rand.Intn(64))

    phaseB[0] = uint8(rand.Intn(64))
    phaseB[1] = uint8(rand.Intn(64))
    phaseB[2] = uint8(rand.Intn(64))
    phaseB[3] = uint8(rand.Intn(64))
    phaseB[4] = phaseB[0]
    phaseB[5] = uint8(rand.Intn(64))
    phaseB[6] = uint8(rand.Intn(64))
    phaseB[7] = uint8(rand.Intn(64))
    phaseB[8] = phaseB[0]
    phaseB[9] = uint8(rand.Intn(64))
    phaseB[10] = uint8(rand.Intn(64))
    phaseB[11] = uint8(rand.Intn(64))
    phaseB[12] = phaseB[0]
    phaseB[13] = uint8(rand.Intn(64))
    phaseB[14] = uint8(rand.Intn(64))
    phaseB[15] = uint8(rand.Intn(64))

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

func store_specimen(specimen *Specimen) {
    connStr := fmt.Sprintf(
        "dbname=" + dbname +
            " host=" + host + " sslmode=disable" +
            " user=" + user + " password=" + password)
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    statement := `
    INSERT INTO mc_result (
        color_count,
        max_vp_size,
        filter_desc,
        target_desc,
        ensemble,
        run_id,
        worker_id,
        duration,
        rounds
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id`

    id := 0
    if err = db.QueryRow(statement,
        specimen.Score["color_count"],
        specimen.Score["max_vp_size"],
        filter_desc,
        target_desc,
        specimen.Ensemble,
        run_id,
        "testbed_1",
        20,
        50,
    ).Scan(&id); err != nil {
        panic(err)
    }
    fmt.Println("specimen saved, record_id: ", id)
}

func yield_specimen(a, b *[16]uint8) *Specimen {
    // initialize a specimen wrapper given a phase pair
    specimen := Specimen{}
    specimen.PhaseA = a
    specimen.PhaseB = b
    specimen.Score = make(map[string]int)

    ensemble := fmt.Sprintf("%s", strings.ToUpper(hex.EncodeToString([]byte( (*a)[:] ))))
    ensemble += fmt.Sprintf("%s", strings.ToUpper(hex.EncodeToString([]byte( (*b)[:] ))))

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

    specimen.Score["color_count"] = 0
    specimen.Score["contrast_amt"] = 0
    specimen.Score["max_vp_size"] = 0
    unique_colors := make(map[uint16]bool)

    vps_full := yield_vps_full(specimen.PhaseA, specimen.PhaseB)
    for _, cc := range vps_full {
        if filter[cc] && !unique_colors[cc] {
            specimen.Score["color_count"] += 1
        }
        unique_colors[cc] = true
    }
}

func adjudicate_specimen(specimen *Specimen) bool {
    // judge a specimen according to criteria specified
    // and determine whether it should be saved

    // if the number of available colors meets
    // a threshold, return affirmative
    if specimen.Score["color_count"] > 22 {
        return true
    }
    return false
}

func println_to_log(line_to_log string) {
    log_file.Write([]byte(line_to_log))
    log_file.Write([]byte("\n"))
}

func main() {
    // RTL: Improve error handling in logging functionality.
    log_file, gerr = os.OpenFile("montecolore.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if gerr != nil {
        log.Fatal(gerr)
    }
    println_to_log("Monte Colore process starting...")

    
    build_ultra()

    target_desc = "color_count > 22"

    rand.Seed(time.Now().UnixNano())

    bytes := make([]byte, 3)
    rand.Read(bytes)
    run_id = os.Args[1]
    rounds_total, _ = strconv.ParseInt(os.Args[2], 10, 64)
    worker_count, _ = strconv.ParseInt(os.Args[3], 10, 64)

    filter_rank_edict := os.Args[4]
    filter_desc = os.Args[5]
    load_filter_csv("filter.csv", filter_rank_edict)

    for i := int64(1); i < rounds_total; i++ {
        a, b := pick_phase_pair()
        x := yield_specimen(&a, &b)
        appraise_specimen(x)
        is_worthy := adjudicate_specimen(x)
        if is_worthy {
            print_phase_pair(x.PhaseA, x.PhaseB)
            store_specimen(x)
        }
    }
    fmt.Println()
    println_to_log("Monte Colore process terminating normally...")

}

