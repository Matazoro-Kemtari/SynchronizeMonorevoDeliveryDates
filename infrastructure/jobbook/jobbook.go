package jobbook

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type OrderDbConfig struct {
	server   string
	database string
	user     string
	password string
}

func NewOrderDbConfig() *OrderDbConfig {
	return &OrderDbConfig{
		os.Getenv("DB_SERVER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	}
}

type Repository struct {
	sugar         *zap.SugaredLogger
	orderDbConfig *OrderDbConfig
}

type JobBookModel struct {
	WorkedNumber string    `gorm:"column:作業NO;unique;not null"`
	DeliveryDate time.Time `gorm:"column:納期"`
}

// JobBookのテーブル名を定義する
func (JobBookModel) TableName() string {
	return "M作業台帳"
}

func NewRepository(
	sugar *zap.SugaredLogger,
	orderDbConfig *OrderDbConfig,
) *Repository {
	return &Repository{
		sugar:         sugar,
		orderDbConfig: orderDbConfig,
	}
}

func (r *Repository) FetchAll() ([]orderdb.JobBook, error) {
	db, err := open(r.orderDbConfig)
	if err != nil {
		r.sugar.Error("データベースに接続できませんでした", err)
		return nil, fmt.Errorf("データベースに接続できませんでした error: %v", err)
	}

	jobBookModels := []JobBookModel{}
	result := db.Find(&jobBookModels, "納期 is not null AND 状態 = '受注'")
	if result.Error != nil {
		m := fmt.Sprintf("M作業台帳を取得できませんでした error: %v", result.Error)
		r.sugar.Error(m)
		return nil, fmt.Errorf(m)
	}
	fmt.Println("jobBook:", jobBookModels)

	// domain.modelに詰め替え
	jobBooks := []orderdb.JobBook{}
	for _, v := range jobBookModels {
		jobBooks = append(
			jobBooks,
			orderdb.JobBook{
				WorkedNumber: v.WorkedNumber,
				DeliveryDate: v.DeliveryDate,
			},
		)
	}

	return jobBooks, nil
}

func open(orderDbPram *OrderDbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%v:%v@%v?database=%v",
		orderDbPram.user,
		orderDbPram.password,
		orderDbPram.server,
		orderDbPram.database,
	)
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}
