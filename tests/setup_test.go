package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	controllers "github.com/Thodoris-el/Coonectly-Project/api/Controllers"
	entity "github.com/Thodoris-el/Coonectly-Project/api/Entity"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var server = controllers.Server{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatal("Error getting env \n", err)
	}

	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	Db_user := os.Getenv("TEST_DB_USERNAME")
	Db_password := os.Getenv("TEST_DB_PASSWORD")
	Db_host := os.Getenv("TEST_DB_HOST")
	Db_port := os.Getenv("TEST_DB_PORT")
	Db_name := os.Getenv("TEST_DB_NAME")

	dsn := Db_user + ":" + Db_password + "@tcp" + "(" + Db_host + ":" + Db_port + ")/" + Db_name + "?" + "parseTime=true&loc=Local"

	server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("cant connect to mysql")
		log.Fatal(err)
	} else {
		fmt.Printf("connected to mysql")
	}
}

func refreshTables() {
	server.DB.Migrator().DropTable(&entity.Customer{})
	server.DB.AutoMigrate(&entity.Customer{})

	server.DB.Migrator().DropTable(&entity.Review{})
	server.DB.AutoMigrate(&entity.Review{})

	server.DB.Migrator().DropTable(&entity.Conversation{})
	server.DB.AutoMigrate(&entity.Conversation{})

	server.DB.Migrator().DropTable(&entity.Template{})
	server.DB.AutoMigrate(&entity.Template{})

}

func createATemplate() (entity.Template, error) {
	template := entity.Template{
		Title:        "How would you rate our product?",
		Placeholder:  "Give additional feedback",
		Subtitle:     "Let us know about our product?",
		Button_Title: "Review",
		Language:     "eng",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := server.DB.Model(&entity.Template{}).Create(&template).Error
	if err != nil {
		return entity.Template{}, err
	}
	return template, nil
}

func createACustomer() (entity.Customer, error) {
	customer := entity.Customer{
		First_name:  "John",
		Last_name:   "Dir",
		Facebook_id: "6706612322695175",
		Language:    "eng",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := server.DB.Model(&entity.Customer{}).Create(&customer).Error
	if err != nil {
		return entity.Customer{}, err
	}
	return customer, nil
}

func createMessages(msg []string) []entity.FacebookMessage {

	var fbMessages []entity.FacebookMessage
	for _, mess := range msg {
		testMessaging := entity.MessagingType{
			Sender:    entity.SenderType{ID: "6706612322695175"},
			Recipient: entity.RecipientType{ID: "101209549545169"},
			Timestamp: 1673894538114,
			Message: entity.MessageType{
				Mid:  "m_6EjoW5-Hh83gn-XI7xyamUqIpuPv8gV21Q6xzapWR7EyfSaOaVai6BEaa0qyeYNA24MQugba4W2YxajYKhMzsQ",
				Text: mess,
			},
		}

		testmess := []entity.MessagingType{}
		testmess = append(testmess, testMessaging)

		testEnty := entity.EntryType{
			ID:        "101209549545169",
			Time:      1673894538501,
			Messaging: testmess,
		}
		testentinies := []entity.EntryType{}
		testentinies = append(testentinies, testEnty)

		reqBody := entity.FacebookMessage{
			Object: "page",
			Entry:  testentinies,
		}

		fbMessages = append(fbMessages, reqBody)
	}
	return fbMessages
}

func createMessageReview() entity.FacebookMessage {

	feed := entity.FeScType{
		ScreenID: 0,
		Questions: entity.QuesTypeRes{
			Myquestion1: entity.MyQuestionType{
				Type:    "csat",
				Payload: "3",
				FollowUp: entity.FollowUpTypeRes{
					Type:    "free_form",
					Payload: "very good",
				},
			},
		},
	}

	var feeds []entity.FeScType
	feeds = append(feeds, feed)

	testMessaging := entity.MessagingType{
		Sender:    entity.SenderType{ID: "6706612322695175"},
		Recipient: entity.RecipientType{ID: "101209549545169"},
		Timestamp: 1673894538114,
		Message: entity.MessageType{
			Mid: "m_6EjoW5-Hh83gn-XI7xyamUqIpuPv8gV21Q6xzapWR7EyfSaOaVai6BEaa0qyeYNA24MQugba4W2YxajYKhMzsQ",
		},
		Messaging_Feedback: entity.MesFeedType{
			FeedbackScreens: feeds,
		},
	}

	testmess := []entity.MessagingType{}
	testmess = append(testmess, testMessaging)

	testEnty := entity.EntryType{
		ID:        "101209549545169",
		Time:      1673894538501,
		Messaging: testmess,
	}
	testentinies := []entity.EntryType{}
	testentinies = append(testentinies, testEnty)

	reqBody := entity.FacebookMessage{
		Object: "page",
		Entry:  testentinies,
	}

	return reqBody
}

func createTwoCustomers() ([]entity.Customer, error) {
	customers := []entity.Customer{
		{
			First_name:  "John",
			Last_name:   "Dir",
			Facebook_id: "6706612322695175",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			First_name:  "Nick",
			Last_name:   "Dir",
			Facebook_id: "8006612322695175",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	for _, customer := range customers {
		err := server.DB.Model(&entity.Customer{}).Create(&customer).Error
		if err != nil {
			return []entity.Customer{}, err
		}
	}

	return customers, nil
}

func createAConversation() (entity.Conversation, error) {
	conversation := entity.Conversation{
		Facebook_id: "6706612322695175",
		Type:        "Review",
		Stage:       "None",
		Product:     "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := server.DB.Model(&entity.Conversation{}).Create(&conversation).Error
	if err != nil {
		return entity.Conversation{}, err
	}

	return conversation, nil
}

func createTwoConversations() ([]entity.Conversation, error) {
	conversations := []entity.Conversation{
		{
			Facebook_id: "6706612322695175",
			Type:        "Review",
			Stage:       "None",
			Product:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Facebook_id: "8006612322695175",
			Type:        "",
			Stage:       "None",
			Product:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	for _, conversation := range conversations {
		err := server.DB.Model(&entity.Conversation{}).Create(&conversation).Error
		if err != nil {
			return []entity.Conversation{}, err
		}
	}

	return conversations, nil
}

func createAReview() (entity.Review, error) {
	review := entity.Review{
		Customer_id: "6706612322695175",
		Text:        "1",
		Score:       1,
		Product:     "car",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := server.DB.Model(&entity.Review{}).Create(&review).Error
	if err != nil {
		return entity.Review{}, err
	}

	return review, nil
}

func createTwoReviews() ([]entity.Review, error) {
	reviews := []entity.Review{
		{
			Customer_id: "6706612322695175",
			Text:        "1",
			Score:       1,
			Product:     "car",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Customer_id: "8006612322695175",
			Text:        "2",
			Score:       2,
			Product:     "bike",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	for _, review := range reviews {
		err := server.DB.Model(&entity.Review{}).Create(&review).Error
		if err != nil {
			return []entity.Review{}, err
		}
	}

	return reviews, nil
}

func createTwoTemplates() ([]entity.Template, error) {
	template := []entity.Template{
		{
			Title:       "How would you rate our product?",
			Placeholder: "Give additional feedback",
			Language:    "eng",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			Title:       "How would you rate our product?",
			Placeholder: "Give additional feedback",
			Language:    "eng",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
	err := server.DB.Model(&entity.Template{}).Create(&template).Error
	if err != nil {
		return []entity.Template{}, err
	}
	return template, nil
}
