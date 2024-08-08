// Package gormdb build connection to db with gorm and ddtrace
package gormdb

import (
	"testing"
)

func TestBuildMySQLConnectionString(t *testing.T) {
	type args struct {
		c *Config
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when location is default",
			args: args{c: &Config{DBUser: "user", DBPass: "password", DBHost: "dbhost.com", DBName: "dbname"}},
			want: "user:password@tcp(dbhost.com:3306)/dbname?parseTime=true",
		},
		{
			name: "when location is specified",
			args: args{c: &Config{DBUser: "user", DBPass: "password", DBHost: "dbhost.com", DBName: "dbname", DBLocation: "Asia/Tokyo"}},
			want: "user:password@tcp(dbhost.com:3306)/dbname?loc=Asia%2FTokyo&parseTime=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildMySQLConnectionString(tt.args.c)
			if err != nil {
				t.Errorf("BuildMySQLConnectionString() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("BuildMySQLConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}
