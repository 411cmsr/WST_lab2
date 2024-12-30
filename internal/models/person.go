package models

type Person struct {
	ID        uint   `gorm:"primaryKey; not null" xml:"id,omitempty" yaml:"id,omitempty"`
	Name      string `gorm:"type:varchar(200)" xml:"name" yaml:"name"`
	Surname   string `gorm:"type:varchar(200)" xml:"surname" yaml:"surname"`
	Age       int    `gorm:"age,omitempty" xml:"age" yaml:"age"`
	Email     string `gorm:"type:varchar(200); uniqueIndex; not null" xml:"email" yaml:"email"`
	Telephone string `gorm:"type:varchar(200); not null" xml:"telephone" yaml:"telephone"`
}
