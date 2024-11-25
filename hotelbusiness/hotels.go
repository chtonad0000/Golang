//go:build !solution

package hotelbusiness

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	maxDate := 0
	for _, guest := range guests {
		if guest.CheckOutDate > maxDate {
			maxDate = guest.CheckOutDate
		}
	}
	days := make([]int, maxDate+1)
	load := make([]Load, 0)
	for _, guest := range guests {
		for i := guest.CheckInDate; i < guest.CheckOutDate; i++ {
			days[i]++
		}
	}
	prev := 0
	for i := 0; i < len(days); i++ {
		if days[i] != prev {
			load = append(load, Load{i, days[i]})
			prev = days[i]
		}
	}
	return load
}
