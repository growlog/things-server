package controllers

import (
	"github.com/growlog/things-server/internal/models"
	"github.com/growlog/things-server/internal/services"
)


type ThingServer struct{
    DAL *models.DataAccessLayer
	IAM *services.IAMClient
}
