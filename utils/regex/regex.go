package regex

import "regexp"

// MustEmail 检查给定的字符串是否是一个有效的电子邮件地址
func MustEmail(account string) bool {
	// 定义正则表达式模式
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_+&*-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`)
	// 使用MatchString方法检查字符串是否匹配正则表达式
	return pattern.MatchString(account)
}

func MustPhone(account string) bool {
	pattern := regexp.MustCompile("^1\\d{10}$")
	return pattern.MatchString(account)
}
