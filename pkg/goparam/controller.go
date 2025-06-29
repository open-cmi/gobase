package goparam

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) map[string]interface{} {
	//sess, exist := c.Get("session")
	u, exist := c.Get("user")
	if exist {
		// session := sess.(*sessions.Session)
		// user, ok := session.Values["user"].(map[string]interface{})
		// if ok {
		// 	return user
		// }
		user, ok := u.(map[string]interface{})
		if ok {
			return user
		}
		return user
	}

	return nil
}

// ParseParams parse param
func ParseParams(c *gin.Context) *Param {
	var option Param

	user := GetUser(c)
	if user != nil {
		userID := user["id"].(string)
		option.UserID = userID

		role := user["role"].(string)
		option.Role = role
	}

	devID := c.DefaultQuery("dev_id", "")
	option.DevID = devID

	pagestr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pagestr)
	if err != nil {
		page = 1
	}

	if page <= 0 {
		page = 1
	}

	option.PageParam.Page = page - 1

	// page size
	pagesizestr := c.DefaultQuery("page_size", "25")
	pagesize, err := strconv.Atoi(pagesizestr)
	if err != nil {
		pagesize = 25
	}
	option.PageParam.PageSize = pagesize

	option.OrderParam.Order = c.Query("order")
	option.OrderParam.OrderBy = c.Query("order_by")

	return &option
}
