// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build example

package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"

	. "github.com/fogleman/fauxgl"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/nfnt/resize"
)

const (
	screenWidth  = 640
	screenHeight = 480

	frameOX     = 0
	frameOY     = 0
	frameWidth  = 640
	frameHeight = 480
	frameNum    = 8

	//3D Mesh
	scale  = 1   // optional supersampling
	width  = 640 // output width in pixels
	height = 480 // output height in pixels
	fovy   = 30  // vertical field of view in degrees
	near   = 1   // near clipping plane
	far    = 10  // far clipping plane
)

var (
	runnerImage *ebiten.Image

	//3D Mesh
	mesh, err = LoadSTL("test.stl")
	context   = NewContext(width*scale, height*scale)
	eye       = V(-3, 1, -0.75)               // camera position
	center    = V(0, -0.07, 0)                // view center position
	up        = V(0, 1, 0)                    // up vector
	light     = V(-0.75, 1, 0.25).Normalize() // light direction
	color     = HexColor("#FF0000")           // object color
)

type Game struct {
	count int
}

func (g *Game) Update() error {
	g.count++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	//3D render
	context.ClearColorBufferWith(HexColor("#334FFF"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	eye = V(-3, 1, rand.Float64()*-1)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	renderedImage := context.Image()
	renderedImage = resize.Resize(width, height, renderedImage, resize.Bilinear)
	runnerImage = ebiten.NewImageFromImage(renderedImage)

	//op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	//op.GeoM.Translate(screenWidth/2, screenHeight/2)
	//i := (g.count / 5) % frameNum
	sx, sy := frameOX*frameWidth, frameOY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// load a mesh
	// mesh, err := LoadSTL("test.stl")
	// if err != nil {
	// 	panic(err)
	// }

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(Radians(30))

	// create a rendering context
	//context := NewContext(width*scale, height*scale)
	context.ClearColorBufferWith(HexColor("#334FFF"))

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := LookAt(eye, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := NewPhongShader(matrix, light, eye)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	renderedImage := context.Image()
	renderedImage = resize.Resize(width, height, renderedImage, resize.Bilinear)

	// save image
	//SavePNG("out.png", image)
	// Decode image from a byte slice instead of a file so that
	// this example works in any working directory.
	// If you want to use a file, there are some options:
	// 1) Use os.Open and pass the file to the image decoder.
	//    This is a very regular way, but doesn't work on browsers.
	// 2) Use ebitenutil.OpenFile and pass the file to the image decoder.
	//    This works even on browsers.
	// 3) Use ebitenutil.NewImageFromFile to create an ebiten.Image directly from a file.
	//    This also works on browsers.
	//img, _, err := image.Decode(renderedImage)
	img := renderedImage
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Game")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
