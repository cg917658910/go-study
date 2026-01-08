package main

import (
	"testing"
)

func TestResumeClone(t *testing.T) {
	// 创建原型
	original := &Resume{
		Name: "张三",
		Age:  25,
		Experience: &WorkExperience{
			Company:  "ABC公司",
			Position: "软件工程师",
			Years:    2,
		},
		Skills: []string{"Go", "Python"},
	}

	// 克隆
	clone := original.Clone().(*Resume)

	// 修改克隆对象
	clone.Name = "李四"
	clone.Age = 30
	clone.Experience.Company = "XYZ公司"
	clone.Skills = append(clone.Skills, "Java")

	// 验证原对象未被修改
	if original.Name != "张三" {
		t.Errorf("原对象的Name被修改了: got %s, want 张三", original.Name)
	}

	if original.Age != 25 {
		t.Errorf("原对象的Age被修改了: got %d, want 25", original.Age)
	}

	if original.Experience.Company != "ABC公司" {
		t.Errorf("原对象的Company被修改了: got %s, want ABC公司", original.Experience.Company)
	}

	if len(original.Skills) != 2 {
		t.Errorf("原对象的Skills长度被修改了: got %d, want 2", len(original.Skills))
	}

	// 验证克隆对象的值
	if clone.Name != "李四" {
		t.Errorf("克隆对象的Name错误: got %s, want 李四", clone.Name)
	}

	if clone.Experience.Company != "XYZ公司" {
		t.Errorf("克隆对象的Company错误: got %s, want XYZ公司", clone.Experience.Company)
	}
}

func TestWorkExperienceClone(t *testing.T) {
	original := &WorkExperience{
		Company:  "ABC公司",
		Position: "工程师",
		Years:    3,
	}

	clone := original.Clone()

	// 修改克隆对象
	clone.Company = "XYZ公司"
	clone.Years = 5

	// 验证原对象未被修改
	if original.Company != "ABC公司" {
		t.Errorf("原对象的Company被修改了: got %s, want ABC公司", original.Company)
	}

	if original.Years != 3 {
		t.Errorf("原对象的Years被修改了: got %d, want 3", original.Years)
	}
}

func TestResumeCloneNilFields(t *testing.T) {
	// 测试nil字段的克隆
	original := &Resume{
		Name:       "张三",
		Age:        25,
		Experience: nil,
		Skills:     nil,
	}

	clone := original.Clone().(*Resume)

	if clone.Experience != nil {
		t.Error("克隆对象的Experience应该为nil")
	}

	if clone.Skills != nil {
		t.Error("克隆对象的Skills应该为nil")
	}
}

func TestResumeCloneIndependence(t *testing.T) {
	// 创建原型
	original := &Resume{
		Name: "张三",
		Age:  25,
		Experience: &WorkExperience{
			Company:  "ABC公司",
			Position: "工程师",
			Years:    2,
		},
		Skills: []string{"Go"},
	}

	// 创建多个克隆
	clone1 := original.Clone().(*Resume)
	clone2 := original.Clone().(*Resume)

	// 修改clone1
	clone1.Experience.Company = "XYZ公司"
	clone1.Skills[0] = "Python"

	// 验证clone2和original未受影响
	if clone2.Experience.Company != "ABC公司" {
		t.Error("clone2不应受clone1的修改影响")
	}

	if original.Skills[0] != "Go" {
		t.Error("original不应受clone1的修改影响")
	}
}

func BenchmarkResumeClone(b *testing.B) {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = original.Clone()
	}
}

func BenchmarkResumeNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &Resume{
			Name: "张三",
			Age:  25,
			Experience: &WorkExperience{
				Company:  "ABC公司",
				Position: "软件工程师",
				Years:    2,
			},
			Skills: []string{"Go", "Python", "Java"},
		}
	}
}
