package models

import (
	"database/sql"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	ID         uint32         `gorm:"primary_key;auto_increment" json:"id"`
	Title      string         `gorm:"size:100;not null;" json:"title"`
	Deskripsi  string         `gorm:"type:text;size:1000;not null;" json:"deskripsi"`
	FileUpload sql.NullString `gorm:"size:255;" json:"file"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Todo) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Deskripsi = html.EscapeString(strings.TrimSpace(p.Deskripsi))
	// p.FileUpload = p.FileUpload
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Todo) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Deskripsi == "" {
		return errors.New("Required Deskripsi")
	}
	return nil
}

func (u *Todo) SaveTodo(db *gorm.DB) (*Todo, error) {

	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &Todo{}, err
	}
	return u, nil
}

func (u *Todo) FindAllTodos(db *gorm.DB) (*[]Todo, error) {
	var err error
	todos := []Todo{}
	err = db.Debug().Model(&Todo{}).Limit(100).Find(&todos).Error
	if err != nil {
		return &[]Todo{}, err
	}
	return &todos, err
}

func (u *Todo) FindTodoByID(db *gorm.DB, uid uint32) (*Todo, error) {
	var err error
	err = db.Debug().Model(Todo{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Todo{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Todo{}, errors.New("Todo Not Found")
	}
	return u, err
}

func (u *Todo) UpdateATodo(db *gorm.DB, uid uint32) (*Todo, error) {
	var err error
	db = db.Debug().Model(&Todo{}).Where("id = ?", uid).Take(&Todo{}).UpdateColumns(
		map[string]interface{}{
			"title":      u.Title,
			"deskripsi":  u.Deskripsi,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Todo{}, db.Error
	}
	// This is the display the updated todo
	err = db.Debug().Model(&Todo{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &Todo{}, err
	}
	return u, nil
}

func (u *Todo) DeleteATodo(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&Todo{}).Where("id = ?", uid).Take(&Todo{}).Delete(&Todo{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
