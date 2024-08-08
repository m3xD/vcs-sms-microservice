package repo

import "gorm.io/gorm"

type Database interface {
	Find(out interface{}, where ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Order(value interface{}) *gorm.DB
	Table(name string, o ...interface{}) *gorm.DB
	Update(column string, value interface{}) *gorm.DB
	Limit(limit int) *gorm.DB
	Offset(offset int) *gorm.DB
	Preload(column string, conditions ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
	Association(column string) *gorm.Association
}
