package database

import (
	"database/sql"
	"fmt"
)

// ResultCheck 检查 result 是否存在错误
func ResultCheck(result sql.Result) error {
	_, err := result.RowsAffected()
	return err
}

func ResultCheckErr(result sql.Result, err error) error {
	if err == nil {
		_, err = result.RowsAffected()
	}

	return err
}

// ResultExpectN 检查 result、err 是否存在错误，并检查本次数据库操作影响的行数是否与 n 匹配
func ResultExpectN(result sql.Result, err error, n int) error {
	_, err2 := ResultExpectNBut(result, err, n)
	return err2
}

// ResultExpectNBut 检查 result、err 是否存在错误，并检查本次数据库操作影响的行数是否与 n 匹配, 并返回实际匹配的行数
func ResultExpectNBut(result sql.Result, err error, n int) (int, error) {
	var effect int = 0
	var effectN int64 = 0
	if err == nil {
		effectN, err = result.RowsAffected()
		effect = int(effectN)
	}

	if err != nil {
		return effect, fmt.Errorf("getting rows affected from result with error %s", err)
	}

	if effect != n {
		return effect, fmt.Errorf("rows affected doesn't math expected %d", n)
	}

	return effect, nil
}

// ResultGreaterN 检查 result、err 是否存在错误，并检查本次数据库操作影响的行数是否大于 n
func ResultGreaterN(result sql.Result, err error, n int) error {
	var effect int64 = 0
	if err == nil {
		effect, err = result.RowsAffected()
	}

	if err != nil {
		return fmt.Errorf("getting rows affected from result with error %s", err)
	}

	if int(effect) <= n {
		return fmt.Errorf("rows affected less then expected %d", n)
	}

	return nil
}
