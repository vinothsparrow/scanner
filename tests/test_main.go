package tests

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vinothsparrow/scanner/config"
	"github.com/vinothsparrow/scanner/controllers"
)

func Test(t *testing.T) { Testing(t) }

var _ = Suite(&ScanSuite{})

type ScanSuite struct {
	config *viper.Viper
	router *gin.Engine
}

func (s *ScanSuite) SetUpTest(c *C) {
	config.Init("test")
	s.config = config.GetConfig()
	s.router = SetupRouter()
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	health := new(controllers.HealthController)
	v1 := router.Group("v1")
	{
		scanGroup := v1.Group("scan")
		{
			scan := new(controllers.ScanController)
			scanGroup.GET("/:id", scan.Retrieve)
		}
	}
	return router
}

func TestMain(m *testing.M) {
	SetupRouter()
}
