package arbor

import (
    "strconv"
    "strings"
    "testing"
    "time"
)

func TestMitigationSuite(t *testing.T) {
    if out, _ := Dispatch("show mitigation summary"); !strings.Contains(out, "Mitigation Summary") {
        t.Fatalf("summary missing: %q", out)
    }
    if out, _ := Dispatch("show mitigation detail 102"); !strings.Contains(out, "ID:            102") {
        t.Fatalf("detail missing id: %q", out)
    }
    if out, _ := Dispatch("show mitigation counters"); !strings.Contains(out, "Mitigation Counters") {
        t.Fatalf("counters missing: %q", out)
    }
    if out, _ := Dispatch("show mitigation statistics 101"); !strings.Contains(out, "ID: 101") {
        t.Fatalf("counters by id missing: %q", out)
    }
    if out, _ := Dispatch("mitigation list"); !strings.Contains(out, "Mitigation List") {
        t.Fatalf("list missing: %q", out)
    }
    if out, _ := Dispatch("clear mitigation counters"); !strings.Contains(out, "all") {
        t.Fatalf("clear all missing: %q", out)
    }
    if out, _ := Dispatch("clear mitigation counters 101"); !strings.Contains(out, "101") {
        t.Fatalf("clear id missing: %q", out)
    }
}

func TestGroupsDiversion(t *testing.T) {
    if out, _ := Dispatch("show tms groups"); !strings.Contains(out, "TMS Groups") {
        t.Fatalf("groups missing: %q", out)
    }
    if out, _ := Dispatch("show tms group edge1"); !strings.Contains(out, "TMS Group \"edge1\"") {
        t.Fatalf("group detail missing: %q", out)
    }
    if out, _ := Dispatch("show diversion status"); !strings.Contains(out, "Diversion Status") {
        t.Fatalf("diversion status missing: %q", out)
    }
}

func TestPolicyCrudAndCatalog(t *testing.T) {
    // ensure empty
    out, _ := Dispatch("show policies")
    if !strings.Contains(out, "[none]") {
        t.Fatalf("expected no policies initially: %q", out)
    }
    // create + set countermeasures
    if out, _ = Dispatch("policy Alpha"); !strings.Contains(out, "ensured") {
        t.Fatalf("ensure failed: %q", out)
    }
    if out, _ = Dispatch("policy Alpha countermeasures drop fragments" ); !strings.Contains(out, "countermeasures set") {
        t.Fatalf("countermeasures failed: %q", out)
    }
    out, _ = Dispatch("show policies")
    if !strings.Contains(out, "Alpha") || !strings.Contains(out, "drop fragments") {
        t.Fatalf("catalog missing updates: %q", out)
    }
    // delete
    if out, _ = Dispatch("no policy Alpha"); !strings.Contains(out, "removed") {
        t.Fatalf("delete failed: %q", out)
    }
}

func TestSystemHealthAndNetworking(t *testing.T) {
    if out, _ := Dispatch("show system"); !strings.Contains(out, "System Health") {
        t.Fatalf("system missing: %q", out)
    }
    if out, _ := Dispatch("show processes"); !strings.Contains(out, "Processes") {
        t.Fatalf("processes missing: %q", out)
    }
    if out, _ := Dispatch("show license"); !strings.Contains(out, "License") {
        t.Fatalf("license missing: %q", out)
    }
    if out, _ := Dispatch("show interfaces brief"); !strings.Contains(out, "Interfaces (brief)") {
        t.Fatalf("interfaces brief missing: %q", out)
    }
    if out, _ := Dispatch("show bgp summary"); !strings.Contains(out, "BGP Summary") {
        t.Fatalf("bgp summary missing: %q", out)
    }
    if out, _ := Dispatch("show bgp neighbors"); !strings.Contains(out, "BGP Neighbors") {
        t.Fatalf("bgp neighbors missing: %q", out)
    }
    if out, _ := Dispatch("show flowspec status"); !strings.Contains(out, "FlowSpec Status") {
        t.Fatalf("flowspec status missing: %q", out)
    }
}

func TestUsabilityAndClockAndDiff(t *testing.T) {
    if out, _ := Dispatch("terminal length 0"); !strings.Contains(out, "Terminal length set") {
        t.Fatalf("terminal length set missing: %q", out)
    }
    out, _ := Dispatch("show clock")
    year := strconv.Itoa(time.Now().UTC().Year())
    if !strings.Contains(out, year) {
        t.Fatalf("clock missing year: %q", out)
    }
    if out, _ = Dispatch("show configuration differences"); !strings.Contains(out, "No configuration differences") {
        t.Fatalf("config diff missing: %q", out)
    }
}

