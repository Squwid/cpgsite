package cpgaws

import (
	"testing"
)

func TestSendInquiry(t *testing.T) {
	inquiry := CreateInquiry("bwhitelaw@sbcglobal.net", "248-752-0682", "Ben Whitelaw", "There are some issues that we need to speak about\n\n\n. <h2> Lets meet saturday </h2>")
	err := inquiry.Send()
	if err != nil {
		t.Logf("error occurred sending test inquiry : %v\n", err)
		t.Fail()
	}
}
