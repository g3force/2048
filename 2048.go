package main

import "fmt"
import "strings"
import "strconv"
import "math/rand"
import "github.com/nsf/termbox-go"
import "os"

type Field struct {
  field [][]int
  width,height int
  cellWidth int
}

func (f *Field) emptyline() string {
  elPlaceholder := strings.Repeat("-",f.cellWidth)
  return strings.Repeat("|"+elPlaceholder,f.width)+"|"
}

func (f *Field) Print() {
  fmtLine := "|%"+strconv.Itoa(f.cellWidth)+"d"
  for h:=0;h<f.height;h++ {
    fmt.Printf(f.emptyline()+"\n")
    for c:=range f.field[h] {
      fmt.Printf(fmtLine, f.field[h][c])
    }
    fmt.Printf("|\n")
  }
  fmt.Printf(f.emptyline()+"\n")
}

func log2(val int) int {
  i:=0
  for ;val > 1;i++ {
    val /= 2
  }
  return i
} 

func getColor(val int) termbox.Attribute {
  if val == 0 {
    return termbox.ColorBlack
  }
  logVal := log2(val)
  
  colors := []termbox.Attribute{
  termbox.ColorRed,
  termbox.ColorGreen,
  termbox.ColorYellow,
  termbox.ColorBlue,
  termbox.ColorMagenta,
  termbox.ColorCyan,
  termbox.ColorWhite}
  return colors[logVal%len(colors)]
}

func (f *Field) Print_tb() {
  fmtLine := "%"+strconv.Itoa(f.cellWidth)+"d"
  h_off := 1
  for h:=0;h<f.height;h++ {
    printf_tb(1, h_off, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorBlack, 
      f.emptyline())
    h_off++
    offset := 1
    for c:=range f.field[h] {
      printf_tb(offset, h_off, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorBlack, 
        "|")
      printf_tb(offset+1, h_off, termbox.ColorWhite|termbox.AttrBold, getColor(f.field[h][c]), 
        fmtLine, f.field[h][c])
      offset += 1 + f.cellWidth
    }
    printf_tb(offset, h_off, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorBlack, 
      "|")
    h_off++
  }
  
  printf_tb(1, h_off, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorBlack, 
    f.emptyline())
    h_off++
}

func (f *Field) init() {
  f.field = make([][]int, f.height)
  for i := range f.field {
    f.field[i] = make([]int, f.width);
    for j := range f.field[i] {
      f.field[i][j] = 0;
    }
  }
  f.putRandom24();
  f.putRandom24();
}

func (f *Field) putRandom(val int) {
  for i:=0;i<100000;i++{
    x := rand.Int31n(int32(f.width))
    y := rand.Int31n(int32(f.height))
    if f.field[y][x] == 0 {
      f.field[y][x] = val
      return
    }
  }
  fmt.Println("Could not place any new tile. You've lost")
  os.Exit(1);
}

func (f *Field) up() {
  for col := 0 ; col < f.width ; col++ {
    vec := make([]int, f.height)
    for row :=0 ; row < f.height ; row++ {
      vec[row] = f.field[row][col]
    }
    res := merge(vec)
    for row :=0 ; row < f.height ; row++ {
      f.field[row][col] = res[row]
    }
  }
}

func (f *Field) down() {
  for col := 0 ; col < f.width ; col++ {
    vec := make([]int, f.height)
    for row :=0 ; row < f.height ; row++ {
      vec[f.height-row-1] = f.field[row][col]
    }
    res := merge(vec)
    for row :=0 ; row < f.height ; row++ {
      f.field[row][col] = res[f.height-row-1]
    }
  }
}

func (f *Field) left() {
  for col := 0 ; col < f.width ; col++ {
    f.field[col] = merge(f.field[col])
  }
}

func (f *Field) right() {
  for col := 0 ; col < f.width ; col++ {
    f.field[col] = reverse(merge(reverse(f.field[col])))
  }
}

type move func()

func (f *Field) lost() bool {
  for y:=0; y<f.height; y++ {
    for x:=0; x<f.width; x++ {
      if (f.field[y][x] == 0) {
        return false
      }
    }
  }
  return true
}

func (f *Field) makeMove(fnMove move) {
  if f.lost() {
    printf_tb(0, f.height*2+3, termbox.ColorWhite|termbox.AttrBold, termbox.ColorBlack, "You've lost.")
    os.Exit(0)
  }
  field := copy(f.field)
  fnMove()
  for y:=0; y<f.height; y++ {
    for x:=0; x<f.width; x++ {
      if field[y][x] != f.field[y][x] {
        f.putRandom24();
        return;
      }
    }
  }
}

func copy(arr [][]int) [][]int {
  narr := make([][]int, len(arr))
  for i:=0;i<len(arr);i++ {
    narr[i] = make([]int, len(arr[i]))
    for j:=0;j<len(arr[i]);j++ {
      narr[i][j] = arr[i][j]
    }
  }
  return narr
}

func reverse(vec []int) []int {
  res := make([]int, len(vec))
  for i:=0; i<len(vec); i++ {
    res[len(vec)-i-1] = vec[i]
  }
  return res
}

func merge(vec []int) []int {
  offset := 0
  res := make([]int,len(vec)) 
  for i:=0;i<len(vec);i++ {
    if vec[i] != 0 {
      res[offset] = vec[i]
      if offset > 0 && res[offset-1] == vec[i] {
        res[offset-1] *= 2
        res[offset] = 0
      }
      offset++
    }
  }
  offset = 0
  for i:=0;i<len(res);i++ {
    if res[i] != 0 {
      res[offset] = res[i]
      offset++
    }
  }
  for i:=offset;i<len(vec);i++ {
    res[i] = 0
  }
  return res
}

func (f *Field) putRandom24() {
  val := rand.Int31n(2)
  if val == 0 {
    f.putRandom(2)
  } else {
    f.putRandom(4)
  }
}


func print_tb(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func printf_tb(x, y int, fg, bg termbox.Attribute, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	print_tb(x, y, fg, bg, s)
}

func main() {
  var f Field;
  f.width = 4
  f.height = 4
  f.cellWidth = 6
  f.init();
  
  err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	//termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
  f.Print_tb()
  termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyArrowUp {
        f.makeMove(f.up)
			} else if ev.Key == termbox.KeyArrowLeft {
        f.makeMove(f.left)
			} else if ev.Key == termbox.KeyArrowRight {
        f.makeMove(f.right)
			} else if ev.Key == termbox.KeyArrowDown {
        f.makeMove(f.down)
			} else if ev.Key == termbox.KeyEsc {
        return
      }
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
      f.Print_tb()
			termbox.Flush()
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
      f.Print_tb()
			termbox.Flush()
		case termbox.EventMouse:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
      f.Print_tb()
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
