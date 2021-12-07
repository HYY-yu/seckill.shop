package mysqlerr_helper

import (
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

// err helpers

// IsMysqlDupEntryError Mysql返回唯一索引冲突
func IsMysqlDupEntryError(err error) bool {
	if driverErr, ok := err.(*mysql.MySQLError); ok {
		if driverErr.Number == mysqlerr.ER_DUP_ENTRY {
			// 唯一索引错误
			return true
		}
	}
	return false
}
