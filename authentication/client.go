package authentication

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect(connectionString string) (*gorm.DB, error) {
	instance, dbError := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
		return nil, dbError
	}
	log.Println("Connected to Database!")
	return instance, nil
}

func Migrate(instance *gorm.DB) error {
	err := instance.AutoMigrate(&Profesor{})
	err = instance.AutoMigrate(&Student{})
	err = instance.AutoMigrate(&Exam{})
	err = instance.AutoMigrate(&Clasa{})
	err = instance.AutoMigrate(&Exercitiu{})
	err = instance.AutoMigrate(&Calificativ{})
	if err != nil {
		return err
	}
	log.Println("Database Migration Completed!")
	return nil
}
