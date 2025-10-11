package game

import (
	"math/rand"
	"time"
)

func generateMap(width, height, maxRooms int) ([][]*room, int, int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	worldMap := make([][]*room, height)
	for i := range worldMap {
		worldMap[i] = make([]*room, width)
	}

	startX, startY := width/2, height/2
	currentX, currentY := startX, startY
	roomsCreated := 0

	commonRoomTypes := []roomType{Empty, Enemy, Tresure, Shop}

	for roomsCreated < maxRooms {
		if worldMap[currentY][currentX] == nil {
			randomIndex := rand.Intn(len(commonRoomTypes))
			randomType := commonRoomTypes[randomIndex]
			worldMap[currentY][currentX] = &room{Type: randomType}
			roomsCreated++
		}

		dx, dy := 0, 0
		switch rand.Intn(4) {
		case 0:
			dy = -1
		case 1:
			dy = 1
		case 2:
			dx = -1
		case 3:
			dx = 1
		}

		if currentX+dx >= 0 && currentX+dx < width && currentY+dy >= 0 && currentY+dy < height {
			currentX += dx
			currentY += dy
		}
	}

	worldMap[startY][startX].Type = Empty

	for {
		randX, randY := rand.Intn(width), rand.Intn(height)
		if worldMap[randY][randX] != nil && (randX != startX || randY != startY) {
			worldMap[randY][randX].Type = StairsUp
			break
		}
	}

	return worldMap, startX, startY
}
