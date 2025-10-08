package arbor

import (
	"fmt"
	"mock-arbor/src/util"
	"strings"
)

type policy struct {
	Name            string
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
		if cm == "" {
			cm = "(none)"
		}
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
