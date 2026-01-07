package main

// Proxy Pattern 代理模式 结构型设计模式
// 为其他对象提供一种代理以控制对这个对象的访问
// 适用于需要在访问对象时 添加额外操作的场景

type Image interface {
	Display()
}

type RealImage struct {
	filename string
}

func NewRealImage(filename string) *RealImage {
	// 模拟加载图片的耗时操作
	println("Loading image from disk:", filename)
	return &RealImage{filename: filename}
}

func (ri *RealImage) Display() {
	println("Displaying image:", ri.filename)
}

type ProxyImage struct {
	realImage *RealImage
	filename  string
}

func NewProxyImage(filename string) *ProxyImage {
	return &ProxyImage{filename: filename}
}

func (pi *ProxyImage) Display() {
	if pi.realImage == nil {
		pi.realImage = NewRealImage(pi.filename)
	}
	pi.realImage.Display()
}

func main() {
	var img Image = NewProxyImage("test_image.jpg")
	// 图像未加载时 代理不会加载真实图像
	println("First call to Display():")
	img.Display()
	// 图像已加载时 代理直接调用真实图像的Display方法
	println("Second call to Display():")
	img.Display()
}

// 优点: 控制对真实对象的访问 可以在不修改真实对象的情况下 添加额外功能
// 缺点: 增加了系统的复杂度 需要额外的代理类
// 适用场景: 需要在访问对象时 添加额外操作 如缓存 安全检查 日志记录等
// 对比装饰器模式: 装饰器模式关注于动态地给对象添加职责 而代理模式关注于控制对对象的访问
// 对比适配器模式: 适配器模式关注于接口转换 使得不兼容的接口可以协同工作 而代理模式关注于控制对对象的访问
// 举例: 虚拟代理 用于延迟加载资源 密集型资源如图像或文件 只有在真正需要时才加载
// Web领域例子: 缓存代理 用于缓存Web资源 减少对服务器的请求 提高响应速度
// 数据库领域例子: 连接池代理 用于管理数据库连接池 控制对数据库连接的访问 提高性能
