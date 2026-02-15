package formatting

import (
	"dickobrazz/src/shared/datetime"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var glitchMarks = []rune{
	'\u0335', '\u0336', '\u0337', '\u0338',
	'\u0300', '\u0301', '\u0302', '\u0303',
	'\u0304', '\u0305', '\u0306', '\u0307',
	'\u0308', '\u0309', '\u030A', '\u030B',
	'\u0310', '\u0311', '\u0312', '\u0313',
	'\u0334', '\u034F', '\u0350', '\u0351',
	'\u0352', '\u0353', '\u0354', '\u0355', '\u0356',
}

var mathFancy = map[int]string{
	0: "sin(0)", 1: "0!", 2: "C(2,1)", 3: "1! + 2!", 4: "2²", 5: "√25",
	6: "3!", 7: "3! + 1", 8: "2³", 9: "3²", 10: "C(5,2)", 11: "(1011)₂",
	12: "4! / 2", 13: "F₇", 14: "Cat₄", 15: "C(6,2)", 16: "2⁴", 17: "√289",
	18: "3! · 3", 19: "3³ − 2³", 20: "5! / 6", 21: "F₈", 22: "⌊π^e⌋", 23: "⌈π^e⌉",
	24: "4!", 25: "5²", 26: "4! + 2!", 27: "3³", 28: "T₇ = 7·8/2", 29: "2⁵ − 3",
	30: "2 · 5!!", 31: "2⁵ − 1", 32: "2⁵", 33: "4! + 3! + 2! + 0!", 34: "F₉",
	35: "C(7,3)", 36: "6²", 37: "⌊12π⌋", 38: "(100110)₂", 39: "3³ + 2·3!",
	40: "5! / 3", 41: "n² + n + 41 |_{n=0}", 42: "Cat₅", 43: "⌊14π⌋", 44: "⌊√2000⌋",
	45: "C(10,2)", 46: "4! + 4! − 2!", 47: "⌊15π⌋", 48: "4! · 2", 49: "7²",
	50: "⌊16π⌋", 51: "4! + 3³", 52: "6!! + 2²", 53: "⌊17π⌋", 54: "3³ + 3³",
	55: "F₁₀", 56: "C(8,3)", 57: "4! + 3! + 3³", 58: "6!! + C(5,2)", 59: "⌊19π⌋",
	60: "5! / 2", 61: "√3721",
}

var (
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rndMu sync.Mutex
)

func isMathDay(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 14
}

func isProgrammersDay(t time.Time) bool {
	return t.YearDay() == 256
}

func toProgrammersNotation(n int) string {
	rndMu.Lock()
	useBinary := rnd.Intn(2) == 0
	rndMu.Unlock()

	if useBinary {
		if n < 0 {
			return "-0b" + strconv.FormatUint(uint64(-n), 2)
		}
		return "0b" + strconv.FormatUint(uint64(n), 2)
	}
	if n < 0 {
		return fmt.Sprintf("-0x%X", -n)
	}
	return fmt.Sprintf("0x%X", n)
}

func glitchify(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		sb.WriteRune(ch)
		rndMu.Lock()
		count := rnd.Intn(3) + 1
		marks := make([]rune, count)
		for i := 0; i < count; i++ {
			marks[i] = glitchMarks[rnd.Intn(len(glitchMarks))]
		}
		rndMu.Unlock()

		for _, mark := range marks {
			sb.WriteRune(mark)
		}
	}
	return sb.String()
}

func fancyMathOrDefault(n int) string {
	if s, ok := mathFancy[n]; ok {
		return s
	}
	return strconv.Itoa(n)
}

// FormatCockSizeForDate форматирует размер в зависимости от текущей даты
func FormatCockSizeForDate(size int) string {
	displaySize := size
	now := datetime.NowTime()

	if now.Month() == time.April && now.Day() == 1 {
		displaySize = -size
	}

	if isMathDay(now) {
		return fancyMathOrDefault(displaySize)
	}

	if isProgrammersDay(now) {
		return toProgrammersNotation(displaySize)
	}

	if now.Month() == time.October && now.Day() == 31 {
		return glitchify(strconv.Itoa(displaySize))
	}

	return strconv.Itoa(displaySize)
}
