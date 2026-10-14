package techpalace
import "strings"

// WelcomeMessage returns a welcome message for the customer.
func WelcomeMessage(customer string) string {
	// panic("Please implement the WelcomeMessage() function")
    cc := "Welcome to the Tech Palace, "+strings.ToUpper(customer)
    return cc
}

// AddBorder adds a border to a welcome message.
func AddBorder(welcomeMsg string, numStarsPerLine int) string {
	// panic("Please implement the AddBorder() function")
    c := strings.Repeat("*", numStarsPerLine)
    cc := c+"\n"+welcomeMsg+"\n"+c
    return cc
}

// CleanupMessage cleans up an old marketing message.
func CleanupMessage(oldMsg string) string {
	cleaned := strings.ReplaceAll(oldMsg, "*", "")
    return strings.TrimSpace(cleaned)
}
