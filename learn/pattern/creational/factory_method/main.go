package main

// Factory Method Pattern 工厂方法模式 创建型设计模式
// 定义一个创建对象的接口 由子类决定实例化哪一个类
// 将类的实例化推迟到子类

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

type CarFactory interface {
	CreateCar() Car
}

type ToyotaFactory struct {
}

func (tf *ToyotaFactory) CreateCar() Car {
	return &Toyota{}
}

type BMWFactory struct {
}

func (bf *BMWFactory) CreateCar() Car {
	return &BMW{}
}

func main() {
	var factory CarFactory

	factory = &ToyotaFactory{}
	car1 := factory.CreateCar()
	car1.Drive()

	factory = &BMWFactory{}
	car2 := factory.CreateCar()
	car2.Drive()
}

// 优点: 遵循开闭原则 新增品牌只需要新增工厂类
// 缺点: 类的个数增加了 复杂度增加
// 适用场景: 需要创建的对象比较复杂 需要通过子类来决定创建哪一个类的实例
// 对比注册表工厂模式: 注册表工厂模式通过注册和查找函数指针来创建对象 更加灵活 但也更复杂
// 对比简单工厂模式: 简单工厂模式集中管理对象创建 但违反开闭原则 不易扩展
// 对比抽象工厂模式: 抽象工厂模式用于创建相关或依赖对象的家族 而工厂方法模式用于创建单一对象
