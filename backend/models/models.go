package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Nickname  string         `json:"nickname"`
	Avatar    string         `json:"avatar"`
	Role      string         `json:"role" gorm:"default:admin"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Nickname  string         `json:"nickname"`
	Avatar    string         `json:"avatar"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Gender    int            `json:"gender" gorm:"default:0"`
	Birthday  *time.Time     `json:"birthday"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CourseType struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Sort        int            `json:"sort" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Course struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	TypeID      uint           `json:"type_id" gorm:"not null"`
	Description string         `json:"description"`
	Cover       string         `json:"cover"`
	Author      string         `json:"author"`
	Price       float64        `json:"price" gorm:"default:0"`
	IsFree      int            `json:"is_free" gorm:"default:1"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	CollectCount int           `json:"collect_count" gorm:"default:0"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	CourseType  CourseType     `json:"course_type,omitempty" gorm:"foreignKey:TypeID"`
}

type Chapter struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CourseID  uint           `json:"course_id" gorm:"not null"`
	Title     string         `json:"title" gorm:"not null"`
	Sort      int            `json:"sort" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Sections  []Section      `json:"sections,omitempty" gorm:"foreignKey:ChapterID"`
}

type Section struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ChapterID   uint           `json:"chapter_id" gorm:"not null"`
	CourseID    uint           `json:"course_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	ContentType string         `json:"content_type" gorm:"default:video"`
	VideoURL    string         `json:"video_url"`
	BookContent string        `json:"book_content"`
	Duration    int            `json:"duration" gorm:"default:0"`
	Sort        int            `json:"sort" gorm:"default:0"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Favorite struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_target"`
	FavoriteType string    `json:"favorite_type" gorm:"not null;uniqueIndex:idx_user_target"`
	TargetID     uint      `json:"target_id" gorm:"not null;uniqueIndex:idx_user_target"`
	CreatedAt    time.Time `json:"created_at"`
}

type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	SectionID uint           `json:"section_id" gorm:"not null"`
	CourseID  uint           `json:"course_id" gorm:"not null"`
	Content   string         `json:"content" gorm:"not null"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type UserBehavior struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	CourseID     uint      `json:"course_id" gorm:"not null;index"`
	BehaviorType string    `json:"behavior_type" gorm:"not null"`
	Rating       int       `json:"rating" gorm:"default:1"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type QuestionBank struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description"`
	QuestionCount int            `json:"question_count" gorm:"default:0"`
	Status        int            `json:"status" gorm:"default:1"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type Question struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	BankID       uint           `json:"bank_id" gorm:"not null;index"`
	QuestionType string         `json:"question_type" gorm:"not null"`
	Content      string         `json:"content" gorm:"not null"`
	Options      string         `json:"options" gorm:"type:json"`
	Answer       string         `json:"answer" gorm:"not null"`
	Analysis     string         `json:"analysis"`
	Score        int            `json:"score" gorm:"default:2"`
	Difficulty   int            `json:"difficulty" gorm:"default:2"`
	Status       int            `json:"status" gorm:"default:1"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Paper struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Title         string         `json:"title" gorm:"not null"`
	Description   string         `json:"description"`
	TotalScore    int            `json:"total_score" gorm:"default:100"`
	PassScore     int            `json:"pass_score" gorm:"default:60"`
	Duration      int            `json:"duration" gorm:"default:60"`
	QuestionCount int            `json:"question_count" gorm:"default:0"`
	Status        int            `json:"status" gorm:"default:1"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type PaperQuestion struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PaperID    uint      `json:"paper_id" gorm:"not null;index"`
	QuestionID uint      `json:"question_id" gorm:"not null;index"`
	Sort       int       `json:"sort" gorm:"default:0"`
	Score      int       `json:"score" gorm:"default:2"`
	CreatedAt  time.Time `json:"created_at"`
}

type Exam struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	PaperID     uint           `json:"paper_id" gorm:"not null;index"`
	Description string         `json:"description"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	Duration    int            `json:"duration" gorm:"default:60"`
	TotalScore  int            `json:"total_score" gorm:"default:100"`
	PassScore   int            `json:"pass_score" gorm:"default:60"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Paper       Paper          `json:"paper,omitempty" gorm:"foreignKey:PaperID"`
}

type ExamRecord struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	ExamID     uint           `json:"exam_id" gorm:"not null;index"`
	PaperID    uint           `json:"paper_id" gorm:"not null"`
	Score      int            `json:"score" gorm:"default:0"`
	TotalScore int            `json:"total_score" gorm:"default:100"`
	IsPass     int            `json:"is_pass" gorm:"default:0"`
	StartTime  *time.Time     `json:"start_time"`
	SubmitTime *time.Time     `json:"submit_time"`
	TimeSpent  int            `json:"time_spent" gorm:"default:0"`
	Status     string         `json:"status" gorm:"default:pending"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Exam       Exam           `json:"exam,omitempty" gorm:"foreignKey:ExamID"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type UserAnswer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RecordID   uint      `json:"record_id" gorm:"not null;index"`
	QuestionID uint      `json:"question_id" gorm:"not null"`
	UserAnswer string    `json:"user_answer"`
	IsCorrect  int       `json:"is_correct" gorm:"default:0"`
	Score      int       `json:"score" gorm:"default:0"`
	CreatedAt  time.Time `json:"created_at"`
}

type Forum struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserID     uint           `json:"user_id" gorm:"not null;index"`
	Title      string         `json:"title" gorm:"not null"`
	Content    string         `json:"content" gorm:"not null"`
	ViewCount  int            `json:"view_count" gorm:"default:0"`
	ReplyCount int            `json:"reply_count" gorm:"default:0"`
	Status     int            `json:"status" gorm:"default:1"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type ForumReply struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ForumID   uint           `json:"forum_id" gorm:"not null;index"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Content   string         `json:"content" gorm:"not null"`
	ParentID  uint           `json:"parent_id" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Announcement struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"not null"`
	Author    string         `json:"author"`
	IsTop     int            `json:"is_top" gorm:"default:0"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	ViewCount int            `json:"view_count" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Banner struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`
	Image     string         `json:"image" gorm:"not null"`
	Link      string         `json:"link"`
	Sort      int            `json:"sort" gorm:"default:0"`
	Status    int            `json:"status" gorm:"default:1"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
