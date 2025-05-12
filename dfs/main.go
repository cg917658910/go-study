package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Graph 表示图结构
type Graph struct {
	vertices map[int]*Vertex
	edges    []*Edge
}

// Vertex 表示图的顶点
type Vertex struct {
	key     int
	x, y    float32 // 顶点坐标
	visited bool
	parent  *Vertex
	rect    *canvas.Rectangle // GUI元素
	text    *canvas.Text      // GUI元素
}

// Edge 表示图的边
type Edge struct {
	from, to *Vertex
	line     *canvas.Line // GUI元素
	text     *canvas.Text // GUI元素
}

// NewGraph 创建新图
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[int]*Vertex),
		edges:    make([]*Edge, 0),
	}
}

// AddVertex 添加顶点
func (g *Graph) AddVertex(k int, x, y float32) *Vertex {
	if _, ok := g.vertices[k]; !ok {
		v := &Vertex{
			key:  k,
			x:    x,
			y:    y,
			rect: canvas.NewRectangle(color.White),
			text: canvas.NewText(fmt.Sprintf("%d", k), color.Black),
		}
		v.rect.SetMinSize(fyne.NewSize(30, 30))
		v.text.TextStyle.Bold = true
		v.text.Alignment = fyne.TextAlignCenter
		v.text.Move(fyne.NewPos(x+10, y+5))
		v.rect.Move(fyne.NewPos(x, y))
		g.vertices[k] = v
		return v
	}
	return g.vertices[k]
}

// AddEdge 添加边
func (g *Graph) AddEdge(from, to int) {
	fromVertex := g.vertices[from]
	toVertex := g.vertices[to]

	if fromVertex == nil || toVertex == nil {
		return
	}

	// 检查边是否已存在
	for _, e := range g.edges {
		if (e.from == fromVertex && e.to == toVertex) || (e.from == toVertex && e.to == fromVertex) {
			return
		}
	}

	line := canvas.NewLine(color.Gray{0x99})
	line.StrokeWidth = 2
	line.Position1 = fyne.NewPos(fromVertex.x+15, fromVertex.y+15)
	line.Position2 = fyne.NewPos(toVertex.x+15, toVertex.y+15)

	edge := &Edge{
		from: fromVertex,
		to:   toVertex,
		line: line,
	}
	g.edges = append(g.edges, edge)
}

// ResetVisited 重置所有顶点的访问状态
func (g *Graph) ResetVisited() {
	for _, v := range g.vertices {
		v.visited = false
		v.parent = nil
		v.rect.FillColor = color.White
		v.text.Color = color.Black
	}
	for _, e := range g.edges {
		e.line.StrokeColor = color.Gray{0x99}
	}
}

// DFSVisualization 执行DFS并可视化
func (g *Graph) DFSVisualization(startKey int, container *fyne.Container, statusLabel *widget.Label) {
	start := g.vertices[startKey]
	if start == nil {
		statusLabel.SetText("起始顶点不存在")
		return
	}

	g.ResetVisited()
	statusLabel.SetText("开始DFS搜索...")
	time.Sleep(time.Second)

	g.dfsVisual(start, container, statusLabel)

	statusLabel.SetText("DFS搜索完成")
}

// dfsVisual 带可视化的DFS实现
func (g *Graph) dfsVisual(v *Vertex, container *fyne.Container, statusLabel *widget.Label) {
	v.visited = true
	v.rect.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // 绿色表示已访问
	v.text.Color = color.White
	container.Refresh()

	statusLabel.SetText(fmt.Sprintf("访问顶点 %d", v.key))
	time.Sleep(1 * time.Second)

	for _, e := range g.edges {
		if e.from == v || e.to == v {
			var neighbor *Vertex
			if e.from == v {
				neighbor = e.to
			} else {
				neighbor = e.from
			}

			if !neighbor.visited {
				e.line.StrokeColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // 红色表示正在探索的边
				container.Refresh()
				time.Sleep(500 * time.Millisecond)

				neighbor.parent = v
				g.dfsVisual(neighbor, container, statusLabel)

				e.line.StrokeColor = color.RGBA{R: 0, G: 0, B: 255, A: 255} // 蓝色表示已探索的边
				container.Refresh()
			}
		}
	}

	v.rect.FillColor = color.RGBA{R: 0, G: 0, B: 255, A: 255} // 蓝色表示已完成
	container.Refresh()
}

// 自定义主题，加载中文字体
type CustomTheme struct {
}

// 返回自定义字体
func (c CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	fontPath := "C:/Windows/Fonts/msyh.ttc" // Windows 字体路径
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		fontPath = "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc" // Linux/macOS
	}
	font, _ := fyne.LoadResourceFromPath(fontPath)
	return font
}

func (c CustomTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (c CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (c CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("DFS")
	myWindow.Resize(fyne.NewSize(600, 500))

	graph := NewGraph()

	// 添加顶点和边（示例图）
	graph.AddVertex(0, 100, 100)
	graph.AddVertex(1, 250, 100)
	graph.AddVertex(2, 100, 250)
	graph.AddVertex(3, 250, 250)
	graph.AddVertex(4, 400, 250)

	graph.AddEdge(0, 1)
	graph.AddEdge(0, 2)
	graph.AddEdge(1, 3)
	graph.AddEdge(2, 3)
	graph.AddEdge(3, 4)

	// 创建图形容器
	graphContainer := container.NewWithoutLayout()

	// 添加边到容器
	for _, e := range graph.edges {
		graphContainer.Add(e.line)
	}

	// 添加顶点到容器
	for _, v := range graph.vertices {
		graphContainer.Add(v.rect)
		graphContainer.Add(v.text)
	}

	// 创建控制面板
	statusLabel := widget.NewLabel("read status")
	startButton := widget.NewButton("Start DFS", func() {
		go graph.DFSVisualization(0, graphContainer, statusLabel)
	})
	resetButton := widget.NewButton("重置", func() {
		graph.ResetVisited()
		statusLabel.SetText("已重置图状态")
		graphContainer.Refresh()
	})

	controls := container.NewHBox(
		startButton,
		resetButton,
		layout.NewSpacer(),
		statusLabel,
	)

	content := container.NewBorder(controls, nil, nil, nil, graphContainer)
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}
