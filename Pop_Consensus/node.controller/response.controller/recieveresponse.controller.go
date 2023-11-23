package responsecontroller

import (
	"pop_v1/models"

	"github.com/gofiber/fiber/v2"
)

func Groupresponse(c *fiber.Ctx) error {

	var response models.Response
	if err:=c.BodyParser(&response);err!=nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	models.Lock.Lock()
	models.Groupmap[response.Ip]=response
	models.VoteCount++
	models.Lock.Unlock()
	
	return c.Status(200).JSON("Response Send to the admin")
	
}