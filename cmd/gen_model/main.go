package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
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

	ignore := make(map[string]string)
	ignore["gorm"] = "-"

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.GenerateModel("test"gen.WithMethod())
	g.ApplyBasic(
		// Generate struct `User` based on table `users`
		g.GenerateModel("users"),

		// Generate struct `Sessions` based on table `sessions`
		g.GenerateModel("sessions"),

		// Generate struct `Clients` based on table `clients`
		g.GenerateModel("clients"),

		// Generate struct `Tickets` based on table `tickets`
		g.GenerateModelAs("tickets", "Ticket", gen.FieldNew("OriginalPath", "string", ignore)),

		// Generate struct `Links` based on table `links`
		g.GenerateModel("links"),

		// Generate struct `Comments` based on table `comments`
		g.GenerateModel("comments"),

		// Generate struct `Histories` based on table `histories`
		g.GenerateModel("histories"),

		// Generate struct `NasRequests` based on table `nas_requests`
		g.GenerateModel("nas_requests"),

		// Generate struct `NasRequests` based on table `nas_servers`
		g.GenerateModel("nas_servers"),
	)

	// Generate the code
	g.Execute()
}
