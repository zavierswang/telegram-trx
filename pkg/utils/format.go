package utils

import (
	"fmt"
	"telegram-trx/pkg/core/cst"
	"time"
)

func Duration(blocks int64) string {
	var str string
	durationHours := blocks * 3 / 3600
	d := durationHours / 24
	h := durationHours % 24
	if d >= 1 {
		if h > 0 {
			str = fmt.Sprintf("%d天%d小时", d, h)
		} else {
			str = fmt.Sprintf("%d天", d)
		}
	} else {
		str = fmt.Sprintf("%d小时", h)
	}
	return str
}

func DurationSec(days int64) string {
	var str string
	durationHours := days / 3600
	d := durationHours / 24
	h := durationHours % 24
	if d >= 1 {
		if h > 0 {
			str = fmt.Sprintf("%d天%d小时", d, h)
		} else {
			str = fmt.Sprintf("%d天", d)
		}
	} else {
		str = fmt.Sprintf("%d小时", h)
	}
	return str
}

func DateTime(t time.Time) string {
	return t.Format(cst.DateTimeFormatter)
}

func EnergyCount(energy int64) string {
	var count float64
	count = float64(energy) / 32000
	return fmt.Sprintf("%.1f", count)
}

func BalanceAdmin(balance string) string {
	return fmt.Sprintf("3️⃣ TRX钱包余额：%s TRX", balance)
}

func AddIndex(idx int) string {
	listIcon := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
	return listIcon[idx]
}

func FormatTime(t time.Time) string {
	return t.Format(cst.DateTimeFormatter)
}
