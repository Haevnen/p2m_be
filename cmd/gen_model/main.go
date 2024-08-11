package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"

	"github.com/Haevnen/p2m_be/internal/pkg/model"
)

func main() {
	conf := gen.Config{
		OutPath:      "../../internal/pkg/dal",
		ModelPkgPath: "../../internal/pkg/model",
		Mode:         gen.WithDefaultQuery,
	}
	g := gen.NewGenerator(conf)

	gormdb, _ := gorm.Open(mysql.Open("root:root@(127.0.0.1:33306)/p2m_db?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(gormdb) // reuse your gorm db

	type Querier interface {
		// GetByRoles query data by roles and return it as *slice of pointer*
		//   (The below blank line is required to comment for the generated method)
		//
		// SELECT * FROM @@table WHERE role IN @rolesName
		GetByRoles(rolesName ...string) ([]*gen.T, error)
	}

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.GenerateModel("test"gen.WithMethod())
	g.ApplyInterface(func(Querier) {}, model.User{},
		// Generate struct `User` based on table `users`
		// g.GenerateModel("users"),

		//
		// Generate struct `Sessions` based on table `users`
		g.GenerateModel("sessions"),
	)
	// Generate the code
	g.Execute()

}
