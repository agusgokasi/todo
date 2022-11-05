package models

import (
	"database/sql"
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type TodoList struct {
	ID         uint64         `gorm:"primary_key;auto_increment" json:"id"`
	Title      string         `gorm:"size:100;not null;" json:"title"`
	Deskripsi  string         `gorm:"type:text;size:1000;not null;" json:"deskripsi"`
	FileUpload sql.NullString `gorm:"size:255;" json:"file"`
	Todo       Todo           `json:"todo"`
	TodoID     uint32         `gorm:"not null" json:"todo_id"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *TodoList) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Deskripsi = html.EscapeString(strings.TrimSpace(p.Deskripsi))
	// p.FileUpload = p.FileUpload
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *TodoList) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}
	if p.Deskripsi == "" {
		return errors.New("Required Deskripsi")
	}
	if p.TodoID < 1 {
		return errors.New("Required TodoID")
	}
	return nil
}

func (p *TodoList) SaveTodoList(db *gorm.DB) (*TodoList, error) {

	var err error
	err = db.Debug().Create(&p).Error
	if err != nil {
		return &TodoList{}, err
	}
	return p, nil
}

func (p *TodoList) FindAllTodoLists(db *gorm.DB) (*[]TodoList, error) {
	var err error
	todolists := []TodoList{}
	err = db.Debug().Model(&TodoList{}).Limit(100).Find(&todolists).Error
	if err != nil {
		return &[]TodoList{}, err
	}
	return &todolists, err
}

func (p *TodoList) FindTodoListByID(db *gorm.DB, pid uint64) (*TodoList, error) {
	var err error
	err = db.Debug().Model(TodoList{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &TodoList{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &TodoList{}, errors.New("TodoList Not Found")
	}
	return p, err
}

func (p *TodoList) UpdateATodoList(db *gorm.DB, pid uint64) (*TodoList, error) {
	var err error
	db = db.Debug().Model(&TodoList{}).Where("id = ?", pid).Take(&TodoList{}).UpdateColumns(
		map[string]interface{}{
			"title":      p.Title,
			"deskripsi":  p.Deskripsi,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &TodoList{}, db.Error
	}
	// This is the display the updated todolist
	err = db.Debug().Model(&TodoList{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &TodoList{}, err
	}
	return p, nil
}

func (p *TodoList) DeleteATodoList(db *gorm.DB, pid uint64) (int64, error) {

	db = db.Debug().Model(&TodoList{}).Where("id = ?", pid).Take(&TodoList{}).Delete(&TodoList{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
