package main

// Abstract Factory Pattern 抽象工厂模式 创建型设计模式
// 提供一个创建一系列相关或相互依赖对象的接口 而无需指定它们具体的类
// 适用于产品族比较固定的场景

type Car interface {
	Drive()
}

type Bike interface {
	Ride()
}

type ToyotaCar struct {
}

func (tc *ToyotaCar) Drive() {
	println("Driving a Toyota car")
}

type ToyotaBike struct {
}

func (tb *ToyotaBike) Ride() {
	println("Riding a Toyota bike")
}

type BMWCar struct {
}

func (bc *BMWCar) Drive() {
	println("Driving a BMW car")
}

type BMWBike struct {
}

func (bb *BMWBike) Ride() {
	println("Riding a BMW bike")
}

type VehicleFactory interface {
	CreateCar() Car
	CreateBike() Bike
}

type ToyotaFactory struct {
}

func (tf *ToyotaFactory) CreateCar() Car {
	return &ToyotaCar{}
}

func (tf *ToyotaFactory) CreateBike() Bike {
	return &ToyotaBike{}
}

type BMWFactory struct {
}

func (bf *BMWFactory) CreateCar() Car {
	return &BMWCar{}
}

func (bf *BMWFactory) CreateBike() Bike {
	return &BMWBike{}
}

func main() {
	var factory VehicleFactory

	factory = &ToyotaFactory{}
	car1 := factory.CreateCar()
	car1.Drive()
	bike1 := factory.CreateBike()
	bike1.Ride()

	factory = &BMWFactory{}
	car2 := factory.CreateCar()
	car2.Drive()
	bike2 := factory.CreateBike()
	bike2.Ride()
}

// 优点: 遵循开闭原则 新增产品族只需要新增工厂类
// 缺点: 类的个数增加了 复杂度增加
// 适用场景: 需要创建一系列相关或相互依赖的对象 且产品族比较固定
// 例如: GUI 工具包 不同操作系统的按钮 文本框等控件 可以使用抽象工厂模式来创建不同操作系统的控件族
// 又例如: 数据库访问层 支持多种数据库 可以使用抽象工厂模式来创建不同数据库的连接 语句等对象族

func example() {

}
