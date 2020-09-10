package generator

import (
	"github.com/appboot/appboot/internal/pkg/database"
	"github.com/gogf/gf/text/gstr"
)

// TableDesc 表结构详情
type TableDesc struct {
	ColumnName       string `gorm:"column:COLUMN_NAME"`    // 数据库原始字段
	OriMysqlType     string `gorm:"column:DATA_TYPE"`      // 数据库原始类型
	ColumnComment    string `gorm:"column:COLUMN_COMMENT"` // 备注
	ColumnKey        string `gorm:"column:COLUMN_KEY"`     // 是否是主键
	IsNullAble       string `gorm:"column:IS_NULLABLE"`    // 是否为空
	DefaultValue     string `gorm:"column:COLUMN_DEFAULT"` // 默认值
	ColumnTypeNumber string `gorm:"column:COLUMN_TYPE"`    // 类型(长度)
}

// CustomTableDesc 自定义表结构详情
type CustomTableDesc struct {
	Index          int
	ColumnName     string // 数据库原始字段
	GoColumnName   string // go使用的字段名称
	OriMysqlType   string // 数据库原始类型
	UpperMysqlType string // 转换大写的类型
	GolangType     string // 转换成golang类型
	//MysqlNullType    string // MYSQL对应的空类型
	PrimaryKey       bool   // 是否是主键
	IsNull           string // 是否为空
	DefaultValue     string // 默认值
	ColumnTypeNumber string // 类型(长度)
	ColumnComment    string // 备注
}

// GetTableDesc 获取表结构详情
func GetTableDesc(tableName string) (reply []*TableDesc, err error) {
	tableDesc := make([]*TableDesc, 0)
	dbName := database.GetDbName()
	if err := database.GetDB().
		Table("COLUMNS").
		Where("TABLE_NAME = ? and TABLE_SCHEMA = ?", tableName, dbName).
		Find(&tableDesc).Error; err != nil {
		return nil, err
	}
	return tableDesc, nil
}

// GenCustomTableDesc 生成自定义表结构
func GenCustomTableDesc(tableName string) ([]*CustomTableDesc, error) {
	desc, err := GetTableDesc(tableName)
	if err != nil {
		return nil, err
	}
	customTableDesc := make([]*CustomTableDesc, 0)
	for i, row := range desc {
		i++
		isPrimaryKey := false
		if row.ColumnKey == "PRI" {
			isPrimaryKey = true
		}
		customTableDesc = append(customTableDesc, &CustomTableDesc{
			Index:          i,
			GoColumnName:   gstr.CamelCase(row.ColumnName),
			UpperMysqlType: gstr.ToUpper(row.OriMysqlType),
			GolangType:     MysqlTypeToGoType[row.OriMysqlType],
			//MysqlNullType:    MysqlTypeToGoNullType[row.OriMysqlType],
			PrimaryKey:       isPrimaryKey,
			ColumnName:       row.ColumnName,
			OriMysqlType:     row.OriMysqlType,
			IsNull:           row.IsNullAble,
			DefaultValue:     row.DefaultValue,
			ColumnTypeNumber: row.ColumnTypeNumber,
			ColumnComment:    row.ColumnComment,
		})
	}
	return customTableDesc, nil
}
