package controller

import (
	"Slot_booking/cache"
	"Slot_booking/entity"
	"Slot_booking/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type SlotController interface {
	FindAll(ctx *gin.Context, startDate time.Time, endTime time.Time) []entity.Slot
	AddSlot(ctx *gin.Context, m map[string]interface{}) error
}

type slotController struct {
	service   service.SlotService
	slotCache cache.SlotCache
}

func NewSlotController(service service.SlotService, cache cache.SlotCache) SlotController {
	return &slotController{
		service:   service,
		slotCache: cache,
	}
}

func (c *slotController) FindAll(ctx *gin.Context, todayTime time.Time, endTime time.Time) []entity.Slot {
	finalSlots := make([]entity.Slot, 0)
	dates := make([]string, 0)

	for date := todayTime; date.After(endTime) == false; date = date.Add(24 * time.Hour) {
		formatDate := date.Format("2006-01-02")

		key := fmt.Sprintf("slots_%v", formatDate)
		slots, err := c.slotCache.GetSlot(ctx, key)
		if err == nil {
			finalSlots = append(finalSlots, slots...)
		} else {
			dates = append(dates, formatDate)
		}
	}

	slots, _ := c.service.FindAll(dates)
	finalSlots = append(finalSlots, slots...)

	for i := 0; i < len(slots); i += 24 {
		date := slots[i].Date
		key := fmt.Sprintf("slots_%v", date)
		slotsForOne := make([]entity.Slot, 0)
		slotsForOne = slots[i : i+24]
		c.slotCache.SetSlot(ctx, key, slotsForOne)
	}
	return finalSlots
}

func (c *slotController) AddSlot(ctx *gin.Context, m map[string]interface{}) error {
	date := m["date"].(string)
	count, err := c.service.GetCount(date)
	if err != nil {
		return err
	}
	if count == 24 {
		err = errors.New("slots are already added")
		return err
	}

	slots := make([]entity.Slot, 24)
	for i := 0; i < 24; i++ {
		slots[i].Date = date
		slots[i].StartTime = entity.StartTimeOfSlot(i)
	}

	slots, err = c.service.AddSlot(slots)
	if err != nil {
		return err
	}
	return nil
}
