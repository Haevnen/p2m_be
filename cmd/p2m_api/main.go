package main

import "C"
import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/configuration"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/core/shipping/bank_service"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/core/shipping/bank_service/repository"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/core/shipping/bank_service/usecase"
	"gitlab.citigo.com.vn/kship/kship-add-on/internal/pkg/models/model"
	_ "go.uber.org/automaxprocs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	_ "net/http/pprof"
	"time"
)

var config configuration.Config

func init() {
	viper.SetConfigFile(`configs/config_check_price.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func main() {

	ctx := context.Context(context.Background())
	gormdb, _ := gorm.Open(mysql.Open("admin:admin@(10.24.22.234:3306)/kvshipping_prelive?charset=utf8mb4&parseTime=True&loc=Local"))
	//var test repository.MysqlShopServiceRepository
	//test := repository2.NewMysqlBankServiceRepository(gormdb)
	//
	s := model.Bank{
		BankCode:   "TFG3",
		BankName:   "Test from Golang 3",
		CityID:     "ca nuoc",
		BranchName: "DONG NAM A (SEA BANK)-317",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	//err := test.Add(ctx, &s)
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println(s.ID)
	//}
	//fmt.Println("done test")

	bankRepo := repository.NewMysqlBankServiceRepository(gormdb)
	bankUc := usecase.NewBankServiceUseCase(bankRepo)

	//bankUc.

	u := interface{}(bankUc).(bank_service.UseCase)
	fmt.Println(u.TestBankUseCase())
	u.Add(ctx, &s)
	fmt.Println(s.ID)
	fmt.Println(u.CountAll(ctx))

}
