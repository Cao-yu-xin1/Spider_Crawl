package ordersn

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func OrderNo() string {
	timeStr := time.Now().Format("20060102150405")
	uuidStr := uuid.NewString()[:8]
	return fmt.Sprintf("%s%s", timeStr, uuidStr)
}
