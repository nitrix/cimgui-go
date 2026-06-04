package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
)

type EnumT struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	CalcValue int    `json:"calc_value"`
}

type FieldT struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type structsAndEnumsT struct {
	Enums   map[string][]EnumT  `json:"enums"`
	Structs map[string][]FieldT `json:"structs"`
}

type typedefsDictT map[string]string

type definitionsT map[string][]definitionT

type definitionT struct {
	Args  string `json:"args"`
	ArgsT []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"argsT"`
	Name           string            `json:"cimguiname"`
	Defaults       map[string]string `json:"defaults"`
	Namespace      string            `json:"namespace"`
	OverloadedName string            `json:"ov_cimguiname"`
	Ret            string            `json:"ret"`
	Signature      string            `json:"signature"`
	StructName     string            `json:"stname"`
	Constructor    bool              `json:"constructor"`
	Destructor     bool              `json:"destructor"`
	Location       string            `json:"location"`
}

type paramIR struct {
	Name          string `json:"name"`
	CType         string `json:"cType"`
	GoType        string `json:"goType,omitempty"`
	CgoType       string `json:"cgoType,omitempty"`
	Default       string `json:"default,omitempty"`
	BufferSizeFor string `json:"bufferSizeFor,omitempty"`
	StringArray   bool   `json:"stringArray,omitempty"`
	Optional      bool   `json:"optional,omitempty"`
}

type functionIR struct {
	CName           string    `json:"cName"`
	OverloadedCName string    `json:"overloadedCName"`
	GoName          string    `json:"goName,omitempty"`
	FullGoName      string    `json:"fullGoName,omitempty"`
	Namespace       string    `json:"namespace,omitempty"`
	StructName      string    `json:"structName,omitempty"`
	Location        string    `json:"location,omitempty"`
	ReturnType      string    `json:"returnType,omitempty"`
	GoReturnType    string    `json:"goReturnType,omitempty"`
	Params          []paramIR `json:"params,omitempty"`
	Variadic        bool      `json:"variadic,omitempty"`
	DefaultWrapper  bool      `json:"defaultWrapper,omitempty"`
	Constructor     bool      `json:"constructor,omitempty"`
	Destructor      bool      `json:"destructor,omitempty"`
	SkipReason      string    `json:"skipReason,omitempty"`
}

type generatorReport struct {
	GeneratedCoreFunctions int          `json:"generatedCoreFunctions"`
	DefaultWrappers        int          `json:"defaultWrappers"`
	VariadicWrappers       int          `json:"variadicWrappers"`
	SkippedCoreFunctions   []functionIR `json:"skippedCoreFunctions,omitempty"`
	BackendFunctions       []functionIR `json:"backendFunctions,omitempty"`
}

var cEnumCastPattern = regexp.MustCompile(`\([A-Za-z_][A-Za-z0-9_]*\)`)

func isUnsupportedFeatureSymbol(name string) bool {
	return strings.Contains(name, "FreeType")
}

func main() {
	content, err := os.ReadFile("thirdparty/cimgui/generator/output/structs_and_enums.json")
	if err != nil {
		panic(err)
	}

	// Parse the JSON file.
	var structsAndEnums structsAndEnumsT
	err = json.NewDecoder(bytes.NewReader(content)).Decode(&structsAndEnums)
	if err != nil {
		panic(err)
	}

	content, err = os.ReadFile("thirdparty/cimgui/generator/output/definitions.json")
	if err != nil {
		panic(err)
	}

	// Parse the JSON file.
	var definitions definitionsT
	err = json.NewDecoder(bytes.NewReader(content)).Decode(&definitions)
	if err != nil {
		panic(err)
	}

	content, err = os.ReadFile("thirdparty/cimgui/generator/output/impl_definitions.json")
	if err != nil {
		panic(err)
	}

	var implDefinitions definitionsT
	err = json.NewDecoder(bytes.NewReader(content)).Decode(&implDefinitions)
	if err != nil {
		panic(err)
	}

	content, err = os.ReadFile("thirdparty/cimgui/generator/output/typedefs_dict.json")
	if err != nil {
		panic(err)
	}

	// Parse the JSON file.
	var typedefsDict map[string]string
	err = json.NewDecoder(bytes.NewReader(content)).Decode(&typedefsDict)
	if err != nil {
		panic(err)
	}

	// POC.
	copyFile("thirdparty/cimgui/cimgui.cpp", "dist/cimgui/cimgui.cpp")
	copyFile("thirdparty/cimgui/cimgui.h", "dist/cimgui/cimgui.h")
	copyFile("thirdparty/cimgui/cimconfig.h", "dist/cimgui/cimconfig.h")
	copyFile("thirdparty/cimgui/cimgui_impl.h", "dist/cimgui/cimgui_impl.h")

	copyFile("thirdparty/cimgui/imgui/imgui.h", "dist/imgui/imgui.h")
	copyFile("thirdparty/cimgui/imgui/imconfig.h", "dist/imgui/imconfig.h")
	copyFile("thirdparty/cimgui/imgui/imgui.cpp", "dist/imgui/imgui.cpp")
	copyFile("thirdparty/cimgui/imgui/imgui_internal.h", "dist/imgui/imgui_internal.h")
	copyFile("thirdparty/cimgui/imgui/imgui_demo.cpp", "dist/imgui/imgui_demo.cpp")
	copyFile("thirdparty/cimgui/imgui/imgui_draw.cpp", "dist/imgui/imgui_draw.cpp")
	copyFile("thirdparty/cimgui/imgui/imgui_tables.cpp", "dist/imgui/imgui_tables.cpp")
	copyFile("thirdparty/cimgui/imgui/imgui_widgets.cpp", "dist/imgui/imgui_widgets.cpp")
	copyFile("thirdparty/cimgui/imgui/imstb_rectpack.h", "dist/imgui/imstb_rectpack.h")
	copyFile("thirdparty/cimgui/imgui/imstb_textedit.h", "dist/imgui/imstb_textedit.h")
	copyFile("thirdparty/cimgui/imgui/imstb_truetype.h", "dist/imgui/imstb_truetype.h")

	copyFile("thirdparty/cimgui/imgui/backends/imgui_impl_glfw.cpp", "backends/glfw/imgui_impl_glfw.cpp")
	copyFile("thirdparty/cimgui/imgui/backends/imgui_impl_glfw.h", "backends/glfw/imgui_impl_glfw.h")

	copyFile("thirdparty/cimgui/imgui/backends/imgui_impl_opengl3.cpp", "backends/opengl3/imgui_impl_opengl3.cpp")
	copyFile("thirdparty/cimgui/imgui/backends/imgui_impl_opengl3.h", "backends/opengl3/imgui_impl_opengl3.h")
	copyFile("thirdparty/cimgui/imgui/backends/imgui_impl_opengl3_loader.h", "backends/opengl3/imgui_impl_opengl3_loader.h")

	copyFile("thirdparty/glfw/include/GLFW/glfw3.h", "backends/glfw/GLFW/glfw3.h")
	copyFile("thirdparty/glfw/include/GLFW/glfw3native.h", "backends/glfw/GLFW/glfw3native.h")

	generateConstants(&structsAndEnums)
	generateTypedefs(typedefsDict)
	generateDefinitions(definitions)
	generateWrappersHeaders(definitions)
	generateWrappersSources(definitions)
	generateReport(definitions, implDefinitions)
}

func copyFile(src, dst string) {
	_ = os.MkdirAll(filepath.Dir(dst), 0750)
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Copied file %s to %s\n", src, dst)
}

func writeGoFile(path string, content string) {
	formatted, err := format.Source([]byte(content))
	if err != nil {
		panic(fmt.Errorf("format %s: %w", path, err))
	}

	err = os.WriteFile(path, formatted, 0644)
	if err != nil {
		panic(err)
	}
}

func generateConstants(structsAndEnums *structsAndEnumsT) {
	constantsContent := strings.Builder{}
	constantsContent.WriteString("//go:build cgo\n")
	constantsContent.WriteString("\n")
	constantsContent.WriteString("package imgui\n")
	constantsContent.WriteString("\n")

	flattenedEnums := make([]EnumT, 0, len(structsAndEnums.Enums))
	for _, enum := range structsAndEnums.Enums {
		flattenedEnums = append(flattenedEnums, enum...)
	}

	sort.Slice(flattenedEnums, func(i, j int) bool {
		return flattenedEnums[i].Name < flattenedEnums[j].Name
	})

	for _, field := range flattenedEnums {
		if isUnsupportedFeatureSymbol(field.Name) {
			continue
		}

		renamedName := strings.ReplaceAll(field.Name, "ImGui", "")
		renamedValue := strings.ReplaceAll(field.Value, "ImGui", "")
		renamedValue = strings.ReplaceAll(renamedValue, "~", "^")
		renamedValue = cEnumCastPattern.ReplaceAllString(renamedValue, "")

		constantsContent.WriteString(fmt.Sprintf("const %s = %s\n", renamedName, renamedValue))
	}

	writeGoFile("imgui_constants.go", constantsContent.String())
}

func isImLike(s string) bool {
	if strings.HasPrefix(s, "Im") {
		r, _ := utf8.DecodeRuneInString(s[2:])
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func generateTypedefs(typedefsDict typedefsDictT) {
	typedefsContent := strings.Builder{}
	typedefsContent.WriteString("//go:build cgo\n")
	typedefsContent.WriteString("\n")
	typedefsContent.WriteString("package imgui\n")
	typedefsContent.WriteString("\n")
	typedefsContent.WriteString("// #define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1\n")
	typedefsContent.WriteString("// #include \"dist/cimgui/cimgui.h\"\n")
	typedefsContent.WriteString("import \"C\"\n")
	typedefsContent.WriteString("\n")

	// sortedNames := make([]string, 0, len(structsAndEnums.Structs))
	// for name := range structsAndEnums.Structs {
	// 	sortedNames = append(sortedNames, name)
	// }
	// sort.Strings(sortedNames)

	sortedNames := make([]string, 0, len(typedefsDict))
	for name := range typedefsDict {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	blacklist := []string{
		"const_iterator", "iterator", "value_type", "STB_TexteditState",
		"stbrp_context_opaque", "stbrp_node", "stbrp_node_im",
		"ImVec1", "ImVec2", "ImVec2ih", "ImVec4", // These are specially handled and replaced by mgl32.Vec2 and mgl32.Vec4.
	}

	for _, name := range sortedNames {
		if slices.Contains(blacklist, name) {
			continue
		}

		if isUnsupportedFeatureSymbol(name) {
			continue
		}

		// s := structsAndEnums.Structs[name]
		// _ = s

		// Use `go tool cgo -godefs` to generate the Go struct from the C struct?

		// reflect.TypeOf()
		// typedefsContent.WriteString(fmt.Sprintf("type %s struct {\n", name))
		// for _, field := range s {
		// typedefsContent.WriteString(fmt.Sprintf("\t%s %s\n", safeIdentifier(field.Name), cToGoType(field.Type)))
		// }
		// typedefsContent.WriteString("}\n\n")

		typedefsContent.WriteString(fmt.Sprintf("type %s C.%s\n", cToGoType(name), name))
	}

	writeGoFile("imgui_typedefs.go", typedefsContent.String())
}

func cToCgoType(cType string) string {
	originalCType := cType

	if strings.HasPrefix(cType, "const ") {
		cType = strings.TrimSpace(strings.TrimPrefix(cType, "const "))
		return cToCgoType(cType)
	}

	switch cType {
	case "char":
		return "C.char"
	case "bool":
		return "C.bool"
	case "int":
		return "C.int"
	case "float":
		return "C.float"
	case "float*":
		return "*C.float"
	case "double":
		return "C.double"
	case "double*":
		return "*C.double"
	case "int*":
		return "*C.int"
	case "unsigned int*":
		return "*C.uint"
	case "unsigned int":
		return "C.uint"
	case "char*":
		return "*C.char"
	case "bool*":
		return "*C.bool"
	case "void*":
		return "unsafe.Pointer"
	case "size_t":
		return "C.size_t"
	case "size_t*":
		return "*C.size_t"
	case "int[2]":
		return "*C.int"
	case "int[3]":
		return "*C.int"
	case "int[4]":
		return "*C.int"
	case "float[2]":
		return "*C.float"
	case "float[3]":
		return "*C.float"
	case "float[4]":
		return "*C.float"
	case "char[5]":
		return "*C.char"
	case "char* const[]":
		return "**C.char"
	case "unsigned char":
		return "C.uchar"
	case "unsigned char*":
		return "*C.uchar"
	case "unsigned char[256]":
		return "*C.uchar"
	}

	if isImLike(cType) || strings.HasPrefix(cType, "ig") {
		if strings.HasSuffix(cType, "**") {
			cType = "**C." + strings.TrimSuffix(cType, "**")
		} else if strings.HasSuffix(cType, "*") {
			cType = "*C." + strings.TrimSuffix(cType, "*")
		} else {
			cType = "C." + cType
		}
		return cType
	}

	panic(fmt.Errorf("unknown C type to convert to Cgo: %s", originalCType))
}

func cToGoType(cType string) string {
	originalCType := cType

	if strings.HasPrefix(cType, "const ") {
		cType = strings.TrimSpace(strings.TrimPrefix(cType, "const "))
		return cToGoType(cType)
	}

	switch cType {
	case "int":
		return "int"
	case "int*":
		return "*int32"
	case "int[2]":
		return "*[2]int32"
	case "int[3]":
		return "*[3]int32"
	case "int[4]":
		return "*[4]int32"
	case "float":
		return "float32"
	case "float*":
		return "*float32"
	case "float[2]":
		return "*[2]float32"
	case "float[3]":
		return "*mgl32.Vec3"
	case "float[4]":
		return "*mgl32.Vec4"
	case "double":
		return "float64"
	case "double*":
		return "*float64"
	case "char":
		return "byte"
	case "bool":
		return "bool"
	case "bool*":
		return "*bool"
	case "char*":
		return "string"
	case "unsigned int*":
		return "*uint32"
	case "unsigned int":
		return "uint"
	case "void*":
		return "unsafe.Pointer"
	case "size_t":
		return "uint"
	case "size_t*":
		return "*uint"
	case "unsigned char[256]":
		return "[256]byte"
	case "unsigned char*":
		return "[]byte"
	case "unsigned char":
		return "byte"
	case "char* const[]":
		return "[]string"
	case "char[5]":
		return "[5]byte"

		// case "short":
		// 	return "int16"
		// case "unsigned short":
		// 	return "uint16"
		// case "bool(*)(ImFontAtlas* atlas)":
		// 	return "func(atlas *ImFontAtlas) bool"
		// case "void(*)(ImGuiContext* ctx,ImGuiDockNode* node,ImGuiTabBar* tab_bar)":
		// 	return "func(ctx *ImGuiContext, node *ImGuiDockNode, tab_bar *ImGuiTabBar)"
	}

	switch cType {
	case "ImRect_c*":
		return "*Rect"
	case "ImRect_c":
		return "Rect"
	case "ImVec2_c*":
		return "*mgl32.Vec2"
	case "ImVec2_c":
		return "mgl32.Vec2"
	case "ImVec4_c":
		return "mgl32.Vec4"
	case "ImVec4_c*":
		return "*mgl32.Vec4"
	case "ImVec2*":
		return "*mgl32.Vec2"
	case "ImVec2":
		return "mgl32.Vec2"
	case "ImVec4":
		return "mgl32.Vec4"
	case "ImVec4*":
		return "*mgl32.Vec4"
	}

	if strings.HasPrefix(cType, "ImGui") {
		cType = strings.TrimPrefix(cType, "ImGui")
		if strings.HasSuffix(cType, "**") {
			cType = "**" + strings.TrimSuffix(cType, "**")
		}
		if strings.HasSuffix(cType, "*") {
			cType = "*" + strings.TrimSuffix(cType, "*")
		}
		return cType
	}

	if isImLike(cType) {
		cType = cType[2:]
		if strings.HasSuffix(cType, "*") {
			cType = "*" + strings.TrimSuffix(cType, "*")
		}
		return cType
	}

	panic(fmt.Errorf("unknown C type to convert go Go: %s", originalCType))
}

func safeIdentifier(s string) string {
	s = strcase.ToLowerCamel(s)

	switch s {
	case "type":
		return "ty"
	case "func":
		return "fn"
	case "map":
		return "m"
	case "range":
		return "rangeArg"
	case "string":
		return "str"
	}

	return s
}

func hasVariadic(def definitionT) bool {
	for _, arg := range def.ArgsT {
		if arg.Name == "..." {
			return true
		}
	}
	return false
}

func shouldGenerateCoreFunction(def definitionT) (bool, string) {
	blacklist := []string{
		"igNewFrame",                 // Needed to be overridden to control the pool memory allocator.
		"igGetAllocatorFunctions",    // Won't be tweaking the allocator functions from Go.
		"igAddDrawListToDrawDataEx",  // TODO: Wants an ImVector_ImDrawListPtr*
		"igDockBuilderCopyDockSpace", // TODO: Wants a ImVector_const_charPtr*
		"igDockBuilderCopyNode",      // TODO: Wants a ImVector_ImGuiID*
	}

	if slices.Contains(blacklist, def.OverloadedName) {
		return false, "blacklisted"
	}

	if !strings.HasPrefix(def.OverloadedName, "ig") {
		return false, "non-imgui namespace"
	}

	if strings.HasPrefix(def.OverloadedName, "igIm") {
		return false, "internal helper"
	}

	renamedName := strings.TrimPrefix(def.OverloadedName, "ig")
	if isImLike(renamedName) {
		renamedName = renamedName[2:]
	}

	if strings.HasPrefix(renamedName, "Debug") {
		return false, "debug helper"
	}

	for _, arg := range def.ArgsT {
		if arg.Type == "va_list" {
			return false, "va_list"
		}
	}

	for _, arg := range def.ArgsT {
		if strings.Contains(arg.Type, "(*)") {
			return false, "function pointer"
		}
	}

	return true, ""
}

func trailingDefaultStart(def definitionT) int {
	if len(def.Defaults) == 0 || hasVariadic(def) {
		return len(def.ArgsT)
	}

	start := len(def.ArgsT)
	for start > 0 {
		arg := def.ArgsT[start-1]
		if _, ok := def.Defaults[arg.Name]; !ok {
			break
		}
		start--
	}

	if start == len(def.ArgsT) {
		return len(def.ArgsT)
	}

	for i := start; i < len(def.ArgsT); i++ {
		if _, ok := def.Defaults[def.ArgsT[i].Name]; !ok {
			return len(def.ArgsT)
		}
	}

	return start
}

func goDefaultExpr(def definitionT, argIndex int) (string, bool) {
	arg := def.ArgsT[argIndex]
	value, ok := def.Defaults[arg.Name]
	if !ok {
		return "", false
	}

	value = strings.TrimSpace(value)
	goType := goArgumentType(def, argIndex)

	switch value {
	case "NULL", "nullptr", "((void*)0)":
		if goType == "string" {
			return "\"\"", true
		}
		return "nil", strings.HasPrefix(goType, "*") || goType == "unsafe.Pointer"
	case "false":
		return "false", true
	case "true":
		return "true", true
	case "0", "0.0f", "0.0":
		if goType == "mgl32.Vec2" {
			return "mgl32.Vec2{}", true
		}
		if goType == "mgl32.Vec4" {
			return "mgl32.Vec4{}", true
		}
		return "0", true
	}

	if value == "ImVec2(0,0)" || value == "ImVec2(0.0f,0.0f)" {
		return "mgl32.Vec2{}", goType == "mgl32.Vec2"
	}
	if value == "ImVec4(0,0,0,0)" || value == "ImVec4(0.0f,0.0f,0.0f,0.0f)" {
		return "mgl32.Vec4{}", goType == "mgl32.Vec4"
	}

	if strings.HasPrefix(value, "ImGui") {
		return strings.TrimPrefix(value, "ImGui"), true
	}

	return "", false
}

func canGenerateDefaultWrapper(def definitionT) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	start := trailingDefaultStart(def)
	if start == len(def.ArgsT) {
		return false
	}

	for i := start; i < len(def.ArgsT); i++ {
		if isCharBufferSizeArgument(def, i) {
			continue
		}
		if _, ok := goDefaultExpr(def, i); !ok {
			return false
		}
	}

	return true
}

func goNameForDefinition(def definitionT) string {
	renamedName := strings.TrimPrefix(def.OverloadedName, "ig")
	if isImLike(renamedName) {
		renamedName = renamedName[2:]
	}
	return renamedName
}

func safeCToGoType(cType string) (goType string, ok bool) {
	defer func() {
		if recover() != nil {
			goType = ""
			ok = false
		}
	}()
	return cToGoType(cType), true
}

func safeCToCgoType(cType string) (cgoType string, ok bool) {
	defer func() {
		if recover() != nil {
			cgoType = ""
			ok = false
		}
	}()
	return cToCgoType(cType), true
}

func buildFunctionIR(def definitionT, backend bool) functionIR {
	ir := functionIR{
		CName:           def.Name,
		OverloadedCName: def.OverloadedName,
		GoName:          goNameForDefinition(def),
		Namespace:       def.Namespace,
		StructName:      def.StructName,
		Location:        def.Location,
		ReturnType:      def.Ret,
		Variadic:        hasVariadic(def),
		Constructor:     def.Constructor,
		Destructor:      def.Destructor,
	}

	ir.FullGoName = ir.GoName

	if def.Ret != "" && def.Ret != "void" {
		if goType, ok := safeCToGoType(def.Ret); ok {
			ir.GoReturnType = goType
		}
	}

	for i, arg := range def.ArgsT {
		param := paramIR{
			Name:  arg.Name,
			CType: arg.Type,
		}

		if value, ok := def.Defaults[arg.Name]; ok {
			param.Default = value
			param.Optional = value == "NULL" || value == "nullptr" || value == "((void*)0)"
		}

		if isCharBufferSizeArgument(def, i) {
			param.BufferSizeFor = safeIdentifier(def.ArgsT[i-1].Name)
		} else if goType, ok := safeCToGoType(arg.Type); ok {
			param.GoType = goType
			param.StringArray = goType == "[]string"
		}

		if cgoType, ok := safeCToCgoType(arg.Type); ok {
			param.CgoType = cgoType
		}

		ir.Params = append(ir.Params, param)
	}

	if backend {
		ir.GoName = ""
		ir.FullGoName = ""
	}

	return ir
}

func flattenDefinitions(definitions definitionsT) []definitionT {
	flattened := make([]definitionT, 0, len(definitions))
	for _, definition := range definitions {
		flattened = append(flattened, definition...)
	}
	sort.Slice(flattened, func(i, j int) bool {
		return flattened[i].OverloadedName < flattened[j].OverloadedName
	})
	return flattened
}

func generateReport(definitions definitionsT, implDefinitions definitionsT) {
	report := generatorReport{}

	for _, def := range flattenDefinitions(definitions) {
		shouldGenerate, reason := shouldGenerateCoreFunction(def)
		ir := buildFunctionIR(def, false)
		if !shouldGenerate {
			ir.SkipReason = reason
			report.SkippedCoreFunctions = append(report.SkippedCoreFunctions, ir)
			continue
		}

		report.GeneratedCoreFunctions++
		if ir.Variadic {
			report.VariadicWrappers++
		}
	}

	for _, def := range flattenDefinitions(implDefinitions) {
		report.BackendFunctions = append(report.BackendFunctions, buildFunctionIR(def, true))
	}

	content, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("generator_report.json", append(content, '\n'), 0644)
	if err != nil {
		panic(err)
	}
}

func goArgumentType(def definitionT, argIndex int) string {
	arg := def.ArgsT[argIndex]
	if arg.Type == "char*" {
		return "[]byte"
	}
	return cToGoType(arg.Type)
}

func isCharBufferSizeArgument(def definitionT, argIndex int) bool {
	if argIndex == 0 {
		return false
	}

	arg := def.ArgsT[argIndex]
	if arg.Type != "int" && arg.Type != "size_t" {
		return false
	}

	previousArg := def.ArgsT[argIndex-1]
	if previousArg.Type != "char*" {
		return false
	}

	return strings.Contains(strings.ToLower(arg.Name), "size")
}

func generateDefinitions(definitions definitionsT) {
	output := strings.Builder{}
	output.WriteString("//go:build cgo\n")
	output.WriteString("\n")
	output.WriteString("package imgui\n")
	output.WriteString("\n")
	output.WriteString("// #define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1\n")
	output.WriteString("// #include \"dist/cimgui/cimgui.h\"\n")
	output.WriteString("// #include \"imgui_wrappers.h\"\n")
	output.WriteString("// #include <stdbool.h>\n")
	output.WriteString("import \"C\"\n")
	output.WriteString("import \"fmt\"\n")
	output.WriteString("import \"unsafe\"\n")
	output.WriteString("import \"github.com/go-gl/mathgl/mgl32\"\n")
	output.WriteString("\n")

	sortedDefinitions := make([]definitionT, 0, len(definitions))
	for _, definition := range definitions {
		sortedDefinitions = append(sortedDefinitions, definition...)
	}
	sort.Slice(sortedDefinitions, func(i, j int) bool {
		return sortedDefinitions[i].OverloadedName < sortedDefinitions[j].OverloadedName
	})

	for _, def := range sortedDefinitions {
		shouldGenerate, _ := shouldGenerateCoreFunction(def)
		if !shouldGenerate {
			continue
		}

		// fmt.Println("=>", def.OverloadedName)

		renamedName := goNameForDefinition(def)

		output.WriteString(fmt.Sprintf("func %s(", renamedName))

		// Parameters.
		defHasVariadic := false
		{
			params := []string{}
			for i, arg := range def.ArgsT {
				isNextArgumentVariadic := len(def.ArgsT) > i+1 && def.ArgsT[i+1].Name == "..."

				if isNextArgumentVariadic {
					params = append(params, "vfmt string, vargs ...interface{}")
					defHasVariadic = true
					break
				} else if isCharBufferSizeArgument(def, i) {
					continue
				} else {
					params = append(params, fmt.Sprintf("%s %s", safeIdentifier(arg.Name), goArgumentType(def, i)))
				}
			}
			output.WriteString(strings.Join(params, ", "))
		}

		output.WriteString(") ")

		hasReturnType := def.Ret != "void" && def.Ret != ""

		// Return type.
		{
			if hasReturnType {
				output.WriteString(cToGoType(def.Ret))
				output.WriteString(" ")
			}
		}

		output.WriteString("{\n")

		// Body.
		{
			for i, arg := range def.ArgsT {
				isNextArgumentVariadic := len(def.ArgsT) > i+1 && def.ArgsT[i+1].Name == "..."

				output.WriteString(fmt.Sprintf("\ta%d := ", i))

				expr := safeIdentifier(arg.Name)

				cType := arg.Type
				goType := goArgumentType(def, i)
				cgoType := cToCgoType(cType)

				if isNextArgumentVariadic {
					expr = "fmt.Sprintf(vfmt, vargs...)"
				}

				if cType == "char*" {
					output.WriteString(fmt.Sprintf("(*C.char)(nil)\n\tif len(%s) > 0 {\n\t\ta%d = (*C.char)(unsafe.Pointer(&%s[0]))\n\t}", expr, i, expr))
				} else if isCharBufferSizeArgument(def, i) {
					previousArgName := safeIdentifier(def.ArgsT[i-1].Name)
					output.WriteString(fmt.Sprintf("(%s)(len(%s))", cgoType, previousArgName))
				} else if goType == "[]string" {
					output.WriteString(fmt.Sprintf("stringPool.StoreCStringArray(%s)", expr))
				} else {
					switch goType {
					case "string":
						expr = fmt.Sprintf("stringPool.StoreCString(%s)", expr)
					case "mgl32.Vec2":
						expr = fmt.Sprintf("mglVec2ToImVec2(%s)", expr)
					case "mgl32.Vec4":
						expr = fmt.Sprintf("mglVec4ToImVec4(%s)", expr)
					case "*mgl32.Vec2":
						output.WriteString(fmt.Sprintf("(%s)(nil)\n\tif %s != nil {\n\t\ta%d = (%s)(unsafe.Pointer(&%s[0]))\n\t}", cgoType, expr, i, cgoType, expr))
						output.WriteString("\n")
						continue
					case "*mgl32.Vec4":
						output.WriteString(fmt.Sprintf("(%s)(nil)\n\tif %s != nil {\n\t\ta%d = (%s)(unsafe.Pointer(&%s[0]))\n\t}", cgoType, expr, i, cgoType, expr))
						output.WriteString("\n")
						continue
					default:
						if strings.HasPrefix(goType, "[") {
							expr = fmt.Sprintf("&%s[0]", expr)
						}
						if strings.HasPrefix(cgoType, "*") {
							expr = fmt.Sprintf("unsafe.Pointer(%s)", expr)
						}
						expr = fmt.Sprintf("(%s)(%s)", cgoType, expr)
					}

					output.WriteString(expr)
				}

				output.WriteString("\n")

				if isNextArgumentVariadic {
					break
				}
			}

			output.WriteString("\t")
			if hasReturnType {
				output.WriteString("call := ")
			}

			if defHasVariadic {
				output.WriteString(fmt.Sprintf("C.wrap_%s(", def.OverloadedName))
			} else {
				output.WriteString(fmt.Sprintf("C.%s(", def.OverloadedName))
			}

			for i := range def.ArgsT {
				if i > 0 {
					output.WriteString(", ")
				}

				output.WriteString(fmt.Sprintf("a%d", i))

				isNextArgumentVariadic := len(def.ArgsT) > i+1 && def.ArgsT[i+1].Name == "..."
				if isNextArgumentVariadic {
					break
				}
			}

			output.WriteString(")\n")

			if hasReturnType {
				expr := "call"

				goType := cToGoType(def.Ret)

				if goType == "string" {
					expr = fmt.Sprintf("C.GoString(%s)", expr)
				} else if goType == "mgl32.Vec2" {
					expr = fmt.Sprintf("imVec2ToMglVec2(%s)", expr)
				} else if goType == "mgl32.Vec4" {
					expr = fmt.Sprintf("imVec4ToMglVec4(%s)", expr)
				} else if goType == "*mgl32.Vec2" {
					expr = fmt.Sprintf("(*mgl32.Vec2)(unsafe.Pointer(%s))", expr)
				} else if goType == "*mgl32.Vec4" {
					expr = fmt.Sprintf("(*mgl32.Vec4)(unsafe.Pointer(%s))", expr)
				} else {
					expr = fmt.Sprintf("(%s)(%s)", goType, expr)
				}

				output.WriteString(fmt.Sprintf("\treturn %s\n", expr))
			}
		}

		output.WriteString("}\n\n")

	}

	writeGoFile("imgui_functions.go", output.String())
}

func wrapperPrototypeForDefinition(def definitionT) string {
	prototype := strings.Builder{}
	prototype.WriteString(fmt.Sprintf("%s wrap_%s(", def.Ret, def.OverloadedName))

	for i, arg := range def.ArgsT {
		if i > 0 {
			prototype.WriteString(", ")
		}

		nextArgIsVariadic := len(def.ArgsT) > i+1 && def.ArgsT[i+1].Name == "..."

		prototype.WriteString(fmt.Sprintf("%s %s", arg.Type, arg.Name))

		if nextArgIsVariadic {
			break
		}
	}

	prototype.WriteString(")")
	return prototype.String()
}

func generateWrappersSources(definitions definitionsT) {
	output := strings.Builder{}
	output.WriteString("#define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1\n")
	output.WriteString("#include \"dist/cimgui/cimgui.h\"\n")
	output.WriteString("\n")

	sortedDefinitions := make([]definitionT, 0, len(definitions))
	for _, definition := range definitions {
		sortedDefinitions = append(sortedDefinitions, definition...)
	}
	sort.Slice(sortedDefinitions, func(i, j int) bool {
		return sortedDefinitions[i].OverloadedName < sortedDefinitions[j].OverloadedName
	})

	for _, def := range sortedDefinitions {
		shouldGenerate, _ := shouldGenerateCoreFunction(def)
		if !shouldGenerate || !hasVariadic(def) {
			continue
		}

		output.WriteString(wrapperPrototypeForDefinition(def))
		output.WriteString(" {\n")
		output.WriteString("\t")

		if def.Ret != "void" && def.Ret != "" {
			output.WriteString("return ")
		}

		output.WriteString(fmt.Sprintf("%s(", def.OverloadedName))
		for i, arg := range def.ArgsT {
			if i > 0 {
				output.WriteString(", ")
			}

			nextArgIsVariadic := len(def.ArgsT) > i+1 && def.ArgsT[i+1].Name == "..."

			if nextArgIsVariadic {
				output.WriteString(arg.Name)
				break
			}

			output.WriteString(arg.Name)
		}

		output.WriteString(");\n")
		output.WriteString("}\n\n")
	}

	err := os.WriteFile("imgui_wrappers.c", []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func generateWrappersHeaders(definitions definitionsT) {
	output := strings.Builder{}
	output.WriteString("#define CIMGUI_DEFINE_ENUMS_AND_STRUCTS 1\n")
	output.WriteString("#include \"dist/cimgui/cimgui.h\"\n")
	output.WriteString("\n")

	sortedDefinitions := make([]definitionT, 0, len(definitions))
	for _, definition := range definitions {
		sortedDefinitions = append(sortedDefinitions, definition...)
	}
	sort.Slice(sortedDefinitions, func(i, j int) bool {
		return sortedDefinitions[i].OverloadedName < sortedDefinitions[j].OverloadedName
	})

	for _, def := range sortedDefinitions {
		shouldGenerate, _ := shouldGenerateCoreFunction(def)
		if !shouldGenerate || !hasVariadic(def) {
			continue
		}

		output.WriteString(wrapperPrototypeForDefinition(def))
		output.WriteString(";\n")
	}

	err := os.WriteFile("imgui_wrappers.h", []byte(output.String()), 0644)
	if err != nil {
		panic(err)
	}
}
