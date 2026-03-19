package runtime

import "fmt"

type AdaptiveController struct {
	Enabled bool
}

func NewAdaptiveController() *AdaptiveController {
	return &AdaptiveController{
		Enabled: false,
	}
}

// включаем overlap только при деградации
func (a *AdaptiveController) Evaluate(current Transport) bool {
	if current.Score < 50 {
		if !a.Enabled {
			fmt.Println("[ADAPTIVE] degradation detected → enabling overlap")
		}
		a.Enabled = true
		return true
	}

	if a.Enabled {
		fmt.Println("[ADAPTIVE] path stable → disabling overlap")
	}

	a.Enabled = false
	return false
}