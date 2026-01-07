package main

// Simple Factory Pattern 简单工厂模式 创建型设计模式
// 定义一个工厂类 根据传入的参数不同 返回不同类的实例
// 通过一个工厂类来创建对象 而不是直接new对象

type CarBrand string

const (
	ToyotaBrand CarBrand = "Toyota"
	BMWBrand    CarBrand = "BMW"
)

type Car interface {
	Drive()
}

type Toyota struct {
}

func (t *Toyota) Drive() {
	println("Driving a Toyota car")
}

type BMW struct {
}

func (b *BMW) Drive() {
	println("Driving a BMW car")
}

type CarFactory struct {
}

func (cf *CarFactory) CreateCar(brand CarBrand) Car {
	switch brand {
	case ToyotaBrand:
		return &Toyota{}
	case BMWBrand:
		return &BMW{}
	default:
		return nil
	}
	// 问题弊端: 每次新增品牌都需要修改工厂类代码 违反开闭原则
	// 解决方案: 使用反射机制 或者 注册表模式
}

func main() {
	factory := &CarFactory{}

	car1 := factory.CreateCar("Toyota")
	car1.Drive()

	car2 := factory.CreateCar("BMW")
	car2.Drive()
}
