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

func updateGuest(date int, inc int, load *[]int) {
	for len(*load) <= date {
		*load = append(*load, 0)
	}
	(*load)[date] += inc
}

func buildDiffArray(load []int) {
	for i := 1; i < len(load); i++ {
		load[i] += load[i-1]
	}
}

func ComputeLoad(guests []Guest) []Load {
	intLoads := make([]int, 0)
	for _, guest := range guests {
		updateGuest(guest.CheckInDate, 1, &intLoads)
		updateGuest(guest.CheckOutDate, -1, &intLoads)
	}
	buildDiffArray(intLoads)
	loads := make([]Load, 0)
	lstCount := 0
	for date, count := range intLoads {
		if lstCount != count {
			lstCount = count
			loads = append(loads, Load{StartDate: date, GuestCount: count})
		}
	}
	return loads
}
