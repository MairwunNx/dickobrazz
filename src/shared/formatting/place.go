package formatting

import (
	"dickobrazz/src/shared/datetime"
	"fmt"
	"time"
)

func GetPlaceEmoji(place int, isCurrentUser bool) string {
	switch place {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		now := datetime.NowTime()
		month := now.Month()

		var emoji string
		switch month {
		case time.March, time.April, time.May:
			emoji = "🫠"
		case time.June, time.July, time.August:
			emoji = "🥵"
		case time.September, time.October, time.November:
			emoji = "🤧"
		default:
			emoji = "🥶"
		}

		if isCurrentUser {
			return fmt.Sprintf("%s *%d*\\.", emoji, place)
		}
		return fmt.Sprintf("%s %d\\.", emoji, place)
	}
}

func GetPlaceEmojiForContext(place int, bold bool) string {
	if bold {
		return fmt.Sprintf("🥀 *%d*\\.", place)
	}
	return fmt.Sprintf("🥀 %d\\.", place)
}

func GetMedalByPosition(position int) string {
	switch position {
	case 0:
		return "🥇"
	case 1:
		return "🥈"
	case 2:
		return "🥉"
	default:
		return ""
	}
}
