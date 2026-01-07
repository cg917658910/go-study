package main

import "fmt"

// Prototype 原型接口
type Prototype interface {
	Clone() Prototype
}

// WorkExperience 工作经历
type WorkExperience struct {
	Company  string
	Position string
	Years    int
}

// Clone 克隆工作经历(深拷贝)
func (w *WorkExperience) Clone() *WorkExperience {
	return &WorkExperience{
		Company:  w.Company,
		Position: w.Position,
		Years:    w.Years,
	}
}

// Resume 简历
type Resume struct {
	Name       string
	Age        int
	Experience *WorkExperience
	Skills     []string
}

// Clone 克隆简历(深拷贝)
func (r *Resume) Clone() Prototype {
	// 克隆基本类型
	newResume := &Resume{
		Name: r.Name,
		Age:  r.Age,
	}

	// 深拷贝引用类型 - WorkExperience
	if r.Experience != nil {
		newResume.Experience = r.Experience.Clone()
	}

	// 深拷贝 slice
	if r.Skills != nil {
		newResume.Skills = make([]string, len(r.Skills))
		copy(newResume.Skills, r.Skills)
	}

	return newResume
}

// Display 显示简历信息
func (r *Resume) Display() {
	fmt.Printf("姓名: %s, 年龄: %d\n", r.Name, r.Age)
	if r.Experience != nil {
		fmt.Printf("工作经历: %s - %s (%d年)\n",
			r.Experience.Company,
			r.Experience.Position,
			r.Experience.Years)
	}
	if len(r.Skills) > 0 {
		fmt.Printf("技能: %v\n", r.Skills)
	}
	fmt.Println("---")
}

func main() {
	// 创建原型对象
	original := &Resume{
		Name: "张三",
		Age:  25,
		Experience: &WorkExperience{
			Company:  "ABC公司",
			Position: "软件工程师",
			Years:    2,
		},
		Skills: []string{"Go", "Python", "Java"},
	}

	fmt.Println("原始简历:")
	original.Display()

	// 克隆简历
	clone1 := original.Clone().(*Resume)
	clone1.Name = "李四"
	clone1.Age = 28
	clone1.Experience.Company = "XYZ公司"
	clone1.Experience.Position = "高级工程师"
	clone1.Experience.Years = 5
	clone1.Skills = append(clone1.Skills, "Rust")

	fmt.Println("克隆简历1 (修改后):")
	clone1.Display()

	// 再次克隆
	clone2 := original.Clone().(*Resume)
	clone2.Name = "王五"
	clone2.Age = 30

	fmt.Println("克隆简历2 (仅修改基本信息):")
	clone2.Display()

	// 验证原型未被修改
	fmt.Println("原始简历 (验证未被修改):")
	original.Display()

	// 演示浅拷贝的问题(仅作对比,不推荐)
	fmt.Println("\n=== 浅拷贝示例(不推荐) ===")
	shallowCopy := *original // 浅拷贝
	shallowCopy.Name = "赵六"
	shallowCopy.Experience.Company = "DEF公司" // 这会影响原对象!

	fmt.Println("浅拷贝后的简历:")
	fmt.Printf("姓名: %s, 年龄: %d, 公司: %s\n",
		shallowCopy.Name,
		shallowCopy.Age,
		shallowCopy.Experience.Company)

	fmt.Println("原始简历(被浅拷贝影响了):")
	original.Display()
}

// 优点: 提高性能 避免重复初始化 运行时动态配置
// 缺点: 需要为每个类实现克隆方法 深拷贝实现复杂 循环引用处理困难
// 适用场景: 对象创建成本高 需要创建大量相似对象 需要保护对象创建过程
// Go特性: 没有类继承 需要手动实现Clone方法 注意引用类型的深拷贝
// 举例: 文档编辑器的复制粘贴 游戏对象克隆 配置对象复制
// Web领域: API响应模板 缓存对象克隆
// 数据库领域: 数据库连接配置复制 查询对象克隆
