package arbor

import (
    "fmt"
    "mock-arbor/src/util"
    "strconv"
    "strings"
)

type policy struct {
    Name           string
    Countermeasures string
}

var (
    policies       = map[string]*policy{}
    terminalLength = 0
)

func formatPolicies() string {
    if len(policies) == 0 {
        return "Policies (emulated)\n  [none]\n\n"
    }
    out := "Policies (emulated)\nName               Countermeasures\n"
    for _, p := range policies {
        cm := p.Countermeasures
        if cm == "" { cm = "(none)" }
        out += fmt.Sprintf("%-18s %s\n", p.Name, cm)
    }
    return out + "\n"
}

func handlePolicy(raw string) string {
    // raw preserves original spacing/case for countermeasures payload
    norm := util.NormalizeSpaces(raw)
    if strings.HasPrefix(norm, "no policy ") {
        name := strings.TrimSpace(raw[len("no policy "):])
        if _, ok := policies[strings.ToLower(name)]; ok {
            delete(policies, strings.ToLower(name))
            return fmt.Sprintf("Policy %q removed (emulated).\n", name)
        }
        return fmt.Sprintf("Policy %q not found (emulated).\n", name)
    }
    // policy <name> [countermeasures ...]
    fields := strings.Fields(raw)
    if len(fields) >= 2 && strings.EqualFold(fields[0], "policy") {
        name := fields[1]
        key := strings.ToLower(name)
        p := policies[key]
        if p == nil {
            p = &policy{Name: name}
            policies[key] = p
        }
        // optional countermeasures
        lower := strings.ToLower(raw)
        idx := strings.Index(lower, "countermeasures ")
        if idx >= 0 {
            cm := strings.TrimSpace(raw[idx+len("countermeasures "):])
            p.Countermeasures = cm
            return getFixture("policy/countermeasures", map[string]string{"NAME": name, "CM": cm}, fmt.Sprintf("Policy %q countermeasures set to %q (emulated).\n", name, cm))
        }
        return getFixture("policy/ensure", map[string]string{"NAME": name}, fmt.Sprintf("Policy %q ensured (emulated).\n", name))
    }
    return "Invalid policy command (emulated).\n"
}

func handleTerminal(cmd string) string {
    // terminal length <n>
    parts := strings.Fields(cmd)
    if len(parts) != 3 {
        return "Usage: terminal length <n>\n"
    }
    n, err := strconv.Atoi(parts[2])
    if err != nil || n < 0 {
        return "Invalid terminal length\n"
    }
    terminalLength = n
    return getFixture("terminal/length", map[string]int{"N": terminalLength}, fmt.Sprintf("Terminal length set to %d (emulated).\n", terminalLength))
}

func mitigationList() string {
    return `Mitigation List (emulated)
ID    Name        State
101   ddos-icmp   IDLE
102   ddos-syn    ACTIVE

`
}

func clearMitigationCounters(cmd string) string {
    // clear mitigation counters [<id>]
    fields := strings.Fields(cmd)
    if len(fields) == 3 {
        return getFixture("mitigation/clear_all", nil, "Mitigation counters cleared (all) (emulated).\n")
    }
    id := fields[3]
    return getFixture("mitigation/clear_id", map[string]string{"ID": id}, fmt.Sprintf("Mitigation counters cleared for %s (emulated).\n", id))
}
