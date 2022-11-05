package seed

import (
	"log"
	"todo/api/models"

	"github.com/jinzhu/gorm"
)

var todos = []models.Todo{
	models.Todo{
		Title:     "Test Title",
		Deskripsi: "just testing",
	},
	models.Todo{
		Title:     "Test Title 2",
		Deskripsi: "just testing 2",
	},
}

var todolist = []models.TodoList{
	models.TodoList{
		Title:     "Title 1",
		Deskripsi: "Hello world 1",
	},
	models.TodoList{
		Title:     "Title 2",
		Deskripsi: "Hello world 2",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.TodoList{}, &models.Todo{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Todo{}, &models.TodoList{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.TodoList{}).AddForeignKey("todo_id", "todos(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range todos {
		err = db.Debug().Model(&models.Todo{}).Create(&todos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed todos table: %v", err)
		}
		todolist[i].TodoID = todos[i].ID

		err = db.Debug().Model(&models.TodoList{}).Create(&todolist[i]).Error
		if err != nil {
			log.Fatalf("cannot seed todolist table: %v", err)
		}
	}
}
