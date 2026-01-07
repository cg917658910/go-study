package main

import "fmt"

// 定义 Car 接口和具体实现
type Car interface {
	Drive()
}

type Toyota struct{}

func (t *Toyota) Drive() {
	fmt.Println("Driving a Toyota car")
}

type BMW struct{}

func (b *BMW) Drive() {
	fmt.Println("Driving a BMW car")
}

// 注册表工厂
type CarFactory struct {
	registry map[string]func() Car
}

// 注册新品牌
func (cf *CarFactory) Register(brand string, constructor func() Car) {
	cf.registry[brand] = constructor
}

// 创建汽车
func (cf *CarFactory) CreateCar(brand string) Car {
	if constructor, exists := cf.registry[brand]; exists {
		return constructor()
	}
	return nil
}

func main() {
	// 初始化工厂
	factory := &CarFactory{registry: make(map[string]func() Car)}

	// 注册品牌
	factory.Register("Toyota", func() Car { return &Toyota{} })
	factory.Register("BMW", func() Car { return &BMW{} })

	// 创建汽车
	car1 := factory.CreateCar("Toyota")
	car1.Drive()

	car2 := factory.CreateCar("BMW")
	car2.Drive()
}
