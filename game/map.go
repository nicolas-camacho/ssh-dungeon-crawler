package game

import (
	"image"
	"math/rand"
	"time"
)

func generateMap(width, height, maxRooms, floorNum int) (*floor, int, int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	worldMap := make([][]*room, height)
	for i := range worldMap {
		worldMap[i] = make([]*room, width)
	}

	var allRoomCoords []image.Point
	startX, startY := width/2, height/2
	currentX, currentY := startX, startY
	roomsCreated := 0

	for roomsCreated < maxRooms {
		if worldMap[currentY][currentX] == nil {
			worldMap[currentY][currentX] = &room{}
			allRoomCoords = append(allRoomCoords, image.Point{X: currentX, Y: currentY})
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

	totalRooms := len(allRoomCoords)

	rand.Shuffle(len(allRoomCoords), func(i, j int) {
		allRoomCoords[i], allRoomCoords[j] = allRoomCoords[j], allRoomCoords[i]
	})

	enemyRatio := 0.4 + rand.Float64()*0.2
	numEnemies := int(float64(totalRooms) * enemyRatio)
	assignedCount := 0
	for i := range numEnemies {
		coord := allRoomCoords[i]
		worldMap[coord.Y][coord.X].Type = Enemy
		assignedCount++
	}

	numTreasures := 1 + rand.Intn(3)
	for i := range numTreasures {
		if assignedCount+i < totalRooms {
			coord := allRoomCoords[assignedCount+i]
			worldMap[coord.Y][coord.X].Type = Tresure
		}
	}
	assignedCount += numTreasures

	if rand.Intn(2) == 0 {
		var potentialShopSpots []image.Point
		for _, coord := range allRoomCoords {
			room := worldMap[coord.Y][coord.X]
			if room.Type == 0 && isAdjacentToEnemy(coord.X, coord.Y, worldMap) {
				potentialShopSpots = append(potentialShopSpots, coord)
			}
		}

		if len(potentialShopSpots) > 0 {
			shopCoord := potentialShopSpots[rand.Intn(len(potentialShopSpots))]
			worldMap[shopCoord.Y][shopCoord.X].Type = Shop
		}
	}

	var startCoords image.Point
	if floorNum > 0 {
		var potentialDownStairsSpots []image.Point
		for _, coord := range allRoomCoords {
			if worldMap[coord.Y][coord.X].Type == 0 {
				potentialDownStairsSpots = append(potentialDownStairsSpots, coord)
			}
		}
		if len(potentialDownStairsSpots) > 0 {
			downStairsCoord := potentialDownStairsSpots[rand.Intn(len(potentialDownStairsSpots))]
			worldMap[downStairsCoord.Y][downStairsCoord.X].Type = StairsDown
			startCoords = downStairsCoord
		} else {
			startCoords = image.Point{X: startX, Y: startY}
		}
	} else {
		startCoords = image.Point{X: startX, Y: startY}
	}

	var potentialUpStairsSpots []image.Point
	for _, coord := range allRoomCoords {
		if worldMap[coord.Y][coord.X].Type == 0 {

			manhattanDistance := abs(coord.X-startX) + abs(coord.Y-startY)
			if manhattanDistance > 1 {
				potentialUpStairsSpots = append(potentialUpStairsSpots, coord)
			}
		}
	}
	if len(potentialUpStairsSpots) > 0 {
		upStairsCoord := potentialUpStairsSpots[rand.Intn(len(potentialUpStairsSpots))]
		worldMap[upStairsCoord.Y][upStairsCoord.X].Type = StairsUp
	}

	for _, coord := range allRoomCoords {
		if worldMap[coord.Y][coord.X].Type == 0 {
			worldMap[coord.Y][coord.X].Type = Empty
		}
	}

	worldMap[startY][startX].Type = Empty

	newFloor := &floor{
		worldMap: worldMap,
	}

	return newFloor, startCoords.X, startCoords.Y
}

func isAdjacentToEnemy(x, y int, worldMap [][]*room) bool {
	directions := []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	width, height := len(worldMap[0]), len(worldMap)

	for _, dir := range directions {
		checkX, checkY := x+dir.X, y+dir.Y
		if checkX >= 0 && checkX < width && checkY >= 0 && checkY < height {
			if neightbor := worldMap[checkY][checkX]; neightbor != nil && neightbor.Type == Enemy {
				return true
			}
		}
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
