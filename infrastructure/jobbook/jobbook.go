package jobbook

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type OrderDbPram struct {
	server   string
	database string
	user     string
	password string
}

type Repository struct {
	sugar       *zap.SugaredLogger
	orderDbPram OrderDbPram
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
	orderDbPram OrderDbPram,
) *Repository {
	return &Repository{
		sugar:       sugar,
		orderDbPram: orderDbPram,
	}
}

func (r *Repository) FetchAll() []orderdb.JobBook {
	db, err := open(r.orderDbPram)
	if err != nil {
		r.sugar.Fatal("データベースに接続できませんでした", err)
	}

	jobBookModels := []JobBookModel{}
	result := db.Find(&jobBookModels, "納期 is not null AND 状態 = '受注'")
	if result.Error != nil {
		r.sugar.Fatal("M作業台帳を取得できませんでした", result.Error)
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

	return jobBooks
}

func open(orderDbPram OrderDbPram) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"sqlserver://%v:%v@%v?database=%v",
		orderDbPram.user,
		orderDbPram.password,
		orderDbPram.server,
		orderDbPram.database,
	)
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}
