package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	conf := gen.Config{
		OutPath:      "tmp",
		ModelPkgPath: "internal/pkg/model",
		//Mode:    gen.WithoutContext, //gen.WithoutContext, //  gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode

	}
	g := gen.NewGenerator(conf)

	gormdb, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:33306)/p2m_db?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.GenerateModel("test"gen.WithMethod())
	g.ApplyBasic(
		// Generate struct `User` based on table `users`
		g.GenerateModel("users"),

		//
		// Generate struct `Employee` based on table `users`
		//g.GenerateModelAs("users", "Employee"),

		// Generate struct `User` based on table `users` and generating options
		//g.GenerateModel("users", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),
	)
	//g.ApplyBasic(
	// Generate structs from all tables of current database
	//g.GenerateAllTable()...,
	//)
	// Generate the code
	g.Execute()
}
