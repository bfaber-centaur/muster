package engine_test

import (
	"reflect"
	"testing"

	"github.com/bubskee/muster/engine"
)

func TestReplayProducesDeterministicTraceAndTerminalState(t *testing.T) {
	scenario := engine.Scenario{
		ID: "test-replay",
		InitialFacts: []string{
			"dataset:untrusted",
		},
		Events: []engine.Event{
			{
				ID: "worker-rce",
				Preconditions: []engine.Condition{
					{
						Fact: "dataset:untrusted",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:worker",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: "node-access",
				Preconditions: []engine.Condition{
					{
						Fact: "credential:node",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "access:node",
						Op:   engine.EffectAdd,
					},
				},
			},
			{
				ID: "service-account-token-read",
				Preconditions: []engine.Condition{
					{
						Fact: "access:worker",
						Op:   engine.ConditionPresent,
					},
				},
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	first := engine.Replay(scenario)
	second := engine.Replay(scenario)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf(
			"replay was not deterministic:\nfirst:  %#v\nsecond: %#v",
			first,
			second,
		)
	}

	wantTerminal := []string{
		"access:worker",
		"credential:k8s-service-account",
		"dataset:untrusted",
	}

	if !reflect.DeepEqual(first.TerminalState, wantTerminal) {
		t.Fatalf(
			"terminal state = %v, want %v",
			first.TerminalState,
			wantTerminal,
		)
	}

	if len(first.Trace) != 3 {
		t.Fatalf("trace length = %d, want 3", len(first.Trace))
	}

	if first.Trace[0].Status != engine.EventApplied {
		t.Errorf("worker-rce status = %q, want applied", first.Trace[0].Status)
	}

	if first.Trace[1].Status != engine.EventSkipped {
		t.Errorf("node-access status = %q, want skipped", first.Trace[1].Status)
	}

	if len(first.Trace[1].Unsatisfied) != 1 {
		t.Fatalf(
			"node-access unsatisfied conditions = %d, want 1",
			len(first.Trace[1].Unsatisfied),
		)
	}

	if first.Trace[2].Status != engine.EventApplied {
		t.Errorf(
			"service-account-token-read status = %q, want applied",
			first.Trace[2].Status,
		)
	}
}

func TestReplayRecordsControlRoles(t *testing.T) {
	scenario := engine.Scenario{
		ID: "control-test",
		Events: []engine.Event{
			{
				ID: "worker-rce",
				Effects: []engine.Effect{
					{
						Fact: "access:worker",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	controls := []engine.Control{
		engine.EventControl{
			ID:       "worker-execution-detector",
			Role:     engine.Vedette,
			EventIDs: []string{"worker-rce"},
		},
		engine.EventControl{
			ID:       "hostpath-admission",
			Role:     engine.Picket,
			EventIDs: []string{"privileged-pod"},
		},
		engine.EventControl{
			ID:       "worker-rebuild",
			Role:     engine.Reserve,
			EventIDs: []string{"worker-rce"},
		},
	}

	result := engine.Replay(scenario, controls...)

	if len(result.Trace) != 1 {
		t.Fatalf("trace length = %d, want 1", len(result.Trace))
	}

	got := result.Trace[0].Controls

	if len(got) != 3 {
		t.Fatalf("control results = %d, want 3", len(got))
	}

	if !got[0].Matched {
		t.Error("expected vedette to match worker-rce")
	}

	if got[0].Role != engine.Vedette {
		t.Errorf("role = %q, want vedette", got[0].Role)
	}

	if got[1].Matched {
		t.Error("picket unexpectedly matched worker-rce")
	}

	if !got[2].Matched {
		t.Error("expected reserve to match worker-rce")
	}

	if got[2].Role != engine.Reserve {
		t.Errorf("role = %q, want reserve", got[2].Role)
	}
}

func TestVedetteObservesWithoutBlocking(t *testing.T) {
	scenario := engine.Scenario{
		ID: "excellent-watchmen-open-gates",
		Events: []engine.Event{
			{
				ID: "credential-access",
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	vedette := engine.EventControl{
		ID:       "token-read-detector",
		Role:     engine.Vedette,
		EventIDs: []string{"credential-access"},
		Action:   engine.ActionObserve,
	}

	result := engine.Replay(scenario, vedette)

	entry := result.Trace[0]

	if entry.Status != engine.EventApplied {
		t.Fatalf("status = %q, want applied", entry.Status)
	}

	if len(entry.ObservedBy) != 1 ||
		entry.ObservedBy[0] != "token-read-detector" {
		t.Fatalf("observed by = %v, want token-read-detector", entry.ObservedBy)
	}

	if len(entry.BlockedBy) != 0 {
		t.Fatalf("blocked by = %v, want none", entry.BlockedBy)
	}

	want := []string{"credential:k8s-service-account"}

	if !reflect.DeepEqual(result.TerminalState, want) {
		t.Fatalf(
			"terminal state = %v, want %v",
			result.TerminalState,
			want,
		)
	}
}

func TestPicketBlocksEventEffects(t *testing.T) {
	scenario := engine.Scenario{
		ID: "picket-test",
		Events: []engine.Event{
			{
				ID: "credential-access",
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	picket := engine.EventControl{
		ID:       "credential-guard",
		Role:     engine.Picket,
		EventIDs: []string{"credential-access"},
		Action:   engine.ActionBlock,
	}

	result := engine.Replay(scenario, picket)

	entry := result.Trace[0]

	if entry.Status != engine.EventBlocked {
		t.Fatalf("status = %q, want blocked", entry.Status)
	}

	if len(entry.BlockedBy) != 1 ||
		entry.BlockedBy[0] != "credential-guard" {
		t.Fatalf("blocked by = %v, want credential-guard", entry.BlockedBy)
	}

	if len(result.TerminalState) != 0 {
		t.Fatalf(
			"terminal state = %v, want empty state",
			result.TerminalState,
		)
	}
}

func TestReserveRespondsAfterEvent(t *testing.T) {
	scenario := engine.Scenario{
		ID:           "reserve-test",
		InitialFacts: []string{"worker:healthy"},
		Events: []engine.Event{
			{
				ID: "worker-compromise",
				Effects: []engine.Effect{
					{
						Fact: "worker:compromised",
						Op:   engine.EffectAdd,
					},
					{
						Fact: "worker:healthy",
						Op:   engine.EffectRemove,
					},
				},
			},
		},
	}

	reserve := engine.EventControl{
		ID:       "rebuild-worker",
		Role:     engine.Reserve,
		EventIDs: []string{"worker-compromise"},
		Action:   engine.ActionRespond,
		Effects: []engine.Effect{
			{
				Fact: "worker:compromised",
				Op:   engine.EffectRemove,
			},
			{
				Fact: "worker:healthy",
				Op:   engine.EffectAdd,
			},
		},
	}

	result := engine.Replay(scenario, reserve)

	entry := result.Trace[0]

	if entry.Status != engine.EventApplied {
		t.Fatalf("status = %q, want applied", entry.Status)
	}

	if len(entry.RespondedBy) != 1 ||
		entry.RespondedBy[0] != "rebuild-worker" {
		t.Fatalf("responded by = %v, want rebuild-worker", entry.RespondedBy)
	}

	want := []string{"worker:healthy"}

	if !reflect.DeepEqual(result.TerminalState, want) {
		t.Fatalf(
			"terminal state = %v, want %v",
			result.TerminalState,
			want,
		)
	}
}

func TestObservationCanEnableResponse(t *testing.T) {
	scenario := engine.Scenario{
		ID: "observation-response-path",
		Events: []engine.Event{
			{
				ID: "credential-access",
				Effects: []engine.Effect{
					{
						Fact: "credential:k8s-service-account",
						Op:   engine.EffectAdd,
					},
				},
			},
		},
	}

	controls := []engine.Control{
		engine.EventControl{
			ID:       "credential-detector",
			Role:     engine.Vedette,
			EventIDs: []string{"credential-access"},
			Action:   engine.ActionObserve,
			Effects: []engine.Effect{
				{
					Fact: "alert:credential-access",
					Op:   engine.EffectAdd,
				},
			},
		},
		engine.EventControl{
			ID:       "soc-review",
			EventIDs: []string{"credential-access"},
			Requires: []engine.Condition{
				{
					Fact: "alert:credential-access",
					Op:   engine.ConditionPresent,
				},
			},
			Action: engine.ActionReview,
			Effects: []engine.Effect{
				{
					Fact: "reviewed:credential-access",
					Op:   engine.EffectAdd,
				},
			},
		},
		engine.EventControl{
			ID:       "criticality-check",
			EventIDs: []string{"credential-access"},
			Requires: []engine.Condition{
				{
					Fact: "reviewed:credential-access",
					Op:   engine.ConditionPresent,
				},
			},
			Action: engine.ActionCritical,
			Effects: []engine.Effect{
				{
					Fact: "critical:credential-access",
					Op:   engine.EffectAdd,
				},
			},
		},
		engine.EventControl{
			ID:       "page-oncall",
			EventIDs: []string{"credential-access"},
			Requires: []engine.Condition{
				{
					Fact: "critical:credential-access",
					Op:   engine.ConditionPresent,
				},
			},
			Action: engine.ActionEscalate,
			Effects: []engine.Effect{
				{
					Fact: "escalated:credential-access",
					Op:   engine.EffectAdd,
				},
			},
		},
		engine.EventControl{
			ID:       "credential-revocation",
			Role:     engine.Reserve,
			EventIDs: []string{"credential-access"},
			Requires: []engine.Condition{
				{
					Fact: "escalated:credential-access",
					Op:   engine.ConditionPresent,
				},
			},
			Action: engine.ActionRespond,
			Effects: []engine.Effect{
				{
					Fact: "credential:k8s-service-account",
					Op:   engine.EffectRemove,
				},
			},
		},
	}

	result := engine.Replay(scenario, controls...)

	if hasFact(result.TerminalState, "credential:k8s-service-account") {
		t.Fatalf(
			"credential remained after response: %v",
			result.TerminalState,
		)
	}

	for _, fact := range []string{
		"alert:credential-access",
		"reviewed:credential-access",
		"critical:credential-access",
		"escalated:credential-access",
	} {
		if !hasFact(result.TerminalState, fact) {
			t.Fatalf(
				"terminal state missing %q: %v",
				fact,
				result.TerminalState,
			)
		}
	}

	if len(result.Trace) != 1 {
		t.Fatalf("trace length = %d, want 1", len(result.Trace))
	}

	entry := result.Trace[0]

	review, ok := findControl(entry.Controls, "soc-review")
	if !ok {
		t.Fatal("soc-review missing from control trace")
	}

	if review.Disposition != engine.DispositionReady {
		t.Fatalf(
			"soc-review disposition = %q, want %q",
			review.Disposition,
			engine.DispositionReady,
		)
	}

	if !review.Acted {
		t.Fatal("soc-review was ready but did not act")
	}
}

func hasFact(facts []string, want string) bool {
	for _, fact := range facts {
		if fact == want {
			return true
		}
	}

	return false
}

func findControl(
	results []engine.ControlResult,
	id string,
) (engine.ControlResult, bool) {
	for _, result := range results {
		if result.ControlID == id {
			return result, true
		}
	}

	return engine.ControlResult{}, false
}

func TestReplayRefreshesSuppressedDispositionAfterEarlierPhaseAddsRequirement(t *testing.T) {
	scenario := engine.Scenario{
		ID: "phase-refresh",
		InitialFacts: []string{
			"review:unavailable",
		},
		Events: []engine.Event{
			{
				ID: "suspicious-action",
			},
		},
	}

	controls := []engine.Control{
		engine.EventControl{
			ID:       "detector",
			EventIDs: []string{"suspicious-action"},
			Action:   engine.ActionObserve,
			Effects: []engine.Effect{
				{
					Fact: "alert:suspicious-action",
					Op:   engine.EffectAdd,
				},
			},
		},
		engine.EventControl{
			ID:       "reviewer",
			EventIDs: []string{"suspicious-action"},
			Requires: []engine.Condition{
				{
					Fact: "alert:suspicious-action",
					Op:   engine.ConditionPresent,
				},
			},
			SuppressedBy: []engine.Condition{
				{
					Fact: "review:unavailable",
					Op:   engine.ConditionPresent,
				},
			},
			Action: engine.ActionReview,
		},
	}

	result := engine.Replay(scenario, controls...)

	entry := result.Trace[0]

	var reviewer engine.ControlResult
	found := false

	for _, control := range entry.Controls {
		if control.ControlID == "reviewer" {
			reviewer = control
			found = true
			break
		}
	}

	if !found {
		t.Fatal("reviewer result not found")
	}

	if reviewer.Disposition != engine.DispositionSuppressed {
		t.Fatalf(
			"reviewer disposition = %q, want suppressed",
			reviewer.Disposition,
		)
	}

	if reviewer.Action != engine.ActionNone {
		t.Fatalf(
			"suppressed reviewer action = %q, want none",
			reviewer.Action,
		)
	}

	if reviewer.Acted {
		t.Fatalf(
			"reviewer acted while suppressed: %+v",
			reviewer,
		)
	}
}

func TestInactivePicketDoesNotBlock(t *testing.T) {
	scenario := engine.Scenario{
		ID:           "inactive-picket",
		InitialFacts: []string{"seed", "picket:down"},
		Events: []engine.Event{{
			ID:            "step",
			Preconditions: []engine.Condition{{Fact: "seed", Op: engine.ConditionPresent}},
			Effects:       []engine.Effect{{Fact: "reached", Op: engine.EffectAdd}},
		}},
	}

	suppressed := engine.EventControl{
		ID: "suppressed-picket", Role: engine.Picket, EventIDs: []string{"step"}, Action: engine.ActionBlock,
		SuppressedBy: []engine.Condition{{Fact: "picket:down", Op: engine.ConditionPresent}},
	}
	waiting := engine.EventControl{
		ID: "waiting-picket", Role: engine.Picket, EventIDs: []string{"step"}, Action: engine.ActionBlock,
		Requires: []engine.Condition{{Fact: "alert:never", Op: engine.ConditionPresent}},
	}
	blind := engine.EventControl{
		ID: "waiting-vedette", Role: engine.Vedette, EventIDs: []string{"step"}, Action: engine.ActionObserve,
		Requires: []engine.Condition{{Fact: "sensor:never", Op: engine.ConditionPresent}},
	}

	for _, control := range []engine.EventControl{suppressed, waiting, blind} {
		result := engine.Replay(scenario, control)
		entry := result.Trace[0]
		if entry.Status != engine.EventApplied {
			t.Errorf("%s: status = %s, want applied", control.ID, entry.Status)
		}
		if len(entry.BlockedBy) != 0 || len(entry.ObservedBy) != 0 {
			t.Errorf("%s: inactive control credited: blocked=%v observed=%v", control.ID, entry.BlockedBy, entry.ObservedBy)
		}
	}
}
