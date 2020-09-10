package generator

import (
	"sort"
	"strings"

	"github.com/appboot/appboot/internal/pkg/database"
)

// TableNameAndComment
type TableNameAndComment struct {
	Index        int    // 索引
	TableName    string `gorm:"column:TABLE_NAME"`    // 表名
	TableComment string `gorm:"column:TABLE_COMMENT"` // 注释
}

// FindDbTables 获取数据库所有表
func FindDbTables() ([]*TableNameAndComment, error) {
	// 获取表名和注释
	var nameAndComments []*TableNameAndComment
	dbName := database.GetDbName()
	// if err := database.GetDB().Table("tables").Select("table_name,table_comment").Where("table_schema = ?", dbName).Find(&nameAndComments).Error; err != nil {
	if err := database.GetDB().Table("tables").Where("table_schema = ?", dbName).Find(&nameAndComments).Error; err != nil {
		return nil, err
	}

	// 添加索引
	for idx, info := range nameAndComments {
		idx++
		info.Index = idx
	}
	//排序, 采用升序
	sort.Slice(nameAndComments, func(i, j int) bool {
		return strings.ToLower(nameAndComments[i].TableName) < strings.ToLower(nameAndComments[j].TableName)
	})
	return nameAndComments, nil
}

// ListTableNameAndComment 返回表名和评论列表
func ListTableNameAndComment() (result []string, err error){
	nameAndComments, err := FindDbTables()
	if err != nil {
		return  nil, err
	}
	for _, v := range nameAndComments {
		result = append(result, v.TableName)
	}
	return result, nil
}
