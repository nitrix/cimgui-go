package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/nitrix/glfw-go"
	"github.com/nitrix/imgui-go"
	bglfw "github.com/nitrix/imgui-go/backends/glfw"
	bopengl3 "github.com/nitrix/imgui-go/backends/opengl3"
)

func main() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Visible, glfw.False)

	window, err := glfw.CreateWindow(1280, 720, "Example", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	glfw.SwapInterval(1)

	ctx := imgui.CreateContext(nil)
	if ctx == nil {
		panic("CreateContext failed")
	}
	defer imgui.DestroyContext(ctx)

	ctx.IO.IniFilename = nil
	ctx.IO.ConfigFlags |= imgui.ConfigFlags_ViewportsEnable
	ctx.IO.ConfigFlags |= imgui.ConfigFlags_DockingEnable
	ctx.IO.ConfigDockingWithShift = true

	if !bglfw.Init(window) {
		panic("ImGui GLFW backend init failed")
	}
	defer bglfw.Shutdown()

	if !bopengl3.Init(window) {
		panic("ImGui OpenGL3 backend init failed")
	}
	defer bopengl3.Shutdown()

	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	window.Show()

	for !window.ShouldClose() {
		glfw.PollEvents()

		gl.ClearColor(0.1, 0.1, 0.1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		bglfw.NewFrame()
		bopengl3.NewFrame()
		imgui.NewFrame()

		imgui.DockSpaceOverViewport(0, nil, imgui.DockNodeFlags_PassthruCentralNode, nil)

		imgui.ShowDemoWindow(nil)

		// Ignore this, it's just how I test some widgets.
		// testing()

		imgui.Render()

		if ctx.IO.ConfigFlags&imgui.ConfigFlags_ViewportsEnable > 0 {
			previous := glfw.GetCurrentContext()
			imgui.UpdatePlatformWindows()
			imgui.RenderPlatformWindowsDefault(nil, nil)
			previous.MakeContextCurrent()
		}

		bopengl3.RenderDrawData(imgui.GetDrawData())

		window.SwapBuffers()
	}
}
