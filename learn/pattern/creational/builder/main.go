package main

// Builder Pattern 建造者模式 创建型设计模式
// 复制对象 希望一步一步创建
// 通过多个简单的对象一步一步构建成一个复杂的对象

type Car struct {
	Brand    string
	Color    string
	Speed    int
	Capacity int
}

type CarBuilder struct {
	car *Car
}

func (car *Car) Drive() {
	println("Driving a", car.Brand, "car with speed", car.Speed)
}

func NewCarBuilder() *CarBuilder {
	return &CarBuilder{car: &Car{}}
}

func (b *CarBuilder) SetBrand(brand string) *CarBuilder {
	b.car.Brand = brand
	return b
}

func (b *CarBuilder) SetColor(color string) *CarBuilder {
	b.car.Color = color
	return b
}

func (b *CarBuilder) SetSpeed(speed int) *CarBuilder {
	b.car.Speed = speed
	return b
}

func (b *CarBuilder) SetCapacity(capacity int) *CarBuilder {
	b.car.Capacity = capacity
	return b
}

func (b *CarBuilder) Build() *Car {
	return b.car
}

func main() {
	builder := NewCarBuilder()
	car := builder.SetBrand("Toyota").
		SetColor("Red").
		SetSpeed(200).
		SetCapacity(5).
		Build()

	// 使用 car 对象
	car.Drive()
}

// 优点: 将复杂对象的构建过程与表示分离 使得同样的构建过程可以创建不同的表示
// 缺点: 产生多余的建造者类 增加系统复杂度
// 适用场景: 需要一步一步构建复杂对象 需要生成不同表示的对象
// 对比工厂方法模式: 工厂方法模式关注创建单一产品 而建造者模式关注创建复杂产品
// 对比抽象工厂模式: 抽象工厂模式用于创建相关或依赖对象的家族 而建造者模式用于一步步构建复杂对象
// 对比原型模式: 原型模式通过复制现有对象来创建新对象 而建造者模式通过一步步构建来创建新对象
// 举例: 创建一辆汽车 需要设置品牌 颜色 速度 容量等属性 使用建造者模式可以一步步设置这些属性 最后构建出完整的汽车对象
// Web领域例子: 构建一个复杂的HTML页面 通过建造者模式一步步添加标题 段落 图片等元素 最后生成完整的HTML页面
// 数据库领域例子: 构建一个复杂的SQL查询语句 通过建造者模式一步步添加选择字段 条件 排序等 最后生成完整的SQL语句
