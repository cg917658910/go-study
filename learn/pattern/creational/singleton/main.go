package main

import "sync"

type UserService struct {
}

var (
	userServiceWithSingleton *UserService
	once                     sync.Once
)

func GetUserServiceInstance() *UserService {
	once.Do(func() {
		userServiceWithSingleton = &UserService{}
	})
	return userServiceWithSingleton
}
func main() {

}
