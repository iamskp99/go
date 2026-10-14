package cars

// CalculateWorkingCarsPerHour calculates how many working cars are
// produced by the assembly line every hour.
func CalculateWorkingCarsPerHour(productionRate int, successRate float64) float64 {
    c := successRate/float64(100)
    res := c*float64(productionRate)
    return res
}

// CalculateWorkingCarsPerMinute calculates how many working cars are
// produced by the assembly line every minute.
func CalculateWorkingCarsPerMinute(productionRate int, successRate float64) int {
    cc := CalculateWorkingCarsPerHour(productionRate,successRate)/float64(60)
    res := int(cc)
	return res
}

// CalculateCost works out the cost of producing the given number of cars.
func CalculateCost(carsCount int) uint {
	g := carsCount/10
    c := carsCount%10
    return uint((g*95000)+(c*10000))
}
