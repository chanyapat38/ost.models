package controllers

import (
	"ost.models/model_ost_db/db/service"

	"fmt"
	"time"

	"ost.models/model_ost_db/structs"
	"ost.models/setting"
	"ost.models/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//DeleteController is for insert logic
type DeleteController struct{}

//DeleteDocument is for
func (u *DeleteController) DeleteDocument(c *gin.Context) (bool, interface{}) {
	var resultStatus bool
	var resultData interface{}
	var jsonbody structs.Jsonbody
	//Check if jsonbody is not following struck format
	if err := c.ShouldBindJSON(&jsonbody); err != nil {
		fmt.Println(err)
		c.JSON(401, err)
		return resultStatus, resultData
	}

	con, e := jsonbody.Condition.(map[string]interface{})
	if e {
	}
	condition := utils.ConvertOperators(con).(map[string]interface{})

	//Check if Multi is empty
	if jsonbody.Multi == nil {
		// not set
		c.JSON(401, gin.H{"error": "'Multi': required field is not set"})
	} else if !(*jsonbody.Multi) {
		// set to false
		data, e := jsonbody.Data.(map[string]interface{})
		if e {
		}
		update := utils.ConvertOperators(data).(map[string]interface{})
		//Check if Data is empty
		if len(data) == 0 {
			userservice := service.DeleteService{}
			id, err, coll := userservice.FindOneAndDelete(condition, jsonbody.Collection)
			if err != nil || !coll {
				if !coll {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": "Collection not found!"})
				} else {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": err.Error()})
				}
			} else {
				// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following users have deleted successfully", "results": id})
				resultStatus = true
				resultData = id
			}
		} else {

			//arrayFilters
			var arrayFilters []interface{}
			updateFilter, e := jsonbody.UpdateFilter.(map[string]interface{})
			if e {
			}
			// fmt.Println("updateFilter : ",updateFilter)
			obj := make(map[string]interface{})
			for k, v := range updateFilter {
				obj[k] = utils.ConvertOperators(v)
				// arrayFilters = append(arrayFilters,  map[string]interface{}{k:utils.ConvertOperators(v)})
			}
			arrayFilters = append(arrayFilters, obj)
			fmt.Println("arrayFilters : ", arrayFilters)
			// jsonString, _ := json.Marshal(arrayFilters)
			// fmt.Println(string(jsonString))
			lastupdate := bson.M{
				"last_updated": time.Now(),
			}
			update["$set"] = lastupdate
			userservice := service.DeleteService{}
			id, err, coll := userservice.FindOneAndUpdate(condition, arrayFilters, update, jsonbody.Collection)
			if err != nil || !coll {
				if !coll {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": "Collection not found!"})
				} else {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": err.Error()})
				}
			} else {
				// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following users have deleted successfully", "results": id})
				resultStatus = true
				resultData = id
			}
		}

	} else {
		data, e := jsonbody.Data.(map[string]interface{})
		if e {
		}
		update := utils.ConvertOperators(data).(map[string]interface{})
		//Check if Data is empty
		if len(data) == 0 {
			userservice := service.DeleteService{}
			id, err, coll := userservice.DeleteMany(condition, jsonbody.Collection)
			if err != nil || !coll {
				if !coll {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": "Collection not found!"})
				} else {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": err.Error()})
				}
			} else {
				// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following users have deleted successfully", "results": id})
				resultStatus = true
				resultData = id
			}
		} else {
			lastupdate := bson.M{
				"last_updated": time.Now(),
			}
			update["$set"] = lastupdate
			userservice := service.DeleteService{}
			id, err, coll := userservice.DeleteManyWithFilter(condition, update, jsonbody.Collection)
			if err != nil || !coll {
				if !coll {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": "Collection not found!"})
				} else {
					c.JSON(500, gin.H{"statusCode": setting.AppSetting.HTTP500, "message": "The following user haven’t deleted", "errors": err.Error()})
				}
			} else {
				// c.JSON(200, gin.H{"statusCode": setting.AppSetting.HTTP200, "message": "The following users have deleted successfully", "results": id})
				resultStatus = true
				resultData = id
			}
		}
	}

	return resultStatus, resultData
}
