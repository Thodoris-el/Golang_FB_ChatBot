package entity

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

//Structure of the Conversation Entity
type Conversation struct {
	ID          int64  `gorm:"primary_key;auto_increment" json:"id"`
	Facebook_id string `gorm:"not null;unique;" json:"facebook_id"`
	Type        string `gorm:"default: None" json:"type"`
	Stage       string `gorm:"default:None" json:"stage"`
	Product     string `gorm:"default:None" json:"product"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Save Conversation to DB
func (conversation *Conversation) SaveConversation(db *gorm.DB) (*Conversation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.WithContext(ctx).Debug().Create(&conversation).Error
	if err != nil {
		log.Println("error while saving conversation")
		return &Conversation{}, err
	}
	return conversation, err
}

//find all conversations fromDB
func (conversation *Conversation) FindAllConversations(db *gorm.DB) (*[]Conversation, error) {

	conversations := []Conversation{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.WithContext(ctx).Debug().Model(&Conversation{}).Limit(10).Find(&conversations).Error

	if err != nil {
		log.Println("Error while finding customers")
		return &[]Conversation{}, err
	}

	return &conversations, err
}

//Find By Id
func (conversation *Conversation) FindById(db *gorm.DB, R_id int64) (*Conversation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.WithContext(ctx).Debug().Model(&Conversation{}).Where("id = ?", R_id).Take(&conversation).Error

	if err != nil {
		log.Println("error while geting conversation by id")
		return &Conversation{}, err
	}

	return conversation, nil
}

//Find Conversation from a specific customer
func (conversation *Conversation) FindByCustomerId(db *gorm.DB, C_id string) (*Conversation, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.WithContext(ctx).Debug().Model(&Conversation{}).Where("facebook_id = ?", C_id).Take(&conversation).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &Conversation{}, err
		} else {
			log.Println("Error while finding conversation from a specific customer", err)
			return &Conversation{}, err
		}
	}
	return conversation, err
}

//Update Conversarion
func (conversation *Conversation) UpdateConversation(db *gorm.DB, C_id int64) (*Conversation, error) {

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	db = db.WithContext(ctx).Debug().Model(&Conversation{}).Where("id = ?", C_id).Take(&Conversation{}).UpdateColumns(
		map[string]interface{}{
			"type":       conversation.Type,
			"stage":      conversation.Stage,
			"product":    conversation.Product,
			"updated_at": time.Now(),
		},
	)
	err = db.WithContext(ctx).Debug().Model(&Conversation{}).Where("id = ?", C_id).Take(&conversation).Error
	if err != nil {
		return &Conversation{}, err
	}
	return conversation, nil
}

//delete conversation
func (conversation *Conversation) DeleteConversation(db *gorm.DB, C_id int64) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	db = db.WithContext(ctx).Debug().Model(&Conversation{}).Where("id = ?", C_id).Take(&Conversation{}).Delete(&Conversation{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
