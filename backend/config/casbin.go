package config

import (
	"log"

	"github.com/casbin/casbin/v3"
	casbinmodel "github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// ACL 模型：sub=帳號, obj=資源, act=動作
const casbinACL = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act
`

func NewEnforcer(db *gorm.DB) *casbin.Enforcer {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		log.Fatalf("[Casbin] adapter: %v", err)
	}
	m, err := casbinmodel.NewModelFromString(casbinACL)
	if err != nil {
		log.Fatalf("[Casbin] model: %v", err)
	}
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		log.Fatalf("[Casbin] enforcer: %v", err)
	}
	log.Println("[Casbin] ready")
	return enforcer
}
