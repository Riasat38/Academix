package main

import "fmt"

var studentPermissions map[string][]string
var teacherPermissions map[string][]string

// Initialize maps inside init() function
func main() {
	studentPermissions = map[string][]string{
		"course":  {"viewAll", "viewOwn", "enroll"},
		"profile": {"view", "edit"},
	}

	teacherPermissions = map[string][]string{
		"course":  {"viewAll", "viewOwn", "modify"},
		"profile": {"view", "edit"},
	}
	fmt.Printf("%T\n", studentPermissions)
	fmt.Println(teacherPermissions["course"])
}
