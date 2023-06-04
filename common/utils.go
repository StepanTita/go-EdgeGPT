package common

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func RunEvery(d time.Duration, fs ...func() error) error {
	if err := runFuncs(CurrentTimestamp(), fs...); err != nil {
		return errors.Wrap(err, "failed to run funcs initial")
	}
	for x := range time.Tick(d) {
		if err := runFuncs(x, fs...); err != nil {
			return errors.Wrap(err, "failed to run funcs")
		}
	}
	return nil
}

func runFuncs(x time.Time, fs ...func() error) error {
	for i, f := range fs {
		if err := f(); err != nil {
			return errors.Wrapf(err, "failed to run function: %v -> %v", x, i)
		}
	}
	return nil
}

// CurrentTimestamp is a utility method to make sure UTC time is used all over the code
func CurrentTimestamp() time.Time {
	return time.Now().UTC()
}

func RandFromRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func MustGenerateRandomHex(length int) string {
	bytes := make([]byte, length/2)
	_, err := crand.Read(bytes)
	if err != nil {
		logrus.Fatal("failed to generate rand bytes")
		return ""
	}
	return hex.EncodeToString(bytes)
}
