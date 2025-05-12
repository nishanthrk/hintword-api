package user_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/datatypes"
	db "hintword.com/api/app/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"hintword.com/api/app/common/utility"
	"hintword.com/api/app/models"

	cfg "hintword.com/api/app/configs"
)

type UserClaims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

type UserDetails struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}

type GoogleUserInfo struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar" validate:"required"`
	GoogleID string `json:"google_id" validate:"required"`
}

func HandleGoogleUserInfo(params GoogleUserInfo) (string, error) {
	var user models.Users
	user, _ = user.FindByEmail(params.Email)
	newUser := false
	if user.UserID == "" {
		user = models.Users{
			UserID:   utility.GenerateUUID(),
			Email:    params.Email,
			Name:     params.Name,
			GoogleID: params.GoogleID,
			Status:   models.StatusActive,
		}
		newUser = true
	}
	user.AvatarURL = params.Avatar

	_, err := user.Save()
	if err != nil {
		return "", fmt.Errorf("failed to save user: %v", err)
	}

	// Generate JWT token
	jwtToken, err := GenerateToken(user.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	if newUser {
		err = PopulateData(user.UserID)
		if err != nil {
			return "", fmt.Errorf("failed to populate data: %v", err)
		}
	}

	return jwtToken, nil
}

func GetUserObject(c interface{}) (details UserDetails) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("recovered from panic %v", r)
		}
		return
	}()

	var user *jwt.Token
	if ctx, ok := c.(*fiber.Ctx); ok {
		user = ctx.Locals("user").(*jwt.Token)
	} else if wCtx, _ok := c.(*websocket.Conn); _ok {
		user = wCtx.Locals("user").(*jwt.Token)
	} else {
		fmt.Println("Content type not found unable to fetch user")
		return
	}

	claims := user.Claims.(jwt.MapClaims)
	jsonData, err := utility.MapToJSON(claims)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return
	}
	err = utility.JSONToStruct(jsonData, &details)
	if err != nil {
		fmt.Println("Error converting JSON to struct:", err)
		return
	}
	return
}

func GenerateToken(userID string) (string, error) {
	// Get user details from database
	var user models.Users
	user, _ = user.FindById(userID)
	if user.UserID == "" {
		return "", errors.New("user not found")
	}

	expireTime := time.Now().Add(time.Hour * 24)

	claims := UserClaims{
		user.UserID,
		user.Email,
		user.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			Issuer:    cfg.GetConfig().JWTIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.GetConfig().JWTAccessSecret))
}

func PopulateData(userID string) error {
	tx := db.MysqlDB.Begin()
	note := models.Notes{}
	note.NoteID = utility.GenerateUUID()
	note.UserID = userID
	note.Title = "🚀 Welcome to Hintword!"
	rawContent := `{"type":"doc","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Thanks for trying out Hintword – your AI-powered notepad for smarter writing, brainstorming, and productivity."}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Whether you're capturing quick ideas, summarizing a meeting, or drafting long-form content, Hintword is here to assist with intelligent suggestions and seamless formatting."}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"heading","attrs":{"textAlign":null,"level":3},"content":[{"type":"text","text":"✨ Features at a Glance"}]},{"type":"blockquote","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"A minimal yet powerful rich text editor with smart tools powered by AI."}]}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"orderedList","attrs":{"start":1,"type":null},"content":[{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"✅ "},{"type":"text","marks":[{"type":"bold"}],"text":"Auto Suggestions:"},{"type":"text","text":" Activate grammar improvement, rephrasing, or tone adjustments with a single click."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"🎙 "},{"type":"text","marks":[{"type":"bold"}],"text":"Voice-to-Text"},{"type":"text","text":" (coming soon): Dictate your notes on the go."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"🧠 "},{"type":"text","marks":[{"type":"bold"}],"text":"AI Note Optimizer"},{"type":"text","text":": Let Hintword structure, improve, or summarize your content."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"🖼 "},{"type":"text","marks":[{"type":"bold"}],"text":"Embed & Organize"},{"type":"text","text":": Add tags, insert links, and organize your notes for easy retrieval."}]}]}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"✍️ How to use"}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"bulletList","content":[{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Just start typing your idea."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Use the "},{"type":"text","marks":[{"type":"bold"}],"text":"🧠 AI"},{"type":"text","text":" icon to optimize."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Tap "},{"type":"text","marks":[{"type":"bold"}],"text":"🔊"},{"type":"text","text":" for speech-related tools (coming soon)."}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"Keep track of all your notes from "},{"type":"text","marks":[{"type":"bold"}],"text":"Recent Notes"},{"type":"text","text":" on the right panel."}]}]}]},{"type":"paragraph","attrs":{"textAlign":null}},{"type":"heading","attrs":{"textAlign":null,"level":3},"content":[{"type":"text","text":"💡 Example Use Cases"}]},{"type":"bulletList","content":[{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"📝 Meeting Minutes"}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"🎯 Goal Planning"}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"📚 Study Notes"}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"💭 Brainstorming Ideas"}]}]},{"type":"listItem","content":[{"type":"paragraph","attrs":{"textAlign":null},"content":[{"type":"text","text":"🧾 Drafting Emails or Reports"}]},{"type":"paragraph","attrs":{"textAlign":null}}]}]}]}`

	var content datatypes.JSON
	err := json.Unmarshal([]byte(rawContent), &content)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to unmarshal note content: %v", err)
	}
	note.Content = content

	if err = tx.Create(&note).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save note: %v", err)
	}

	collection := models.Collections{}
	collection.CollectionID = utility.GenerateUUID()
	collection.Name = "Default"
	collection.UserID = userID

	if err = tx.Create(&collection).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create default collection: %v", err)
	}

	now := time.Now()

	tabs := []models.Tabs{
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Google", URL: "https://www.google.com", FaviconURL: "https://www.google.com/favicon.ico", Sequence: 1, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Whatsapp", URL: "https://web.whatsapp.com", FaviconURL: "https://web.whatsapp.com/favicon.ico", Sequence: 2, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "ChatGPT", URL: "https://chat.openai.com", FaviconURL: "https://chat.openai.com/favicon.ico", Sequence: 3, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Facebook", URL: "https://www.facebook.com", FaviconURL: "https://facebook.com/favicon.ico", Sequence: 4, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Youtube", URL: "https://www.youtube.com", FaviconURL: "https://www.youtube.com/favicon.ico", Sequence: 5, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Google Mail", URL: "https://mail.google.com", FaviconURL: "https://ssl.gstatic.com/ui/v1/icons/mail/rfr/gmail.ico", Sequence: 6, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Spotify", URL: "https://open.spotify.com", FaviconURL: "https://open.spotify.com/favicon.ico", Sequence: 7, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Linkedin", URL: "https://www.linkedin.com", FaviconURL: "https://www.linkedin.com/favicon.ico", Sequence: 8, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Hintword", URL: "https://hintword.com", FaviconURL: "https://www.hintword.com/favicon.ico", Sequence: 9, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
		{TabID: utility.GenerateUUID(), UserID: userID, CollectionID: collection.CollectionID, Title: "Pub.dev", URL: "https://pub.dev", FaviconURL: "https://pub.dev/favicon.ico", Sequence: 10, Status: models.StatusActive, CreatedAt: now, UpdatedAt: now},
	}

	if err = tx.Create(&tabs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save default tabs: %v", err)
	}

	return tx.Commit().Error
}
