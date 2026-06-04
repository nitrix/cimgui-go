package main

import "testing"

func TestCToGoTypeUsesSafePointerShapes(t *testing.T) {
	tests := map[string]string{
		"int*":          "*int32",
		"unsigned int*": "*uint32",
		"int[2]":        "*[2]int32",
		"float[2]":      "*[2]float32",
		"char* const[]": "[]string",
	}

	for input, want := range tests {
		if got := cToGoType(input); got != want {
			t.Fatalf("cToGoType(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestIRRetainsDefaultsWithoutGeneratingWrapperNames(t *testing.T) {
	def := definitionT{
		OverloadedName: "igSetNextWindowPos",
		Ret:            "void",
		Defaults: map[string]string{
			"cond":  "0",
			"pivot": "ImVec2(0,0)",
		},
		ArgsT: []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			{Name: "pos", Type: "const ImVec2"},
			{Name: "cond", Type: "ImGuiCond"},
			{Name: "pivot", Type: "const ImVec2"},
		},
	}

	ir := buildFunctionIR(def, false)
	if ir.DefaultWrapper {
		t.Fatal("expected defaults to stay metadata-only")
	}
	if ir.FullGoName != "SetNextWindowPos" {
		t.Fatalf("FullGoName = %q, want SetNextWindowPos", ir.FullGoName)
	}
	if ir.Params[1].Default != "0" {
		t.Fatalf("cond default = %q, want 0", ir.Params[1].Default)
	}
}

func TestShouldGenerateCoreFunctionReportsSkippedCallbacks(t *testing.T) {
	def := definitionT{
		OverloadedName: "igPlotLines_FnFloatPtr",
		ArgsT: []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}{
			{Name: "values_getter", Type: "float(*)(void* data,int idx)"},
		},
	}

	ok, reason := shouldGenerateCoreFunction(def)
	if ok {
		t.Fatal("expected function-pointer overload to be skipped")
	}
	if reason != "function pointer" {
		t.Fatalf("skip reason = %q, want function pointer", reason)
	}
}
