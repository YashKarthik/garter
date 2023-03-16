package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
    screen, err := tcell.NewScreen()
    if err != nil {
        log.Fatal("Couldn't create screen\n", err);
    }
    if err := screen.Init(); err != nil {
        log.Fatal("Couldn't Init screen\n", err);
    }

    defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
    screen.SetStyle(defStyle);

    game := Game{
        Screen: screen,
    }

    go game.Run()

    for {
        switch event := game.Screen.PollEvent().(type) {
        case *tcell.EventResize:
            game.Screen.Sync()

        case *tcell.EventKey:
            switch key := event.Key(); key {
            case tcell.KeyCtrlC:
                game.Screen.Fini()
                os.Exit(0)

            case tcell.KeyUp:
                game.snakeBody.ChangeDirection(-1, 0)

            case tcell.KeyDown:
                game.snakeBody.ChangeDirection(1, 0)

            case tcell.KeyLeft:
                game.snakeBody.ChangeDirection(0, -1)

            case tcell.KeyRight:
                game.snakeBody.ChangeDirection(0, 1)
            }
        }
    }

}

type Game struct {
    Screen      tcell.Screen
    snakeBody   SnakeBody
    FoodPos     SnakePart
    Score       int
    GameOver    bool
}

func drawParts(s tcell.Screen, parts []SnakePart, style tcell.Style, foodPos SnakePart, foodStyle tcell.Style) {
    s.SetContent(foodPos.X, foodPos.Y, '\u25CF', nil, foodStyle)
    for _, part := range parts {
        s.SetContent(part.X, part.Y, '■', nil, style)
    }
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, text string) {
    row := y1
    col := x1
    style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

    for _, r := range text {
        s.SetContent(col, row, r, nil, style)
        col++

        if col >= x2 {
            row++
            col = x1
        }
        if row > y2 {
            break
        }
    }
}

func checkCollision(parts []SnakePart, otherPart SnakePart) bool {
    for _, part := range parts {
        if part.X == otherPart.X && part.Y == otherPart.Y {
            return true
        }
    }
    return false
}

func (g *Game) UpdateFoodPos(width int, height int) {
    g.FoodPos.X = rand.Intn(width)
    g.FoodPos.Y = rand.Intn(height)

    if g.FoodPos.Y == 1 && g.FoodPos.X < 10 {
        g.UpdateFoodPos(width, height)
    }
}

func (g *Game) Run() {
    defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
    g.Screen.SetStyle(defStyle);
    width, height := g.Screen.Size()
    g.snakeBody.ResetPosition(width, height)
    g.UpdateFoodPos(width, height)
    g.GameOver = false
    g.Score = 0
    snakeStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

    for {
        longerSnake := false
        g.Screen.Clear()

        if checkCollision(g.snakeBody.Parts[len(g.snakeBody.Parts)-1:], g.FoodPos) {
            g.UpdateFoodPos(width, height)
            longerSnake = true
            g.Score++
        }
        if checkCollision(g.snakeBody.Parts[:len(g.snakeBody.Parts)-1], g.snakeBody.Parts[len(g.snakeBody.Parts)-1]) {
            break
        }

        g.snakeBody.Update(width, height, longerSnake)
        drawParts(g.Screen, g.snakeBody.Parts, snakeStyle, g.FoodPos, snakeStyle)
        drawText(g.Screen, 1, 1, 8+len(strconv.Itoa(g.Score)), 1, "Score: "+strconv.Itoa(g.Score))
        time.Sleep(40 * time.Millisecond)
        g.Screen.Show()
    }

    g.GameOver = true
    drawText(g.Screen, width/2-20, height/2, width/2+20, height/2, "Game Over, Score: "+strconv.Itoa(g.Score)+", Play Again? y/n")
    g.Screen.Show()
}

type SnakePart struct {
    X int
    Y int
}

type SnakeBody struct {
    Parts       []SnakePart
    Xspeed      int
    Yspeed      int
}

func (sb *SnakeBody) ChangeDirection(vertical int, horizontal int) {
    sb.Yspeed = vertical
    sb.Xspeed = horizontal
}

func (sb *SnakeBody) Update(width int, height int, longerSnake bool) {
    sb.Parts = append(sb.Parts, sb.Parts[len(sb.Parts)-1].GetUpdatedPart(sb, width, height))
    if !longerSnake {
        sb.Parts = sb.Parts[1:]
    }
}

func (sb *SnakeBody) ResetPosition(width int, height int) {
    snakeParts := []SnakePart{
        { X: int(width / 2), Y: int(width / 2)},
        { X: int(width / 2) + 1, Y: int(width / 2) + 1},
        { X: int(width / 2) + 2, Y: int(width / 2) + 2},
    }
    
    sb.Parts = snakeParts
    sb.Xspeed = 1
    sb.Yspeed = 0
}

func (sp *SnakePart) GetUpdatedPart(sb *SnakeBody, width int, height int) SnakePart {
    newPart := *sp

    newPart.X = (newPart.X + sb.Xspeed) % width
    if newPart.X < 0 {
        newPart.X += width
    }

    newPart.Y = (newPart.Y + sb.Yspeed) % height
    if newPart.Y < 0 {
        newPart.Y += height
    }

    return newPart
}
