package mian

import "context"

// Decorator Pattern 装饰器模式 结构型设计模式
// 动态地给一个对象添加一些额外的职责 就增加功能来说 装饰器模式比生成子类更为灵活
// 适用于需要在不影响其他对象的情况下 给单个对象添加功能的场景

type Car interface {
	Drive()
}

type BasicCar struct {
}

func (bc *BasicCar) Drive() {
	println("Driving a basic car")
}

type CarDecorator struct {
	car Car
}

func (cd *CarDecorator) Drive() {
	cd.car.Drive()
}

type SportsCar struct {
	*CarDecorator
}

func NewSportsCar(car Car) *SportsCar {
	return &SportsCar{&CarDecorator{car: car}}
}

func (sc *SportsCar) Drive() {
	println("Adding sports features")
	sc.CarDecorator.Drive()
}

type LuxuryCar struct {
	*CarDecorator
}

func NewLuxuryCar(car Car) *LuxuryCar {
	return &LuxuryCar{&CarDecorator{car: car}}
}

func (lc *LuxuryCar) Drive() {
	println("Adding luxury features")
	lc.CarDecorator.Drive()
}

func main() {
	var car Car = &BasicCar{}
	car.Drive()

	car = NewSportsCar(car)
	car.Drive()

	car = NewLuxuryCar(car)
	car.Drive()
}

// 优点: 遵循开闭原则 可以在不修改现有代码的情况下扩展对象功能
// 缺点: 可能会产生大量的小类 增加系统复杂度
// 适用场景: 需要在不影响其他对象的情况下 给单个对象添加功能 需要动态地添加或删除功能
// 对比继承: 继承在编译时确定功能 而装饰器模式在运行时动态添加功能 更加灵活
// 对比代理模式: 代理模式主要用于控制对对象的访问 而装饰器模式主要用于扩展对象的功能
// 举例: 给汽车添加不同的功能 如运动套件 豪华套件等 使用装饰器模式可以动态地组合这些功能
// Web领域例子: 给HTTP请求添加不同的中间件 如日志记录 鉴权等 使用装饰器模式可以动态地组合这些中间件
// 数据库领域例子: 给数据库连接添加不同的功能 如缓存 事务等 使用装饰器模式可以动态地组合这些功能

// 函数编程领域例子: 给函数添加不同的装饰器 如缓存 日志等 使用装饰器模式可以动态地组合这些装饰器

type Object func(int) int

func Decorator1(obj Object) Object {
	return func(a int) int {
		println("Decorator1 before")
		result := obj(a)
		println("Decorator1 after")
		return result
	}
}

func Decorator2(obj Object) Object {
	return func(a int) int {
		println("Decorator2 before")
		result := obj(a)
		println("Decorator2 after")
		return result
	}
}

func BaseFunction(a int) int {
	println("BaseFunction:", a)
	return a * 2
}

func main2() {
	var obj Object = BaseFunction
	obj = Decorator1(obj)
	obj = Decorator2(obj)

	result := obj(10)
	println("Result:", result)
}

// 优雅实现中间件

type HandlerFunc func(ctx context.Context)

func Middleware1(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context) {
		println("Middleware1 before")
		next(ctx)
		println("Middleware1 after")
	}
}

func Middleware2(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context) {
		println("Middleware2 before")
		next(ctx)
		println("Middleware2 after")
	}
}

func FinalHandler(ctx context.Context) {
	println("FinalHandler executed")
}

func main3() {
	var handler HandlerFunc = FinalHandler
	handler = Middleware1(handler)
	handler = Middleware2(handler)

	handler(context.Background())
}

// 实现一个优雅重试机制
func RetryMiddleware(retries int, next HandlerFunc) HandlerFunc {
	return func(ctx context.Context) {
		for i := 0; i < retries; i++ {
			println("Retry attempt:", i+1)
			next(ctx)
		}
	}
}

func main4() {
	var handler HandlerFunc = FinalHandler
	handler = RetryMiddleware(3, handler)

	handler(context.Background())
}
