package routes

import (
	"apis-sati/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Routes() *fiber.App {
	crs := cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, DELETE, PUT, OPTIONS",
		AllowHeaders:     "*",
		AllowCredentials: true,
	})
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // Limite de 10Mb para envio de arquivos
	})
	app.Use(crs)

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// AUTHENTICATION ROUTES ==============================================================
			v1.Post("/auth/signin", controllers.SignIn)
			v1.Post("/auth/generate_code", controllers.GenerateCodeRecoverPassword)
			v1.Post("/auth/confirm_code", controllers.ConfirmCode)
			v1.Put("/auth/recover_password", controllers.RecoverPassword)

			// DASHBOARD ROUTES ========================================================================
			v1.Get("/dashboard/units/totalizators", controllers.TotalizatorsUnitsDashboard)
			v1.Get("/dashboard/units/graphic", controllers.GraphicUnitsDashboard)
			v1.Get("/dashboard/patients/totalizators", controllers.TotalizatorsPatientsDashboard)
			v1.Get("/dashboard/patients/graphic", controllers.GraphicPatientsDashboard)
			v1.Get("/dashboard/monthly_occupation/graphic", controllers.GraphicMonthlyDashboard)

			// UNITS ROUTES ============================================================================
			v1.Get("/units/list", controllers.ListUnits)
			v1.Get("/units/totalizators/:id", controllers.TotalizatorsUnit)
		}
	}

	return app
}
